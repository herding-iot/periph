// Copyright 2016 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package tcs3472

import (
	"encoding/binary"
	"errors"
	"fmt"

	"time"

	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/mmr"
)

// FIXME: Expose public symbols as relevant. Do not export more than needed!
// See https://periph.io/x/periph/tree/master/doc/drivers#requirements
// for the expectations.

// Gain is the analog gain of the RGBC meassurement.
type Gain uint8

// Possible gain values.
const (
	G1x  Gain = 0
	G4x  Gain = 1
	G16x Gain = 2
	G60x Gain = 3
)

const (
	chipAddr  = 0x29
	chipID    = 0x44 // TCS34721 and TCS34725
	chipIDAlt = 0x4D // TCS34723 and TCS34727

	regCmd           = 0x80 // 0b10000000
	cmdAutoIncrement = 0x20 // 0b00100000

	// Commands for configuration.
	cmdEnable  = regCmd | 0
	cmdATime   = regCmd | 1
	cmdControl = regCmd | 0x0f
	cmdChipID  = regCmd | 0x12
	cmdStatus  = regCmd | 0x13

	// Commands for reading the light values.
	cmdReadClear = regCmd | cmdAutoIncrement | 0x14
	cmdReadRed   = regCmd | cmdAutoIncrement | 0x16
	cmdReadGreen = regCmd | cmdAutoIncrement | 0x18
	cmdReadBlue  = regCmd | cmdAutoIncrement | 0x1a

	// Flags for the enable command.
	enableInterrupt = 1 << 4
	enableWait      = 1 << 3
	enableRGBC      = 1 << 1
	enablePower     = 1

	// Masks for the status response.
	statusInterrupt = 1 << 4
	statusValid     = 1
)

// Dev is a handle to a TCS34725 Color Sensor.
type Dev struct {
	dev      mmr.Dev8
	maxCount uint32
}

// New opens a handle that communicates over I²C with the TCS34725 Color Sensor.
// The bus supports fast mode with a data rate of up to 400 kbit/s.
func New(b i2c.Bus) (*Dev, error) {
	bus := &i2c.Dev{Bus: b, Addr: chipAddr}
	d := &Dev{
		dev: mmr.Dev8{
			Conn:  bus,
			Order: binary.BigEndian,
		},
	}

	// Check the ID.
	id, err := d.dev.ReadUint8(cmdChipID)
	if err != nil {
		return nil, err
	}
	if id != chipID && id != chipIDAlt {
		return nil, fmt.Errorf("tcs3472: unexpected chip ID 0x%x", id)
	}

	// Enable the device.
	if err := d.dev.WriteUint8(cmdEnable, enablePower|enableRGBC); err != nil {
		return nil, err
	}
	if err := d.SetIntegrationTime(511 * time.Millisecond); err != nil {
		return nil, err
	}

	return d, nil
}

// String implements the String method of the fmt.Stringer interface.
func (d *Dev) String() string {
	return fmt.Sprintf("TCS4372{%s}", d.dev.Conn)
}

// Halt implements the Halt method of the devices.Device interface.
func (d *Dev) Halt() error {
	return nil
}

// SetIntegrationTime sets the integration time of measurements. It can be
// between 2.4 and 612 ms.
func (d *Dev) SetIntegrationTime(dur time.Duration) error {
	// The RGBC timing register controls the internal integration time of the
	// RGBC clear and IR channel ADCs in 2.4-ms increments.
	// Max RGBC Count = (256 − ATIME) × 1024 up to a maximum of 65535
	if dur < 2400*time.Microsecond || dur > 612*time.Millisecond {
		return errors.New("integration time must be between 2.4 and 612 ms")
	}
	atime := 255 - uint8(dur.Nanoseconds()/2400000)
	d.maxCount = uint32(atime) * 1024
	if d.maxCount > 65535 {
		d.maxCount = 65535
	}
	if err := d.dev.WriteUint8(cmdATime, atime); err != nil {
		return err
	}
	return nil
}

// MaxCount returns the max value that can be counted for a channel for the
// choosen integration time. Longer integration times give a larger count.
// The count ranges from 1024 for a 2.4 ms integration time to 65535 for a
// 612 ms integration time.
func (d *Dev) MaxCount() uint32 {
	return d.maxCount
}

func (d *Dev) SetGain(g Gain) error {
	return nil
}

// Light is a light measurement with intensity and color.
type Light struct {
	// The light intensity, up to MaxCount().
	Int uint16
	// The color values as a percentage.
	R, G, B float32
}

// Measure measures the light intensity and color.
func (d *Dev) Measure(l *Light) error {
	var err error
	if l.Int, err = d.dev.ReadUint16(cmdReadClear); err != nil {
		return err
	}
	var r, g, b uint16
	r, err = d.dev.ReadUint16(cmdReadRed)
	if err != nil {
		return err
	}
	g, err = d.dev.ReadUint16(cmdReadGreen)
	if err != nil {
		return err
	}
	b, err = d.dev.ReadUint16(cmdReadBlue)
	if err != nil {
		return err
	}
	// Normalize the color values to the intensity.
	l.R = float32(r) / float32(l.Int)
	l.G = float32(g) / float32(l.Int)
	l.B = float32(b) / float32(l.Int)
	return nil
}

// Valid indicates that the RGBC channels have completed an integration cycle.
func (d *Dev) Valid() (bool, error) {
	status, err := d.dev.ReadUint8(cmdStatus)
	if err != nil {
		return false, err
	}

	return status&statusValid > 0, nil
}
