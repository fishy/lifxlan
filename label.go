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

func (d *device) Label() *Label {
	return &d.label
}

// RawStateLabelPayload defines the struct to be used for encoding and decoding.
//
// https://lan.developer.lifx.com/docs/information-messages#statelabel---packet-25
type RawStateLabelPayload struct {
	Label Label
}

// LabelLength is the length of the raw label used in messages.
const LabelLength = 32

// Label defines raw label in message payloads according to:
//
// https://lan.developer.lifx.com/docs/information-messages#statelabel---packet-25
type Label [LabelLength]byte

var _ flag.Getter = (*Label)(nil)

func (l Label) String() string {
	index := bytes.IndexByte(l[:], 0)
	if index < 0 {
		return string(l[:])
	}
	return string(l[:index])
}

// Set encodes label into Label.
//
// Long labels will be truncated. This function always return nil error.
//
// It also implements flag.Value interface.
func (l *Label) Set(label string) error {
	for i := 0; i < LabelLength; i++ {
		l[i] = 0
	}
	copy((*l)[:], label)
	return nil
}

// Get implements flag.Getter interface.
func (l Label) Get() interface{} {
	return l
}

func (d *device) GetLabel(ctx context.Context, conn net.Conn) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if conn == nil {
		newConn, err := d.Dial()
		if err != nil {
			return err
		}
		defer newConn.Close()
		conn = newConn

		if ctx.Err() != nil {
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

	for {
		resp, err := ReadNextResponse(ctx, conn)
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
