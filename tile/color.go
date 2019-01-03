package tile

import (
	"bytes"
	"context"
	"encoding/binary"
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

// TileState64Width is the width value to be used in *TileState64 messages.
//
// Please note that in most cases we try to avoid hardcoding the tile width and
// height, and use the values returned by StateDeviceChain message instead.
// But for *TileState64 messages it only makes sense to hardcode them,
// as the colors array is hardcoded as size of 64.
const TileState64Width = 8

// GetColor returns the color at the given coordinate.
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
// https://lan.developer.lifx.com/v2.0/docs/tile-messages#section-settilestate64-715
type RawSetTileState64Payload struct {
	TileIndex uint8
	Length    uint8
	_         uint8 // reserved
	X         uint8
	Y         uint8
	Width     uint8
	Duration  uint32
	Colors    [TileState64Width * TileState64Width]lifxlan.Color
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
			Width:     TileState64Width,
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
				index := data.X*TileState64Width + data.Y
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
		go func(payload *RawSetTileState64Payload) {
			defer wg.Done()
			buf := new(bytes.Buffer)
			if err := binary.Write(buf, binary.LittleEndian, payload); err != nil {
				errChan <- err
				return
			}
			seq, err := td.Send(
				ctx,
				conn,
				lifxlan.NotTagged,
				flags,
				SetTileState64,
				buf.Bytes(),
			)
			if err != nil {
				errChan <- err
				return
			}
			sentChan <- seq
		}(payload)
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

	if ack {
		return lifxlan.WaitForAcks(ctx, conn, td.Source(), seqs...)
	}
	return nil
}

// RawGetTileState64Payload defines the struct to be used for encoding and
// decoding.
//
// https://lan.developer.lifx.com/v2.0/docs/tile-messages#section-gettilestate64-707
type RawGetTileState64Payload struct {
	TileIndex uint8
	Length    uint8
	_         uint8 // reserved
	X         uint8
	Y         uint8
	Width     uint8
}

// RawStateTileState64Payload defines the struct to be used for encoding and
// decoding.
//
// https://lan.developer.lifx.com/v2.0/docs/tile-messages#section-statetilestate64-711
type RawStateTileState64Payload struct {
	TileIndex uint8
	_         uint8 // reserved
	X         uint8
	Y         uint8
	Width     uint8
	Colors    [TileState64Width * TileState64Width]lifxlan.Color
}

func (td *device) GetColors(
	ctx context.Context,
	conn net.Conn,
) (ColorBoard, error) {
	select {
	default:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	if conn == nil {
		newConn, err := td.Dial()
		if err != nil {
			return nil, err
		}
		defer newConn.Close()
		conn = newConn

		select {
		default:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	// Send
	payload := &RawGetTileState64Payload{
		TileIndex: td.startIndex,
		Length:    uint8(len(td.tiles)),
		Width:     TileState64Width,
	}
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, payload); err != nil {
		return nil, err
	}
	seq, err := td.Send(
		ctx,
		conn,
		lifxlan.NotTagged,
		0, // flags
		GetTileState64,
		buf.Bytes(),
	)
	if err != nil {
		return nil, err
	}

	// Read responses
	respBuf := make([]byte, lifxlan.ResponseReadBufferSize)
	received := make([]int, len(td.tiles))
	cb := MakeColorBoard(td.Width(), td.Height())
	for {
		select {
		default:
		case <-ctx.Done():
			return nil, ctx.Err()
		}

		if err := conn.SetReadDeadline(lifxlan.GetReadDeadline()); err != nil {
			return nil, err
		}

		n, err := conn.Read(respBuf)
		if err != nil {
			if lifxlan.CheckTimeoutError(err) {
				continue
			}
			return nil, err
		}

		resp, err := lifxlan.ParseResponse(respBuf[:n])
		if err != nil {
			return nil, err
		}
		if resp.Sequence != seq || resp.Source != td.Source() {
			continue
		}
		if resp.Message != StateTileState64 {
			continue
		}

		var raw RawStateTileState64Payload
		r := bytes.NewReader(resp.Payload)
		if err := binary.Read(r, binary.LittleEndian, &raw); err != nil {
			return nil, err
		}

		// tile index
		ti := raw.TileIndex - td.startIndex
		received[ti] = 1
		tile := td.tiles[ti]
		for x := 0; x < int(tile.Width); x++ {
			if x >= TileState64Width {
				continue
			}
			for y := 0; y < int(tile.Height); y++ {
				if y >= TileState64Width {
					continue
				}
				// index is the index on the returned colors array.
				index := x*TileState64Width + y
				// c is the coordinate on the color board.
				c := td.board.ReverseData[ti][x][y]
				cb[c.X][c.Y] = &raw.Colors[index]
			}
		}

		n = 0
		for _, rec := range received {
			n += rec
		}
		if n >= len(td.tiles) {
			// Got responses for all tiles.
			return cb, nil
		}
	}
}
