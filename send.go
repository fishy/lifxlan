package lifxlan

import (
	"context"
	"fmt"
	"net"
)

func (d *device) Send(
	ctx context.Context,
	conn net.Conn,
	flags AckResFlag,
	message MessageType,
	payload []byte,
) (seq uint8, err error) {
	select {
	default:
	case <-ctx.Done():
		err = ctx.Err()
		return
	}

	var msg []byte
	seq = d.NextSequence()
	msg, err = GenerateMessage(
		NotTagged,
		d.Source(),
		d.Target(),
		flags,
		seq,
		message,
		payload,
	)
	if err != nil {
		return
	}

	select {
	default:
	case <-ctx.Done():
		err = ctx.Err()
		return
	}

	var n int
	n, err = conn.Write(msg)
	if err != nil {
		return
	}
	if n < len(msg) {
		err = fmt.Errorf(
			"lifxlan.Device.Send: only wrote %d out of %d bytes",
			n,
			len(msg),
		)
		return
	}

	select {
	default:
	case <-ctx.Done():
		err = ctx.Err()
		return
	}

	return
}
