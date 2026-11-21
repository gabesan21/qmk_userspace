//go:build linux

package main

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestDiscoveryUsesOnlyDedicatedUSBConsole(t *testing.T) {
	classDir, devDir := t.TempDir(), t.TempDir()
	const usbID = "0003:0000A8F8:00001836"
	const console = "\x06\x31\xff\x09\x74\xa1\x01\x09\x75\x81\x02\xc0"
	add := func(node, hidID, desc string) {
		t.Helper()
		sys := filepath.Join(classDir, node, "device")
		if err := os.MkdirAll(sys, 0700); err != nil {
			t.Fatal(err)
		}
		for path, data := range map[string]string{
			filepath.Join(sys, "uevent"):            "DRIVER=hid-generic\nHID_ID=" + hidID + "\n",
			filepath.Join(sys, "report_descriptor"): desc,
			filepath.Join(devDir, node):             node,
		} {
			if err := os.WriteFile(path, []byte(data), 0600); err != nil {
				t.Fatal(err)
			}
		}
	}
	add("hidraw0", "0003:0000A8F8:00000001", console)
	add("hidraw1", "0005:0000A8F8:00001836", console)                // Bluetooth
	add("hidraw2", usbID, "\x06\x60\xff\x09\x61\xa1\x01\xc0")        // VIA
	add("hidraw3", usbID, "\x05\x01\x09\x06\xa1\x01"+console+"\xc0") // keyboard with nested header
	add("hidraw4", usbID, strings.Replace(console, "\x74", "\x75", 1))
	add("hidraw5", usbID, strings.Replace(console, "\xa1\x01", "\xa1\x00", 1))
	add("hidraw6", usbID, console+"\x05\x01\x09\x02\xa1\x01\xc0") // sibling pointer
	add("hidraw7", usbID, console[:len(console)-1])
	if dev, err := openConsole(classDir, devDir); !errors.Is(err, os.ErrNotExist) {
		if dev != nil {
			dev.Close()
		}
		t.Fatalf("expected no console, got %v, %v", dev, err)
	}
	add("hidraw83", usbID, console)
	dev, err := openConsole(classDir, devDir)
	if err != nil {
		t.Fatal(err)
	}
	defer dev.Close()
	if filepath.Base(dev.Name()) != "hidraw83" {
		t.Fatalf("unexpected device %s", dev.Name())
	}
}

func TestLayerDecoder(t *testing.T) {
	var decoder layerDecoder
	var layers []int
	for _, part := range []struct {
		data  string
		valid bool
	}{
		{"DL\x00M", false},
		{"1:\x000\n\x00DLM1:6\nDLM1:6\n", true},
		{"DLM1:7\nDLM1:10\nDLM1:-1\nDLM1:3\r\nDLM1:\n", false},
		{"noiseDLM1:5\n" + strings.Repeat("x", 3000), false},
		{"DLM1:2\nDL", false}, // Discard the oversized line, not just its prefix.
		{"M1:2\n", true},
	} {
		got := decoder.feed([]byte(part.data), func(id int) { layers = append(layers, id) })
		if got != part.valid {
			t.Fatalf("valid = %t, want %t", got, part.valid)
		}
	}
	if !reflect.DeepEqual(layers, []int{0, 6, 6, 2}) {
		t.Fatalf("decoded layers: %v", layers)
	}
}

type scriptedConsole struct {
	reports   []string
	end       error
	deadlines int
	closed    bool
}

func (s *scriptedConsole) Read(buf []byte) (int, error) {
	if len(s.reports) == 0 {
		return 0, s.end
	}
	n := copy(buf, s.reports[0])
	s.reports = s.reports[1:]
	return n, nil
}

func (s *scriptedConsole) SetReadDeadline(time.Time) error { s.deadlines++; return nil }
func (s *scriptedConsole) Close() error                    { s.closed = true; return nil }

func TestWatchConsolesLosesAndRecoversSignal(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	first := &scriptedConsole{reports: []string{"DLM1:3\n", "DLM1:3\n", "noise\nDLM1:"}, end: os.ErrDeadlineExceeded}
	second := &scriptedConsole{reports: []string{"4\n", "DLM1:6\n"}, end: io.EOF}
	attempt := 0
	open := func() (consoleReader, error) {
		attempt++
		switch attempt {
		case 1, 2:
			return nil, os.ErrNotExist
		case 3, 6:
			return nil, os.ErrPermission
		case 4:
			return first, nil
		case 5:
			return second, nil
		default:
			t.Fatal("unexpected reconnect")
			return nil, os.ErrNotExist
		}
	}
	var got []layerEvent
	watchConsoles(ctx, open, 0, time.Hour, func(event layerEvent) {
		got = append(got, event)
		if attempt == 6 {
			cancel()
		}
	})
	want := []layerEvent{
		{Layer: -1, Status: awaitDevice}, {Layer: -1, Status: needAccess},
		{Layer: -1, Status: awaitSignal}, {Layer: 3},
		{Layer: -1, Status: awaitSignal}, {Layer: 6},
		{Layer: -1, Status: awaitDevice}, {Layer: -1, Status: needAccess},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("events = %v, want %v", got, want)
	}
	if first.deadlines != 3 || second.deadlines != 2 || !first.closed || !second.closed {
		t.Fatalf("expected deadline renewal only on valid frames and closed sessions; got %+v, %+v", first, second)
	}
}

type waitingPipe struct {
	*os.File
	entered chan struct{}
}

func (p waitingPipe) Read(buf []byte) (int, error) {
	close(p.entered)
	return p.File.Read(buf)
}

func TestCancelInterruptsSilentReader(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pipe := waitingPipe{r, make(chan struct{})}
	go func() { <-pipe.entered; cancel() }()
	if err := readConsole(ctx, pipe, time.Minute, func(int) { t.Error("unexpected layer") }); !errors.Is(err, os.ErrClosed) {
		t.Fatalf("expected cancellation to close pipe, got %v", err)
	}
}
