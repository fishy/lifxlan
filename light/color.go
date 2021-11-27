package light

import (
	"bytes"
	"context"
	"encoding/binary"
	"net"
	"time"

	"go.yhsif.com/lifxlan"
)

// RawSetColorPayload defines the struct to be used for encoding and decoding.
//
// https://lan.developer.lifx.com/docs/changing-a-device#setcolor---packet-102
type RawSetColorPayload struct {
	_        byte // reserved
	Color    lifxlan.Color
	Duration lifxlan.TransitionTime
}

func (ld *device) SetColor(
	ctx context.Context,
	conn net.Conn,
	color *lifxlan.Color,
	transition time.Duration,
	ack bool,
) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if conn == nil {
		newConn, err := ld.Dial()
		if err != nil {
			return err
		}
		defer newConn.Close()
		conn = newConn

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	var flags lifxlan.AckResFlag
	if ack {
		flags |= lifxlan.FlagAckRequired
	}

	// Send
	seq, err := ld.Send(
		ctx,
		conn,
		flags,
		SetColor,
		&RawSetColorPayload{
			Color:    ld.SanitizeColor(*color),
			Duration: lifxlan.ConvertDuration(transition),
		},
	)
	if err != nil {
		return err
	}

	if ack {
		return lifxlan.WaitForAcks(ctx, conn, ld.Source(), seq)
	}
	return nil
}

// RawStatePayload defines the struct to be used for encoding and decoding.
//
// https://lan.developer.lifx.com/docs/information-messages#lightstate---packet-107
type RawStatePayload struct {
	Color lifxlan.Color
	_     [2]byte // reserved
	Power lifxlan.Power
	Label lifxlan.Label
	_     [8]byte // reserved
}

func (ld *device) GetColor(
	ctx context.Context,
	conn net.Conn,
) (*lifxlan.Color, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if conn == nil {
		newConn, err := ld.Dial()
		if err != nil {
			return nil, err
		}
		defer newConn.Close()
		conn = newConn

		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
	}

	// Send
	seq, err := ld.Send(
		ctx,
		conn,
		0, // flags
		Get,
		nil, // payload
	)
	if err != nil {
		return nil, err
	}

	// Read
	for {
		resp, err := lifxlan.ReadNextResponse(ctx, conn)
		if err != nil {
			return nil, err
		}
		if resp.Sequence != seq || resp.Source != ld.Source() {
			continue
		}
		if resp.Message != State {
			continue
		}

		var raw RawStatePayload
		r := bytes.NewReader(resp.Payload)
		if err := binary.Read(r, binary.LittleEndian, &raw); err != nil {
			return nil, err
		}

		*ld.Label() = raw.Label
		// Make a copy so we don't pin the whole raw payload from gc.
		color := raw.Color
		return &color, nil
	}
}
