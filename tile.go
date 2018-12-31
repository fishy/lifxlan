package lifxlan

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"math"
	"net"
)

// TileDevice is a wrapped Device that provides tile related APIs.
type TileDevice struct {
	dev        *device
	startIndex uint8
	tiles      []*Tile
}

var _ Device = (*TileDevice)(nil)

func (td *TileDevice) String() string {
	return fmt.Sprintf("TileDevice(%v)", td.Target())
}

// Target calls underlying Device's Target function.
func (td *TileDevice) Target() Target {
	return td.dev.Target()
}

// Dial calls underlying Device's Dial function.
func (td *TileDevice) Dial() (net.Conn, error) {
	return td.dev.Dial()
}

// Source calls underlying Device's Source function.
func (td *TileDevice) Source() uint32 {
	return td.dev.Source()
}

// NextSequence calls underlying Device's NextSequence function.
func (td *TileDevice) NextSequence() uint8 {
	return td.dev.NextSequence()
}

// GetTileDevice returns self.
//
// The only possibility it returns an error is that the context is already
// cancelled when entering this function.
func (td *TileDevice) GetTileDevice(ctx context.Context) (*TileDevice, error) {
	select {
	default:
		return td, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// RawStateDeviceChainPayload defines the struct to be used for encoding and
// decoding.
//
// https://lan.developer.lifx.com/v2.0/docs/tile-messages#section-statedevicechain-702
type RawStateDeviceChainPayload struct {
	StartIndex  uint8
	TileDevices [16]RawTileDevice
	TotalCount  uint8
}

func (d *device) GetTileDevice(ctx context.Context) (*TileDevice, error) {
	select {
	default:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	conn, err := d.Dial()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	select {
	default:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	sequence := d.NextSequence()
	msg, err := GenerateMessage(
		NotTagged,
		d.Source(),
		d.Target(),
		0, // flags
		sequence,
		GetDeviceChain,
		nil, // payload
	)
	if err != nil {
		return nil, err
	}

	n, err := conn.Write(msg)
	if err != nil {
		return nil, err
	}
	if n < len(msg) {
		return nil, fmt.Errorf(
			"lifxlan.Device.GetTileDevice: only wrote %d out of %d bytes",
			n,
			len(msg),
		)
	}

	buf := make([]byte, ResponseReadBufferSize)
	for {
		select {
		default:
		case <-ctx.Done():
			return nil, ctx.Err()
		}

		if err := conn.SetReadDeadline(getReadDeadline()); err != nil {
			return nil, err
		}

		n, err := conn.Read(buf)
		if err != nil {
			if timeoutErr, ok := err.(timeouter); ok {
				if timeoutErr.Timeout() {
					continue
				}
			}
			return nil, err
		}

		resp, err := ParseResponse(buf[:n])
		if err != nil {
			return nil, err
		}
		if resp.Sequence != sequence || resp.Source != d.Source() {
			continue
		}
		if resp.Message != StateDeviceChain {
			continue
		}

		var raw RawStateDeviceChainPayload
		r := bytes.NewReader(resp.Payload)
		if err := binary.Read(r, binary.LittleEndian, &raw); err != nil {
			return nil, err
		}
		tiles := make([]*Tile, raw.TotalCount)
		for i := range tiles {
			tiles[i] = ParseTile(&raw.TileDevices[int(raw.StartIndex)+i])
			fmt.Printf("%+v\n", tiles[i])
		}
		return &TileDevice{
			dev:        d,
			startIndex: raw.StartIndex,
			tiles:      tiles,
		}, nil
	}
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
	Rotation TileRotation
}

// ParseTile parses RawTileDevice into a Tile.
func ParseTile(raw *RawTileDevice) *Tile {
	return &Tile{
		UserX:    raw.UserX,
		UserY:    raw.UserY,
		Width:    raw.Width,
		Height:   raw.Height,
		Rotation: ParseTileRotation(raw.AccelMeasX, raw.AccelMeasY, raw.AccelMeasZ),
	}
}

// TileRotation defines the rotation of a single tile.
type TileRotation int

// TileRotation values
const (
	TileRotationRightSideUp TileRotation = iota
	TileRotationRotateRight
	TileRotationRotateLeft
	TileRotationFaceDown
	TileRotationFaceUp
	TileRotationUpsideDown
)

func (r TileRotation) String() string {
	switch r {
	default:
		return fmt.Sprintf("<Invalid value %d>", int(r))
	case TileRotationRightSideUp:
		return "RightSideUp"
	case TileRotationRotateRight:
		return "RotateRight"
	case TileRotationRotateLeft:
		return "RotateLeft"
	case TileRotationFaceDown:
		return "FaceDown"
	case TileRotationFaceUp:
		return "FaceUp"
	case TileRotationUpsideDown:
		return "UpsideDown"
	}
}

// ParseTileRotation parses measurements into TileRotation
func ParseTileRotation(x, y, z int16) TileRotation {
	abs := func(x int16) int16 {
		return int16(math.Abs(float64(x)))
	}

	// Copied from:
	// https://lan.developer.lifx.com/v2.0/docs/tile-messages#section-tile
	absX := abs(x)
	absY := abs(y)
	absZ := abs(z)

	if x == -1 && y == -1 && z == -1 {
		// Invalid data, assume right-side up.
		return TileRotationRightSideUp
	}
	if absX > absY && absX > absZ {
		if x > 0 {
			return TileRotationRotateRight
		}
		return TileRotationRotateLeft
	}

	if absZ > absX && absZ > absY {
		if z > 0 {
			return TileRotationFaceDown
		}
		return TileRotationFaceUp
	}

	if y > 0 {
		return TileRotationUpsideDown
	}
	return TileRotationRightSideUp
}
