package lifxlan

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net"
)

// Response is the parsed response from a lifxlan device.
type Response struct {
	Message  MessageType
	Flags    AckResFlag
	Source   uint32
	Target   Target
	Sequence uint8
	Payload  []byte
}

// ParseResponse parses the response received from a lifxlan device.
func ParseResponse(msg []byte) (*Response, error) {
	if len(msg) < int(HeaderLength) {
		return nil, fmt.Errorf(
			"lifxlan.ParseResponse: response size not enough: %d < %d",
			len(msg),
			HeaderLength,
		)
	}

	var d RawHeader
	r := bytes.NewReader(msg)
	if err := binary.Read(r, binary.LittleEndian, &d); err != nil {
		return nil, err
	}
	if len(msg) != int(d.Size) {
		return nil, fmt.Errorf(
			"lifxlan.ParseResponse: response size mismatch: %d != %d",
			len(msg),
			d.Size,
		)
	}

	payload, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return &Response{
		Message:  d.Type,
		Flags:    d.Flags,
		Source:   d.Source,
		Target:   d.Target,
		Sequence: d.Sequence,
		Payload:  payload,
	}, nil
}

// ReadNextResponse returns the next received response.
//
// It handles read buffer, deadline, context cancellation check,
// and response parsing.
func ReadNextResponse(ctx context.Context, conn net.Conn) (*Response, error) {
	buf := make([]byte, ResponseReadBufferSize)
	for {
		select {
		default:
		case <-ctx.Done():
			return nil, ctx.Err()
		}

		if err := conn.SetReadDeadline(GetReadDeadline()); err != nil {
			return nil, err
		}

		n, err := conn.Read(buf)
		if err != nil {
			if CheckTimeoutError(err) {
				continue
			}
			return nil, err
		}

		return ParseResponse(buf[:n])
	}
}
