package relay

import (
	"bytes"
	"context"
	"encoding/binary"

	"go.yhsif.com/lifxlan"
)

// Wrap tries to wrap a lifxlan.Device into a relay device.
//
// When force is false and d is already a relay device,
// d will be casted and returned directly.
// Otherwise, this function calls a relay device API,
// and only returns a non-nil Device if it supports the API.
//
// If the device is not a relay device,
// the function might block until ctx is cancelled.
func Wrap(ctx context.Context, d lifxlan.Device, force bool) (Device, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if !force {
		if t, ok := d.(Device); ok {
			return t, nil
		}
	}

	conn, err := d.Dial()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	seq, err := d.Send(
		ctx,
		conn,
		0, // flags
		GetRPower,
		&RawGetRPowerPayload{
			Index: 0,
		},
	)
	if err != nil {
		return nil, err
	}

	for {
		resp, err := lifxlan.ReadNextResponse(ctx, conn)
		if err != nil {
			return nil, err
		}
		if resp.Sequence != seq || resp.Source != d.Source() {
			continue
		}

		switch resp.Message {
		case StateRPower:
			var raw RawStateRPowerPayload
			r := bytes.NewReader(resp.Payload)
			if err := binary.Read(r, binary.LittleEndian, &raw); err != nil {
				return nil, err
			}

			return &device{
				Device: d,
			}, nil

		case lifxlan.StateUnhandled:
			var raw lifxlan.RawStateUnhandledPayload
			r := bytes.NewReader(resp.Payload)
			if err := binary.Read(r, binary.LittleEndian, &raw); err != nil {
				return nil, err
			}
			return nil, raw
		}
	}
}
