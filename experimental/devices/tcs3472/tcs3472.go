// Copyright 2016 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package tcs3472

import (
	"fmt"

	"periph.io/x/periph/conn"
	"periph.io/x/periph/conn/i2c"
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

	// Commands for reading the light levels.
	cmdReadClear = regCmd | cmdAutoIncrement | 0x14
	cmdReadRed   = regCmd | cmdAutoIncrement | 0x16
	cmdReadGreen = regCmd | cmdAutoIncrement | 0x18
	cmdReadBlue  = regCmd | cmdAutoIncrement | 0x1a

	// Commands for
	cmdEnable  = regCmd | 0
	cmdATime   = regCmd | 1
	cmdControl = regCmd | 0x0f
	cmdChipID  = regCmd | 0x12
	cmdStatus  = regCmd | 0x13

	enableInterrupt = 1 << 4
	enableWait      = 1 << 3
	enableRGBC      = 1 << 1
	enablePower     = 1
)

// Dev is a handle to a TCS34725 Color Sensor.
type Dev struct {
	c conn.Conn
}

// New opens a handle that communicates over I²C with the TCS34725 Color Sensor.
// The bus supports fast mode with a data rate of up to 400 kbit/s.
func New(i i2c.Bus) (*Dev, error) {
	d := &Dev{
		c: &i2c.Dev{Bus: i, Addr: chipAddr},
	}

	var id [1]byte
	if err := d.c.Tx([]byte{cmdChipID}, id[:]); err != nil {
		return nil, err
	}
	if id[0] != chipID && id[0] != chipIDAlt {
		return nil, fmt.Errorf("tcs3472: unexpected chip ID 0x%x", id)
	}

	// self.i2c_bus.write_byte_data(ADDR, REG_ENABLE, REG_ENABLE_RGBC|REG_ENABLE_POWER)
	// self.set_integration_time_ms(511.2)

	// FIXME: Simulate a setup dance.
	// var b [2]byte
	// if err := d.c.Tx([]byte(), b[:]); err != nil {
	// 	return nil, err
	// }
	// if b[0] != 'I' || b[1] != 'N' {
	// 	return nil, errors.New("driverskeleton: unexpected reply")
	// }

	return d, nil
}

// String implements the String method of the fmt.Stringer interface.
func (d *Dev) String() string {
	return fmt.Sprintf("TCS4372{%s}", d.c)
}

// Halt implements the Halt method of the devices.Device interface.
func (d *Dev) Halt() error {
	return nil
}

type Light struct {
	Int     uint8
	R, G, B uint8
}

// Measure measures the light intensity and color.
func (d *Dev) Measure(l *Light) error {
	// var b [12]byte
	// if err := d.c.Tx([]byte("what"), b[:]); err != nil {
	// 	return err.Error()
	// }
	// return string(b[:])
	return nil
}
