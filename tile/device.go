package tile

import (
	"fmt"
	"net"

	"github.com/fishy/lifxlan"
)

// Device is a wrapped lifxlan.Device that provides tile related APIs.
type Device interface {
	lifxlan.Device

	Board

	// Tiles returns a copy of the tiles in this device.
	Tiles() []Tile
}

type device struct {
	dev        lifxlan.Device
	startIndex uint8
	tiles      []*Tile

	// parsed board data
	board BoardData
}

var _ Device = (*device)(nil)

func (td *device) String() string {
	return fmt.Sprintf("TileDevice(%v)", td.Target())
}

func (td *device) Target() lifxlan.Target {
	return td.dev.Target()
}

func (td *device) Dial() (net.Conn, error) {
	return td.dev.Dial()
}

func (td *device) Source() uint32 {
	return td.dev.Source()
}

func (td *device) NextSequence() uint8 {
	return td.dev.NextSequence()
}

func (td *device) Tiles() []Tile {
	tiles := make([]Tile, len(td.tiles))
	for i := range tiles {
		tiles[i] = *td.tiles[i]
	}
	return tiles
}
