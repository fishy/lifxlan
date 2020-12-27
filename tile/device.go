package tile

import (
	"context"
	"fmt"
	"net"
	"time"

	"go.yhsif.com/lifxlan"
	"go.yhsif.com/lifxlan/light"
)

// Device is a wrapped lifxlan.Device that provides tile related APIs.
type Device interface {
	light.Device

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

	// TileWidth returns the width of the i-th tile.
	//
	// If i is out of bound, it returns the width of the first tile (index 0)
	// instead. If there's no known tiles, it returns 0.
	TileWidth(i int) uint8
}

type device struct {
	light.Device

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
	if parsed := td.HardwareVersion().Parse(); parsed != nil {
		return fmt.Sprintf("%s(%v)", parsed.ProductName, td.Target())
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

func (td *device) TileWidth(i int) uint8 {
	if len(td.tiles) == 0 {
		return 0
	}
	if i < 0 || i >= len(td.tiles) {
		i = 0
	}
	return td.tiles[i].Width
}
