package lifxlan

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
)

// StateService message payload offsets.
const (
	StateServiceServiceOffset = 0                             // 1 byte
	StateServicePortOffset    = StateServiceServiceOffset + 1 // 4 bytes
)

// Default boardcast host and port.
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

	select {
	default:
	case <-ctx.Done():
		return ctx.Err()
	}

	msg := GenerateMessage(
		Tagged,
		0, // source
		AllDevices,
		0, // flags
		0, // sequence
		GetService,
		nil, // payload
	)

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

	select {
	default:
	case <-ctx.Done():
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
		select {
		default:
		case <-ctx.Done():
			return ctx.Err()
		}

		if err := conn.SetReadDeadline(getReadDeadline()); err != nil {
			return err
		}
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			if timeoutErr, ok := err.(timeouter); ok {
				if timeoutErr.Timeout() {
					continue
				}
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

		service := ServiceType(resp.Payload[StateServiceServiceOffset])
		switch service {
		default:
			// Unkown service, ignore.
			continue
		case ServiceUDP:
			port := binary.LittleEndian.Uint32(resp.Payload[StateServicePortOffset:])
			devices <- NewDevice(
				net.JoinHostPort(host, fmt.Sprintf("%d", port)),
				service,
				resp.Target,
			)
		}
	}
}
