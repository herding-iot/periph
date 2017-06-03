// Copyright 2017 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// bmp180 reads the current temperature and pressure from a BMP180.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/conn/pin"
	"periph.io/x/periph/conn/pin/pinreg"
	"periph.io/x/periph/experimental/devices/tcs3472"
	"periph.io/x/periph/host"
)

func printPin(fn string, p pin.Pin) {
	name, pos := pinreg.Position(p)
	if name != "" {
		log.Printf("  %-4s: %-10s found on header %s, #%d\n", fn, p, name, pos)
	} else {
		log.Printf("  %-4s: %-10s\n", fn, p)
	}
}

// TODO: Maybe use an interface for sensing light?
func read(d *tcs3472.Dev, interval time.Duration) error {
	var t *time.Ticker
	if interval != 0 {
		t = time.NewTicker(interval)
	}

	if valid, err := d.Valid(); !valid || err != nil {
		return err
	}

	for {
		var light tcs3472.Light
		if err := d.Measure(&light); err != nil {
			return err
		}
		fmt.Printf("Int: %d, RGB: %.3f, %.3f, %.3f\n", light.Int, light.R, light.G, light.B)
		if t == nil {
			break
		}

		<-t.C
	}
	return nil
}

func mainImpl() error {
	i2cID := flag.String("i2c", "", "I²C bus to use")
	interval := flag.Duration("i", 0, "read data continously with this interval")
	verbose := flag.Bool("v", false, "verbose mode")
	flag.Parse()
	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}
	log.SetFlags(log.Lmicroseconds)

	if _, err := host.Init(); err != nil {
		return err
	}

	bus, err := i2creg.Open(*i2cID)
	if err != nil {
		return err
	}
	defer bus.Close()

	if p, ok := bus.(i2c.Pins); ok {
		printPin("SCL", p.SCL())
		printPin("SDA", p.SDA())
	}

	dev, err := tcs3472.New(bus)
	if err != nil {
		return err
	}

	err = read(dev, *interval)
	err2 := dev.Halt()
	if err != nil {
		return err
	}
	return err2
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "tcs3472: %s.\n", err)
		os.Exit(1)
	}
}
