package main

import "bytes"

const framePrefix = "DLM1:"

type layerEvent struct {
	Layer  int
	Status string
}

// One decoder per connection: a partial line must not cross a disconnect.
type layerDecoder struct {
	line [len(framePrefix) + 1]byte
	used int
}

// feed returns true even for a repeated layer: only valid frames renew the deadline.
func (d *layerDecoder) feed(report []byte, emit func(int)) bool {
	valid := false
	for _, b := range report {
		switch b {
		case 0: // QMK pads console reports, including between fragments.
		case '\n':
			if d.used == len(d.line) && bytes.HasPrefix(d.line[:], []byte(framePrefix)) {
				id := d.line[len(framePrefix)]
				if id >= '0' && id <= '6' {
					emit(int(id - '0'))
					valid = true
				}
			}
			d.used = 0
		default:
			if d.used < len(d.line) {
				d.line[d.used] = b
				d.used++
			} else {
				d.used = len(d.line) + 1 // Discard the entire oversized line until LF.
			}
		}
	}
	return valid
}
