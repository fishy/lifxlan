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

	// GetColors returns the current color board on this tile device.
	//
	// If conn is nil,
	// a new connection will be made and guaranteed to be closed before returning.
	// You should pre-dial and pass in the conn if you plan to call APIs on this
	// device repeatedly.
	//
	// This function will wait for len(Tiles()) response messages.
	// In case of one or more of the responses get dropped on the network,
	// this function will wait until context is cancelled.
	// So it's important to set an appropriate timeout on the context.
	GetColors(ctx context.Context, conn net.Conn) (ColorBoard, error)

	// SetColors sets the tile device with the given color board.
	//
	// If conn is nil,
	// a new connection will be made and guaranteed to be closed before returning.
	// You should pre-dial and pass in the conn if you plan to call APIs on this
	// device repeatedly.
	//
	// If ack is false,
	// this function returns nil error after the API is sent successfully.
	// If ack is true,
	// this function will only return nil error after it received all ack(s) from
	// the device.
	SetColors(ctx context.Context, conn net.Conn, cb ColorBoard, transition time.Duration, ack bool) error
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
	if label := td.Label().String(); label != lifxlan.EmptyLabel {
		return fmt.Sprintf("%s(%v)", label, td.Target())
	}
	return fmt.Sprintf("TileDevice(%v)", td.Target())
}

func (td *device) Tiles() []Tile {
	tiles := make([]Tile, len(td.tiles))
	for i := range tiles {
		tiles[i] = *td.tiles[i]
	}
	return tiles
}
