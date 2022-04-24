package relay_test

import (
	"context"
	"log"
	"time"

	"go.yhsif.com/lifxlan"
	"go.yhsif.com/lifxlan/relay"
)

func Example() {
	// Need proper initialization in real code.
	var (
		device relay.Device
		// Important to set timeout to context when requiring ack.
		timeout time.Duration
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := device.SetRPower(
		ctx,
		nil, // conn, use nil so that SetRPower will maintain it for us
		0,   // index of the relay to controal
		lifxlan.PowerOn,
		true, // ack
	); err != nil {
		log.Fatal(err)
	}
}
