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
}

type device struct {
	dev        lifxlan.Device
	startIndex uint8
	tiles      []*Tile
}

var _ Device = (*device)(nil)

func (td *device) String() string {
	return fmt.Sprintf("TileDevice(%v)", td.Target())
}

// Target calls underlying Device's Target function.
func (td *device) Target() lifxlan.Target {
	return td.dev.Target()
}

// Dial calls underlying Device's Dial function.
func (td *device) Dial() (net.Conn, error) {
	return td.dev.Dial()
}

// Source calls underlying Device's Source function.
func (td *device) Source() uint32 {
	return td.dev.Source()
}

// NextSequence calls underlying Device's NextSequence function.
func (td *device) NextSequence() uint8 {
	return td.dev.NextSequence()
}
