package tile

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/fishy/lifxlan"
)

// Device is a wrapped lifxlan.Device that provides tile related APIs.
type Device interface {
	lifxlan.Device

	Board

	// Tiles returns a copy of the tiles in this device.
	Tiles() []Tile

	// SetColors sets the tile device with the given color board.
	//
	// If conn is nil,
	// a new connection will be made and guaranteed to be closed before returning.
	// You should pre-dial and pass in the conn if you plan to call this function
	// repeatedly.
	//
	// duration is the fade in time duration.
	//
	// If ack is false,
	// the function returns nil error after the API is sent successfully.
	// If ack is true,
	// the function will only return nil error after it received ack from the
	// device.
	SetColors(ctx context.Context, conn net.Conn, cb ColorBoard, duration time.Duration, ack bool) error
}

type device struct {
	lifxlan.Device

	startIndex uint8
	tiles      []*Tile

	// parsed board data
	board BoardData
}

var _ Device = (*device)(nil)

func (td *device) String() string {
	return fmt.Sprintf("TileDevice(%v)", td.Target())
}

func (td *device) Tiles() []Tile {
	tiles := make([]Tile, len(td.tiles))
	for i := range tiles {
		tiles[i] = *td.tiles[i]
	}
	return tiles
}
