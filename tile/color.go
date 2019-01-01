package tile

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/fishy/lifxlan"
)

// ColorBoard represents a board of colors.
//
// The zero value returns nil color on every coordinate.
type ColorBoard [][]*lifxlan.Color

// MakeColorBoard creates a ColorBoard with the given size.
func MakeColorBoard(width, height int) ColorBoard {
	cb := make(ColorBoard, width)
	for i := range cb {
		cb[i] = make([]*lifxlan.Color, height)
	}
	return cb
}

// GetColor returns a color at the given coordinate.
//
// If the given coordinate is out of boundary, nil color will be returned.
func (cb ColorBoard) GetColor(x, y int) *lifxlan.Color {
	if x < 0 || x >= len(cb) {
		return nil
	}
	row := cb[x]
	if y < 0 || y >= len(row) {
		return nil
	}
	return row[y]
}

// RawSetTileState64Payload defines the struct to be used for encoding and
// decoding.
//
// https://lan.developer.lifx.com/docs/tile-messages#section-settilestate64-715
type RawSetTileState64Payload struct {
	TileIndex uint8
	Length    uint8
	_         uint8 // reserved
	X         uint8
	Y         uint8
	Width     uint8
	Duration  uint32
	Colors    [64]lifxlan.Color
}

func (td *device) SetColors(
	ctx context.Context,
	conn net.Conn,
	cb ColorBoard,
	duration time.Duration,
	ack bool,
) error {
	select {
	default:
	case <-ctx.Done():
		return ctx.Err()
	}

	if conn == nil {
		newConn, err := td.Dial()
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

	payloads := make([]*RawSetTileState64Payload, len(td.tiles))
	for i := range payloads {
		payloads[i] = &RawSetTileState64Payload{
			TileIndex: td.startIndex + uint8(i),
			Length:    1,
			Width:     td.tiles[i].Width,
			Duration:  uint32(duration / time.Millisecond),
		}
		// Init with all black colors.
		for j := range payloads[i].Colors {
			payloads[i].Colors[j] = lifxlan.ColorBlack
		}
	}

	for x := 0; x < td.Width(); x++ {
		for y := 0; y < td.Height(); y++ {
			if c := cb.GetColor(x, y); c != nil {
				data := td.board.Data[x][y]
				if data == nil {
					// Not on tile
					continue
				}
				index := data.X*int(td.tiles[data.Index].Width) + data.Y
				payloads[data.Index].Colors[index] = *c
			}
		}
	}

	var flags lifxlan.AckResFlag
	if ack {
		flags |= lifxlan.FlagAckRequired
	}

	var wg sync.WaitGroup
	wg.Add(len(payloads))
	errChan := make(chan error, len(payloads))
	sentChan := make(chan uint8, len(payloads))
	for _, payload := range payloads {
		sequence := td.NextSequence()
		go func(sequence uint8, payload *RawSetTileState64Payload) {
			defer wg.Done()
			buf := new(bytes.Buffer)
			if err := binary.Write(buf, binary.LittleEndian, payload); err != nil {
				errChan <- err
				return
			}
			msg, err := lifxlan.GenerateMessage(
				lifxlan.NotTagged,
				td.Source(),
				td.Target(),
				flags,
				sequence,
				SetTileState64,
				buf.Bytes(),
			)
			if err != nil {
				errChan <- err
				return
			}

			n, err := conn.Write(msg)
			if err != nil {
				errChan <- err
				return
			}
			if n < len(msg) {
				errChan <- fmt.Errorf(
					"lifxlan/tile.SetColors: only wrote %d out of %d bytes",
					n,
					len(msg),
				)
				return
			}
			sentChan <- sequence
		}(sequence, payload)
	}
	wg.Wait()

	seqs := make([]uint8, 0, 0)
	if err := func() error {
		var n int
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err := <-errChan:
				return err
			case seq := <-sentChan:
				n++
				seqs = append(seqs, seq)
				if n >= len(payloads) {
					// All API calls successfully sent.
					return nil
				}
			}
		}
	}(); err != nil {
		return err
	}
	if !ack {
		return nil
	}

	return lifxlan.WaitForAcks(ctx, conn, td, seqs...)
}
