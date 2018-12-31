package lifxlan

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
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

// StateDeviceChain message payload offsets.
const (
	StateDeviceChainStartIndexOffset  = 0                                                  // 1 byte
	StateDeviceChainTileDevicesOffset = StateDeviceChainStartIndexOffset + 1               // TileBufSize*16 bytes
	StateDeviceChainTotalCountOffset  = StateDeviceChainTileDevicesOffset + TileBufSize*16 // 1 byte
)

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
	msg := GenerateMessage(
		NotTagged,
		d.Source(),
		d.Target(),
		0, // flags
		sequence,
		GetDeviceChain,
		nil, // payload
	)
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

		startIndex := uint8(resp.Payload[StateDeviceChainStartIndexOffset])
		numDevices := uint8(resp.Payload[StateDeviceChainTotalCountOffset])
		fmt.Println("TileBufSize:", TileBufSize)
		tileBuf := make([]byte, TileBufSize)
		tiles := make([]*Tile, numDevices)
		for i := range tiles {
			offset := StateDeviceChainTileDevicesOffset + TileBufSize*(i+int(startIndex))
			copy(tileBuf, resp.Payload[offset:])
			tiles[i] = ParseTile(tileBuf)
			fmt.Printf("%+v\n", tiles[i])
		}
		return &TileDevice{
			dev:        d,
			startIndex: startIndex,
			tiles:      tiles,
		}, nil
	}
}

// Tile buf offsets.
const (
	TileAccelMeasXOffset    = 0                            // 2 bytes, signed
	TileAccelMeasYOffset    = TileAccelMeasXOffset + 2     // 2 bytes, signed
	TileAccelMeasZOffset    = TileAccelMeasYOffset + 2     // 2 bytes, signed
	tileReserved1Offset     = TileAccelMeasZOffset + 2     // 2 bytes
	TileUserXOffset         = tileReserved1Offset + 2      // 4 bytes, float
	TileUserYOffset         = TileUserXOffset + 4          // 4 bytes, float
	TileWidthOffset         = TileUserYOffset + 4          // 1 byte
	TileHeightOffset        = TileWidthOffset + 1          // 1 byte
	tileReserved2Offset     = TileHeightOffset + 1         // 1 byte
	TileDeviceVersionOffset = tileReserved2Offset + 1      // 12 bytes
	TileFirmwareOffset      = TileDeviceVersionOffset + 12 // 20 bytes
	tileReserved3Offset     = TileFirmwareOffset + 20      // 4 bytes

	TileBufSize = tileReserved3Offset + 4
)

// Tile defines a single tile inside a TileDevice
type Tile struct {
	UserX    float32
	UserY    float32
	Width    uint8
	Height   uint8
	Rotation TileRotation
}

// ParseTile parses buf into a Tile.
//
// buf must be of the length of TileBufSize,
// otherwise this functino might panic.
func ParseTile(buf []byte) *Tile {
	var x, y, z int16
	tile := &Tile{}
	// map of offset -> pointer
	table := map[int64]interface{}{
		TileAccelMeasXOffset: &x,
		TileAccelMeasYOffset: &y,
		TileAccelMeasZOffset: &z,
		TileUserXOffset:      &tile.UserX,
		TileUserYOffset:      &tile.UserY,
		TileWidthOffset:      &tile.Width,
		TileHeightOffset:     &tile.Height,
	}
	r := bytes.NewReader(buf)
	for offset, pointer := range table {
		// Seek only returns error when whence is invalid,
		// or when absolute offset is negative.
		// Neither should happen here.
		if _, err := r.Seek(offset, io.SeekStart); err != nil {
			panic(err)
		}
		// Read only returns error regarding EOF,
		// which means buf is not big enough.
		if err := binary.Read(r, binary.LittleEndian, pointer); err != nil {
			panic(err)
		}
	}
	tile.Rotation = ParseTileRotation(x, y, z)
	return tile
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
