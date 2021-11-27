package tile

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"

	"go.yhsif.com/lifxlan"
	"go.yhsif.com/lifxlan/light"
)

// Wrap tries to wrap a lifxlan.Device into a tile device.
//
// When force is false and d is already a tile device,
// d will be casted and returned directly.
// Otherwise, this function calls a tile device API,
// and only returns a non-nil Device if it supports the API.
//
// If the device is not a tile device,
// the function might block until ctx is cancelled.
//
// When returning a valid tile device,
// the device's HardwareVersion is guaranteed to be cached.
func Wrap(ctx context.Context, d lifxlan.Device, force bool) (Device, error) {
	if ctx.Err() != nil {
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

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	const msg = GetDeviceChain

	seq, err := d.Send(
		ctx,
		conn,
		0, // flags
		msg,
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

		switch resp.Message {
		case StateDeviceChain:
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

		case lifxlan.StateUnhandled:
			var raw lifxlan.RawStateUnhandledPayload
			r := bytes.NewReader(resp.Payload)
			if err := binary.Read(r, binary.LittleEndian, &raw); err != nil {
				return nil, err
			}
			return nil, raw
		}
	}
}

// RawStateDeviceChainPayload defines the struct to be used for encoding and
// decoding.
//
// https://lan.developer.lifx.com/docs/information-messages#statedevicechain---packet-702
type RawStateDeviceChainPayload struct {
	StartIndex  uint8
	TileDevices [16]RawTileDevice
	TotalCount  uint8
}
