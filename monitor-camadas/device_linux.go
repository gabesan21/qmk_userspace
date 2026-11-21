package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

const (
	awaitDevice = "Aguardando o Dilemma V3."
	awaitSignal = "Aguardando sinal do firmware."
	needAccess  = "Sem acesso ao console. Ajuste a regra udev e reconecte o teclado."
	readFailed  = "Falha na leitura do console. Tentando novamente."
)

// QMK declares its console as a dedicated top-level Application collection.
// Match its header at an item boundary, never a byte pattern inside item data.
func consoleDescriptor(desc []byte) bool {
	header := []byte{0x06, 0x31, 0xff, 0x09, 0x74, 0xa1, 0x01}
	if !bytes.HasPrefix(desc, header) {
		return false
	}
	depth := 1
	for desc = desc[len(header):]; len(desc) > 0; {
		item := desc[0]
		size := int(item & 3)
		if size == 3 {
			size = 4
		}
		if item == 0xfe || len(desc) < 1+size {
			return false // Unknown long item or truncated descriptor: fail closed.
		}
		switch item & 0xfc {
		case 0xa0: // Collection
			if size != 1 || desc[1] == 1 {
				return false // Never read another Application multiplexed on this node.
			}
			depth++
		case 0xc0: // End Collection
			if size != 0 {
				return false
			}
			depth--
			if depth == 0 {
				return len(desc) == 1 // No sibling collection (keyboard, pointer or VIA).
			}
		}
		desc = desc[1+size:]
	}
	return false
}

func consoleIdentity(uevent []byte) bool {
	// Linux HID_ID is BUS:VID:PID; 0003 is USB (not Bluetooth).
	for len(uevent) > 0 {
		var line []byte
		line, uevent, _ = bytes.Cut(uevent, []byte{'\n'})
		if bytes.EqualFold(line, []byte("HID_ID=0003:0000A8F8:00001836")) {
			return true
		}
	}
	return false
}

func openConsole(classDir, devDir string) (*os.File, error) {
	entries, err := os.ReadDir(classDir)
	if err != nil {
		return nil, err
	}
	var scanErr error = os.ErrNotExist
	for _, entry := range entries {
		classNode := filepath.Join(classDir, entry.Name())
		dev, err := openCandidate(classNode, filepath.Join(devDir, entry.Name()))
		if err == nil {
			return dev, nil
		}
		if !errors.Is(err, os.ErrNotExist) && !errors.Is(scanErr, os.ErrPermission) {
			scanErr = err
		}
	}
	return nil, scanErr
}

func openCandidate(classNode, devNode string) (*os.File, error) {
	before, err := os.Stat(classNode)
	if err != nil {
		return nil, err
	}
	device := filepath.Join(classNode, "device")
	identity, err := os.ReadFile(filepath.Join(device, "uevent"))
	if err != nil {
		return nil, err
	}
	if !consoleIdentity(identity) {
		return nil, os.ErrNotExist
	}
	desc, err := os.ReadFile(filepath.Join(device, "report_descriptor"))
	if err != nil {
		return nil, err
	}
	if !consoleDescriptor(desc) {
		return nil, os.ErrNotExist
	}
	dev, err := os.OpenFile(devNode, os.O_RDONLY|syscall.O_NONBLOCK|syscall.O_NOFOLLOW, 0)
	if err != nil {
		return nil, err
	}
	// A hidraw number can be reused on unplug. Never read a replacement node
	// if its sysfs identity changed between discovery and opening the device.
	after, err := os.Stat(classNode)
	if err != nil || !os.SameFile(before, after) {
		dev.Close()
		return nil, os.ErrNotExist
	}
	return dev, nil
}

func offlineStatus(err error) string {
	switch {
	case errors.Is(err, os.ErrPermission):
		return needAccess
	case errors.Is(err, os.ErrDeadlineExceeded):
		return awaitSignal
	case errors.Is(err, os.ErrNotExist), errors.Is(err, io.EOF), errors.Is(err, syscall.ENODEV), errors.Is(err, syscall.ENXIO), errors.Is(err, syscall.EIO):
		return awaitDevice
	default:
		return readFailed
	}
}

type consoleReader interface {
	io.ReadCloser
	SetReadDeadline(time.Time) error
}

func watchLayers(ctx context.Context, send func(layerEvent)) {
	watchConsoles(ctx, func() (consoleReader, error) {
		return openConsole("/sys/class/hidraw", "/dev")
	}, time.Second, 2*time.Second, send)
}

func watchConsoles(ctx context.Context, open func() (consoleReader, error), retry, silence time.Duration, send func(layerEvent)) {
	last := layerEvent{Layer: -2}
	report := func(event layerEvent) {
		if event != last && ctx.Err() == nil {
			last = event
			send(event)
		}
	}
	for ctx.Err() == nil {
		dev, err := open()
		if err == nil {
			report(layerEvent{Layer: -1, Status: awaitSignal})
			err = readConsole(ctx, dev, silence, func(layer int) {
				report(layerEvent{Layer: layer})
			})
		}
		report(layerEvent{Layer: -1, Status: offlineStatus(err)})
		timer := time.NewTimer(retry)
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
		}
	}
}

func readConsole(ctx context.Context, dev consoleReader, silence time.Duration, emit func(int)) error {
	// O_NONBLOCK lets Go park the reader in epoll. Closing interrupts even a
	// silent firmware immediately; the timeout alone is for heartbeat loss.
	stop := context.AfterFunc(ctx, func() { dev.Close() })
	defer stop()
	defer dev.Close()
	var decoder layerDecoder
	var report [64]byte // Full-speed USB interrupt packet.
	err := dev.SetReadDeadline(time.Now().Add(silence))
	for err == nil {
		var n int
		n, err = dev.Read(report[:])
		if decoder.feed(report[:n], emit) && err == nil {
			err = dev.SetReadDeadline(time.Now().Add(silence))
		}
	}
	return err
}
