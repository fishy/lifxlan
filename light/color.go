package light

import (
	"bytes"
	"context"
	"encoding/binary"
	"net"
	"time"

	"github.com/fishy/lifxlan"
)

// RawSetColorPayload defines the struct to be used for encoding and decoding.
//
// https://lan.developer.lifx.com/v2.0/docs/light-messages#section-setcolor-102
type RawSetColorPayload struct {
	_        uint8 // reserved
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
	select {
	default:
	case <-ctx.Done():
		return ctx.Err()
	}

	if conn == nil {
		newConn, err := ld.Dial()
		if err != nil {
			return err
		}
		defer newConn.Close()
		conn = newConn

		select {
		default:
		case <-ctx.Done():
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
// https://lan.developer.lifx.com/v2.0/docs/light-messages#section-state-107
type RawStatePayload struct {
	Color lifxlan.Color
	_     int16 // reserved
	Power lifxlan.Power
	Label lifxlan.Label
	_     uint64 // reserved
}

func (ld *device) GetColor(
	ctx context.Context,
	conn net.Conn,
) (*lifxlan.Color, error) {
	select {
	default:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	if conn == nil {
		newConn, err := ld.Dial()
		if err != nil {
			return nil, err
		}
		defer newConn.Close()
		conn = newConn

		select {
		default:
		case <-ctx.Done():
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
	buf := make([]byte, lifxlan.ResponseReadBufferSize)
	for {
		select {
		default:
		case <-ctx.Done():
			return nil, ctx.Err()
		}

		if err := conn.SetReadDeadline(lifxlan.GetReadDeadline()); err != nil {
			return nil, err
		}

		n, err := conn.Read(buf)
		if err != nil {
			if lifxlan.CheckTimeoutError(err) {
				continue
			}
			return nil, err
		}

		resp, err := lifxlan.ParseResponse(buf[:n])
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
