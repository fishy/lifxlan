package light

import (
	"context"

	"github.com/fishy/lifxlan"
)

// Wrap tries to wrap a lifxlan.Device into a light device.
//
// When force is false and d is already a tile device,
// d will be casted and returned directly.
// Otherwise, this function might call a device API to determine whether it's a
// light device.
//
// If the device is not a light device,
// the function might block until ctx is cancelled.
//
// Currently all lifx devices are light device,
// so this function doesn't really call an API and just do a naive wrapping.
// But that might change in the future so you shouldn't be relying on that.
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

	return &device{
		Device: d,
	}, nil
}
