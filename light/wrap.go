package light

import (
	"bytes"
	"context"
	"encoding/binary"

	"github.com/fishy/lifxlan"
)

// Wrap tries to wrap a lifxlan.Device into a light device.
//
// When force is false and d is already a light device,
// d will be casted and returned directly.
// Otherwise, this function calls a light device API,
// and only returns a non-nil Device if it supports the API.
//
// If the device is not a light device,
// the function might block until ctx is cancelled.
//
// When returning a valid light device,
// the device's Label is guaranteed to be cached.
func Wrap(ctx context.Context, d lifxlan.Device, force bool) (Device, error) {
	select {
	default:
	case <-ctx.Done():
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

	select {
	default:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	const msg = Get

	seq, err := d.Send(
		ctx,
		conn,
		0, // flags
		msg,
		nil, // payload
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
		case State:
			var raw RawStatePayload
			r := bytes.NewReader(resp.Payload)
			if err := binary.Read(r, binary.LittleEndian, &raw); err != nil {
				return nil, err
			}

			ld := &device{
				Device: d,
			}
			*ld.Label() = raw.Label
			return ld, nil

		case lifxlan.StateUnhandled:
			var raw lifxlan.RawStateUnhandledPayload
			r := bytes.NewReader(resp.Payload)
			if err := binary.Read(r, binary.LittleEndian, &raw); err != nil {
				return nil, err
			}

			return nil, raw.GenerateError(msg)
		}
	}
}
