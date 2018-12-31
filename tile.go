package lifxlan

import (
	"context"
	"fmt"
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
	TileBufSize = 55

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
		tileBuf := make([]byte, TileBufSize)
		tiles := make([]*Tile, numDevices)
		for i := range tiles {
			offset := StateDeviceChainTileDevicesOffset + TileBufSize*(i+int(startIndex))
			copy(tileBuf, resp.Payload[offset:])
			tiles[i] = ParseTile(tileBuf)
		}
		return &TileDevice{
			dev:        d,
			startIndex: startIndex,
			tiles:      tiles,
		}, nil
	}
}

// Tile defines a single tile inside a TileDevice
type Tile struct {
	// TODO
}

// ParseTile parses buf into a Tile.
//
// buf must be of the length of TileBufSize,
// otherwise this functino might panic.
func ParseTile(buf []byte) *Tile {
	// TODO
	fmt.Printf("%d: % x\n", len(buf), buf)
	return &Tile{}
}
