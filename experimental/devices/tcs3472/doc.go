// Copyright 2016 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package tcs3472 controls a AMS TCS34725 Color Sensor over I²C.
//
// The TCS3472 device provides a digital return of red, green, blue (RGB), and
// clear light sensing values. An IR blocking filter, integrated on-chip and
// localized to the color sensing photodiodes, minimizes the IR spectral component
// of the incoming light and allows color measurements to be made accurately.
//
// Datasheet
//
// http://ams.com/eng/content/download/319364/1117183/file/TCS3472_Datasheet_EN_v2.pdf
package tcs3472
