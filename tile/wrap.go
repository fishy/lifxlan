package tile

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"

	"github.com/fishy/lifxlan"
	"github.com/fishy/lifxlan/light"
)

// Wrap tries to wrap a lifxlan.Device into a tile device.
//
// When force is false and d is already a tile device,
// d will be casted and returned directly.
// Otherwise, this function calls a tile device API,
// and only return a non-nil Device if it supports the API.
//
// If the device is not a tile device,
// the function might block until ctx is cancelled.
//
// When returning a valid tile device,
// the device's HardwareVersion is guaranteed to be cached.
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

	ld, err := light.Wrap(ctx, d, force)
	if err != nil {
		return nil, err
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

	seq, err := d.Send(
		ctx,
		conn,
		0, // flags
		GetDeviceChain,
		nil, // payload
	)
	if err != nil {
		return nil, err
	}

	for {
		resp, err := lifxlan.ReadNextResponse(ctx, conn)
		if err != nil {
			return nil, err
		}
		if resp.Sequence != seq || resp.Source != d.Source() {
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
		if raw.TotalCount == 0 {
			return nil, errors.New("lifxlan/tile.Wrap: no tiles found")
		}
		*d.HardwareVersion() = raw.TileDevices[int(raw.StartIndex)].HardwareVersion
		td := &device{
			Device:     ld,
			startIndex: raw.StartIndex,
			tiles:      make([]*Tile, raw.TotalCount),
		}
		for i := range td.tiles {
			td.tiles[i] = ParseTile(&raw.TileDevices[int(raw.StartIndex)+i])
		}
		td.parseBoard()
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
