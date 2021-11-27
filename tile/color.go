package tile

import (
	"bytes"
	"context"
	"encoding/binary"
	"net"
	"sync"
	"time"

	"go.yhsif.com/lifxlan"
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
// https://lan.developer.lifx.com/docs/changing-a-device#set64---packet-715
type RawSetTileState64Payload struct {
	TileIndex uint8
	Length    uint8
	_         byte // reserved
	X         uint8
	Y         uint8
	Width     uint8
	Duration  lifxlan.TransitionTime
	Colors    [64]lifxlan.Color
}

func (td *device) SetColors(
	ctx context.Context,
	conn net.Conn,
	cb ColorBoard,
	transition time.Duration,
	ack bool,
) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if conn == nil {
		newConn, err := td.Dial()
		if err != nil {
			return err
		}
		defer newConn.Close()
		conn = newConn

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	payloads := make([]*RawSetTileState64Payload, len(td.tiles))
	sanitizedBlack := td.SanitizeColor(lifxlan.ColorBlack)
	for i := range payloads {
		payloads[i] = &RawSetTileState64Payload{
			TileIndex: td.startIndex + uint8(i),
			Length:    1,
			Width:     td.TileWidth(i),
			Duration:  lifxlan.ConvertDuration(transition),
		}
		// Init with all black colors.
		for j := range payloads[i].Colors {
			payloads[i].Colors[j] = sanitizedBlack
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
				colorIndex := data.X*int(td.TileWidth(data.Index)) + data.Y
				payloads[data.Index].Colors[colorIndex] = td.SanitizeColor(*c)
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
			seq, err := td.Send(
				ctx,
				conn,
				flags,
				SetTileState64,
				payload,
			)
			if err != nil {
				errChan <- err
				return
			}
			sentChan <- seq
		}(payload)
	}
	wg.Wait()

	seqs := make([]uint8, 0)
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
// https://lan.developer.lifx.com/docs/querying-the-device-for-data#get64---packet-707
type RawGetTileState64Payload struct {
	TileIndex uint8
	Length    uint8
	_         byte // reserved
	X         uint8
	Y         uint8
	Width     uint8
}

// RawStateTileState64Payload defines the struct to be used for encoding and
// decoding.
//
// https://lan.developer.lifx.com/docs/information-messages#state64---packet-711
type RawStateTileState64Payload struct {
	TileIndex uint8
	_         byte // reserved
	X         uint8
	Y         uint8
	Width     uint8
	Colors    [64]lifxlan.Color
}

func (td *device) GetColors(
	ctx context.Context,
	conn net.Conn,
) (ColorBoard, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if conn == nil {
		newConn, err := td.Dial()
		if err != nil {
			return nil, err
		}
		defer newConn.Close()
		conn = newConn

		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
	}

	// Send
	seq, err := td.Send(
		ctx,
		conn,
		0, // flags
		GetTileState64,
		&RawGetTileState64Payload{
			TileIndex: td.startIndex,
			Length:    uint8(len(td.tiles)),
			Width:     td.TileWidth(0),
		},
	)
	if err != nil {
		return nil, err
	}

	// Read responses
	received := make([]int, len(td.tiles))
	cb := MakeColorBoard(td.Width(), td.Height())
	for {
		resp, err := lifxlan.ReadNextResponse(ctx, conn)
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
			for y := 0; y < int(tile.Height); y++ {
				// c is the coordinate on the color board.
				c := td.board.ReverseData[ti][x][y]
				cb[c.X][c.Y] = &raw.Colors[x*int(tile.Width)+y]
			}
		}

		n := 0
		for _, rec := range received {
			n += rec
		}
		if n >= len(td.tiles) {
			// Got responses for all tiles.
			return cb, nil
		}
	}
}
