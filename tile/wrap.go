package tile

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"

	"github.com/fishy/lifxlan"
)

// Wrap tries to wrap a lifxlan.Device into a tile device.
//
// When force is false and d is already a tile Device,
// d will be casted and returned directly.
// Otherwise, this function calls a tile device API,
// and only return a non-nil Device if it supports the API.
//
// If the device is not a tile device,
// the function might return nil Device and nil error.
// The function might also block until ctx is cancelled.
func Wrap(ctx context.Context, d lifxlan.Device, force bool) (Device, error) {
	select {
	default:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	if !force {
		if t, ok := d.(Device); ok {
			return t, nil
		}
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
	msg, err := lifxlan.GenerateMessage(
		lifxlan.NotTagged,
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

	buf := make([]byte, lifxlan.ResponseReadBufferSize)
	for {
		select {
		default:
		case <-ctx.Done():
			return nil, ctx.Err()
		}

		if err := conn.SetReadDeadline(lifxlan.GetReadDeadline()); err != nil {
			return nil, err
		}

		n, err := conn.Read(buf)
		if err != nil {
			if lifxlan.CheckTimeoutError(err) {
				continue
			}
			return nil, err
		}

		resp, err := lifxlan.ParseResponse(buf[:n])
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
		td := &device{
			dev:        d,
			startIndex: raw.StartIndex,
			tiles:      tiles,
		}
		td.initBoard()
		return td, nil
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
