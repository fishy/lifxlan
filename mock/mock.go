package mock

import (
	"bytes"
	"context"
	"encoding/binary"
	"net"
	"sync"
	"testing"

	"github.com/fishy/lifxlan"
	"github.com/fishy/lifxlan/tile"
)

// ListenAddr is the addr to listen on this device.
const ListenAddr = "127.0.0.1:"

// Target is the mocked device target.
const Target lifxlan.Target = 1

// Service is a mocked device listening on localhost.
//
// All service functions require TB to be non-nil, or they will panic.
type Service struct {
	// Testing context
	TB testing.TB

	// When AcksToDrop > 0 and it's supposed to send an ack,
	// the ack won't be send and AcksToDrop will decrease by 1.
	AcksToDrop int

	// Payloads to response.
	RawStateLabelPayload        *lifxlan.RawStateLabelPayload
	RawStateVersionPayload      *lifxlan.RawStateVersionPayload
	RawStateDeviceChainPayload  *tile.RawStateDeviceChainPayload
	RawStateTileState64Payloads []*tile.RawStateTileState64Payload

	// The service context.
	//
	// Please note that it's different from the context of the API calls.
	Context context.Context
	Cancel  context.CancelFunc

	wg      sync.WaitGroup
	started bool
}

// StartService starts a mock service, returns the service and the device.
func StartService(tb testing.TB) (*Service, lifxlan.Device) {
	tb.Helper()

	s := &Service{
		TB: tb,
	}
	return s, s.Start()
}

// Start starts the service and returns the device.
func (s *Service) Start() lifxlan.Device {
	s.TB.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	s.Context = ctx
	s.Cancel = cancel

	conn, err := net.ListenPacket("udp", ListenAddr)
	if err != nil {
		s.TB.Fatal(err)
	}

	s.wg.Add(1)
	go s.handler(conn)

	return lifxlan.NewDevice(
		conn.LocalAddr().String(),
		lifxlan.ServiceUDP,
		Target,
	)
}

// Stop stops the mocked device service.
//
// It won't response to any requests after stopped.
func (s *Service) Stop() {
	s.TB.Helper()

	s.Cancel()
	s.wg.Wait()
}

// Reply replies a request.
func (s *Service) Reply(
	conn net.PacketConn,
	addr net.Addr,
	orig *lifxlan.Response,
	message lifxlan.MessageType,
	payload []byte,
) {
	select {
	default:
	case <-s.Context.Done():
		return
	}

	msg, err := lifxlan.GenerateMessage(
		lifxlan.NotTagged,
		orig.Source,
		Target,
		orig.Flags,
		orig.Sequence,
		message,
		payload,
	)
	if err != nil {
		s.TB.Log(err)
		return
	}

	select {
	default:
	case <-s.Context.Done():
		return
	}

	n, err := conn.WriteTo(msg, addr)
	if err != nil {
		s.TB.Log(err)
		return
	}
	if n < len(msg) {
		s.TB.Logf(
			"lifxlan/mock.Reply: only wrote %d out of %d bytes",
			n,
			len(msg),
		)
	}
	return
}

func (s *Service) handler(conn net.PacketConn) {
	defer s.wg.Done()
	defer conn.Close()

	buf := make([]byte, lifxlan.ResponseReadBufferSize)
	for {
		select {
		default:
		case <-s.Context.Done():
			s.TB.Log(s.Context.Err())
			return
		}

		if err := conn.SetReadDeadline(lifxlan.GetReadDeadline()); err != nil {
			s.TB.Log(err)
			continue
		}
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			if !lifxlan.CheckTimeoutError(err) {
				s.TB.Log(err)
			}
			continue
		}

		resp, err := lifxlan.ParseResponse(buf[:n])
		if err != nil {
			s.TB.Log(err)
			continue
		}

		if !resp.Target.Matches(Target) {
			s.TB.Logf("Ignoring unmatched target %v", resp.Target)
			continue
		}

		if resp.Flags|lifxlan.FlagAckRequired != 0 {
			if s.AcksToDrop > 0 {
				s.AcksToDrop--
			} else {
				s.Reply(conn, addr, resp, lifxlan.Acknowledgement, nil)
			}
		}

		switch resp.Message {
		default:
			s.TB.Logf("Ignoring unknown message %v", resp.Message)
			continue

		case lifxlan.GetLabel:
			buf := new(bytes.Buffer)
			if err := binary.Write(
				buf,
				binary.LittleEndian,
				s.RawStateLabelPayload,
			); err != nil {
				s.TB.Log(err)
				continue
			}
			s.Reply(conn, addr, resp, lifxlan.StateLabel, buf.Bytes())

		case lifxlan.GetVersion:
			buf := new(bytes.Buffer)
			if err := binary.Write(
				buf,
				binary.LittleEndian,
				s.RawStateVersionPayload,
			); err != nil {
				s.TB.Log(err)
				continue
			}
			s.Reply(conn, addr, resp, lifxlan.StateVersion, buf.Bytes())

		case tile.GetDeviceChain:
			buf := new(bytes.Buffer)
			if err := binary.Write(
				buf,
				binary.LittleEndian,
				s.RawStateDeviceChainPayload,
			); err != nil {
				s.TB.Log(err)
				continue
			}
			s.Reply(conn, addr, resp, tile.StateDeviceChain, buf.Bytes())

		case tile.GetTileState64:
			for _, payload := range s.RawStateTileState64Payloads {
				buf := new(bytes.Buffer)
				if err := binary.Write(
					buf,
					binary.LittleEndian,
					payload,
				); err != nil {
					s.TB.Log(err)
					continue
				}
				s.Reply(conn, addr, resp, tile.StateTileState64, buf.Bytes())
			}
		}
	}
}
