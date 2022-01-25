package lifxlan

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"math/rand"
	"net"
)

const (
	EchoPayloadLength = 64
)

// RawEchoResponsePayload defines echo response payload according to:
//
// https://lan.developer.lifx.com/docs/information-messages#echoresponse---packet-59
type RawEchoResponsePayload struct {
	Echoing [EchoPayloadLength]byte
}

func (d *device) Echo(ctx context.Context, conn net.Conn, payload []byte) error {
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

	body := make([]byte, EchoPayloadLength)
	copy(body, payload)
	if len(payload) < EchoPayloadLength {
		rand.Read(body[len(payload):])
	}

	seq, err := d.Send(
		ctx,
		conn,
		0, // flags
		EchoRequest,
		body,
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
		if resp.Message != EchoResponse {
			continue
		}

		var raw RawEchoResponsePayload
		r := bytes.NewReader(resp.Payload)
		if err := binary.Read(r, binary.LittleEndian, &raw); err != nil {
			return err
		}

		if !bytes.Equal(raw.Echoing[:], body) {
			return errors.New("unexpected echo response value")
		}

		return nil
	}
}
