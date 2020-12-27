package light_test

import (
	"context"
	"log"
	"time"

	"go.yhsif.com/lifxlan"
	"go.yhsif.com/lifxlan/light"
)

// This example demonstrates how to set color on a light device.
func Example() {
	// Need proper initialization on real code.
	var (
		device light.Device
		// Important to set timeout to context when requiring ack.
		timeout time.Duration
		// Color to set.
		color lifxlan.Color
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := device.SetColor(
		ctx,
		nil, // conn, use nil so that SetColors will maintain it for us
		&color,
		0,    // fade in duration
		true, // ack
	); err != nil {
		log.Fatal(err)
	}
}
