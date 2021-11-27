package tile_test

import (
	"context"
	"log"
	"time"

	"go.yhsif.com/lifxlan/tile"
)

// This example demonstrates how to draw a single frame on a tile device.
func Example_draw() {
	// Need proper initialization on real code.
	var (
		device tile.Device
		// Important to set timeout to context when requiring ack.
		timeout time.Duration
		// ColorBoard to draw.
		cb tile.ColorBoard
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := device.SetColors(
		ctx,
		nil, // conn, use nil so that SetColors will maintain it for us
		cb,
		0,    // fade in duration
		true, // ack
	); err != nil {
		log.Fatal(err)
	}
}

// This example demonstrates how to draw frames continuously on a tile device.
func Example_drawContinuously() {
	// Need proper initialization on real code.
	var (
		device tile.Device
		// Important to set timeout to context when requiring ack.
		timeout time.Duration
		// Interval between frames.
		interval time.Duration
		// Function to return the next frame.
		nextFrame func() tile.ColorBoard
	)

	conn, err := device.Dial()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for range time.Tick(interval) {
		// Use lambda to make sure defer works as expected.
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			if err := device.SetColors(
				ctx,
				conn,
				nextFrame(),
				0,    // fade in duration
				true, // ack
			); err != nil {
				log.Fatal(err)
			}
		}()
	}
}
