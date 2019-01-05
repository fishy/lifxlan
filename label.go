package lifxlan

import (
	"bytes"
	"context"
	"encoding/binary"
	"flag"
	"net"
)

// EmptyLabel is the constant to be compared against Device.Label().String().
const EmptyLabel = ""

func (d *device) Label() *RawLabel {
	return &d.label
}

// RawStateLabelPayload defines the struct to be used for encoding and decoding.
//
// https://lan.developer.lifx.com/v2.0/docs/device-messages#section-statelabel-25
type RawStateLabelPayload struct {
	Label RawLabel
}

// RawLabelLength is the length of the raw label used in messages.
const RawLabelLength = 32

// RawLabel defines raw label in message payloads according to:
//
// https://lan.developer.lifx.com/v2.0/docs/device-messages#section-labels
type RawLabel [RawLabelLength]byte

var _ flag.Value = (*RawLabel)(nil)

func (l RawLabel) String() string {
	index := bytes.IndexByte(l[:], 0)
	if index < 0 {
		return string(l[:])
	}
	return string(l[:index])
}

// Set encodes label into RawLabel.
//
// Long labels will be truncated. This function always return nil error.
//
// It also implements flag.Value interface.
func (l *RawLabel) Set(label string) error {
	for i := 0; i < RawLabelLength; i++ {
		l[i] = 0
	}
	copy((*l)[:], label)
	return nil
}

func (d *device) GetLabel(ctx context.Context, conn net.Conn) error {
	select {
	default:
	case <-ctx.Done():
		return ctx.Err()
	}

	if conn == nil {
		newConn, err := d.Dial()
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

	seq, err := d.Send(
		ctx,
		conn,
		0, // flags
		GetLabel,
		nil, // payload
	)
	if err != nil {
		return err
	}

	buf := make([]byte, ResponseReadBufferSize)
	for {
		select {
		default:
		case <-ctx.Done():
			return ctx.Err()
		}

		if err := conn.SetReadDeadline(GetReadDeadline()); err != nil {
			return err
		}

		n, err := conn.Read(buf)
		if err != nil {
			if CheckTimeoutError(err) {
				continue
			}
			return err
		}

		resp, err := ParseResponse(buf[:n])
		if err != nil {
			return err
		}
		if resp.Sequence != seq || resp.Source != d.Source() {
			continue
		}
		if resp.Message != StateLabel {
			continue
		}

		var raw RawStateLabelPayload
		r := bytes.NewReader(resp.Payload)
		if err := binary.Read(r, binary.LittleEndian, &raw); err != nil {
			return err
		}

		d.label = raw.Label
		return nil
	}
}
