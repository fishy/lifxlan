package tile

import (
	"math"
)

// Coordinate defines a simple 2D coordinate.
type Coordinate struct {
	X, Y int
}

// RawTileDevice defines the struct to be used for encoding and decoding.
//
// https://lan.developer.lifx.com/v2.0/docs/tile-messages#section-tile
type RawTileDevice struct {
	AccelMeasX int16
	AccelMeasY int16
	AccelMeasZ int16
	_          int16 // reserved
	UserX      float32
	UserY      float32
	Width      uint8
	Height     uint8
	_          uint8  // reserved
	_          uint32 // device_version_vendor
	_          uint32 // device_version_product
	_          uint32 // device_version_version
	_          uint64 // firmware_build
	_          uint64 // reserved
	_          uint32 // firmware_versio
	_          uint32 // reserved
}

// Tile defines a single tile inside a TileDevice
type Tile struct {
	UserX    float32
	UserY    float32
	Width    uint8
	Height   uint8
	Rotation Rotation
}

// ParseTile parses RawTileDevice into a Tile.
func ParseTile(raw *RawTileDevice) *Tile {
	return &Tile{
		UserX:    raw.UserX,
		UserY:    raw.UserY,
		Width:    raw.Width,
		Height:   raw.Height,
		Rotation: ParseRotation(raw.AccelMeasX, raw.AccelMeasY, raw.AccelMeasZ),
	}
}

// BoardCoordinates returns non-normalized coordinates of the pixels on this
// tile on the board.
//
// "non-normalized" means that the coornidate might be negative.
//
// The returned coordinates are guaranteed to be of the size of Width*Height.
func (t Tile) BoardCoordinates() (
	coordinates [][]Coordinate,
	min Coordinate,
	max Coordinate,
) {
	min.X = int(math.MaxInt32)
	min.Y = int(math.MaxInt32)
	max.X = int(math.MinInt32)
	max.Y = int(math.MinInt32)

	coordinates = make([][]Coordinate, t.Width)
	for i := range coordinates {
		coordinates[i] = make([]Coordinate, t.Height)
	}

	baseX := int(float32(t.Width) * t.UserX)
	baseY := int(float32(t.Height) * t.UserY)
	for i := 0; i < int(t.Width); i++ {
		for j := 0; j < int(t.Height); j++ {
			switch t.Rotation {
			default:
				// TODO: handle other rotations correctly.
				fallthrough
			case RotationRightSideUp:
				x := i + baseX
				if x < min.X {
					min.X = x
				}
				if x+1 > max.X {
					max.X = x + 1
				}
				y := j + baseY
				if y < min.Y {
					min.Y = y
				}
				if y+1 > max.Y {
					max.Y = y + 1
				}
				coordinates[i][j] = Coordinate{
					X: x,
					Y: y,
				}
			}
		}
	}
	return
}