package lifxlan

import (
	"context"
	"net"
)

// WaitForAcks helps device API implementations to wait for acks.
//
// It blocks until acks for all sequences are received,
// in which case it returns nil error.
// Or until the context is cancelled.
//
// There shouldn't be more than one WaitForAcks functions running for the same
// connection at the same time.
// It might cause some sequences waited in one goroutine being dropped by other
// goroutines.
func WaitForAcks(
	ctx context.Context,
	conn net.Conn,
	d Device,
	sequences ...uint8,
) error {
	seqMap := make(map[uint8]bool)
	for _, seq := range sequences {
		seqMap[seq] = true
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
		if resp.Source != d.Source() || resp.Message != Acknowledgement {
			continue
		}
		if seqMap[resp.Sequence] {
			delete(seqMap, resp.Sequence)
			if len(seqMap) == 0 {
				// All ack received.
				return nil
			}
		}
	}
}
