package lifxlan

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"net"
)

// RawStateServicePayload defines the struct to be used for encoding and
// decoding.
//
// https://lan.developer.lifx.com/docs/information-messages#stateservice---packet-3
type RawStateServicePayload struct {
	Service ServiceType
	Port    uint32
}

// Default broadcast host and port.
const (
	DefaultBroadcastHost = "255.255.255.255"
	DefaultBroadcastPort = "56700"
)

// Discover discovers lifx products in the lan.
//
// When broadcastHost is empty (""), DefaultBroadcastHost will be used instead.
// In most cases that should just work.
// But if your network has special settings, you can override it via the arg.
//
// The function will write discovered devices into devices channel.
// It's the caller's responsibility to read from channel timely to avoid
// blocking writing.
// The function is guaranteed to close the channel upon retuning,
// so the caller could just range over the channel for reading, e.g.
//
//     devices := make(chan Device)
//     go func() {
//       if err := Discover(ctx, devices, ""); err != nil {
//         if err != context.DeadlineExceeded {
//           // handle error
//         }
//       }
//     }()
//     for device := range devices {
//       // Do something with device
//     }
//
// The function will only return upon error or when ctx is cancelled.
// It's the caller's responsibility to make sure that the context is cancelled
// (e.g. Use context.WithTimeout).
func Discover(
	ctx context.Context,
	devices chan Device,
	broadcastHost string,
) error {
	defer close(devices)

	if ctx.Err() != nil {
		return ctx.Err()
	}

	msg, err := GenerateMessage(
		Tagged,
		0, // source
		AllDevices,
		0, // flags
		0, // sequence
		GetService,
		nil, // payload
	)
	if err != nil {
		return err
	}

	conn, err := net.ListenPacket("udp", ":"+DefaultBroadcastPort)
	if err != nil {
		return err
	}
	defer conn.Close()

	if broadcastHost == "" {
		broadcastHost = DefaultBroadcastHost
	}
	broadcast, err := net.ResolveUDPAddr(
		"udp",
		net.JoinHostPort(broadcastHost, DefaultBroadcastPort),
	)
	if err != nil {
		return err
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	n, err := conn.WriteTo(msg, broadcast)
	if err != nil {
		return err
	}
	if n < len(msg) {
		return fmt.Errorf(
			"lifxlan.Discover: only wrote %d out of %d bytes",
			n,
			len(msg),
		)
	}

	buf := make([]byte, ResponseReadBufferSize)
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if err := conn.SetReadDeadline(GetReadDeadline()); err != nil {
			return err
		}
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			if CheckTimeoutError(err) {
				continue
			}
			return err
		}

		host, _, err := net.SplitHostPort(addr.String())
		if err != nil {
			return err
		}

		resp, err := ParseResponse(buf[:n])
		if err != nil {
			return err
		}
		if resp.Message != StateService {
			continue
		}

		var d RawStateServicePayload
		r := bytes.NewReader(resp.Payload)
		if err := binary.Read(r, binary.LittleEndian, &d); err != nil {
			return err
		}
		switch d.Service {
		default:
			// Unknown service, ignore.
			continue
		case ServiceUDP:
			devices <- NewDevice(
				net.JoinHostPort(host, fmt.Sprintf("%d", d.Port)),
				d.Service,
				resp.Target,
			)
		}
	}
}
