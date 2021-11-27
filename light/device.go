package light

import (
	"context"
	"fmt"
	"net"
	"time"

	"go.yhsif.com/lifxlan"
)

// Device is a wrapped lifxlan.Device that provides light related APIs.
type Device interface {
	lifxlan.Device

	// GetColor returns the current color on this light device.
	//
	// If conn is nil,
	// a new connection will be made and guaranteed to be closed before returning.
	// You should pre-dial and pass in the conn if you plan to call APIs on this
	// device repeatedly.
	GetColor(ctx context.Context, conn net.Conn) (*lifxlan.Color, error)

	// SetColor sets the light device with the given color.
	//
	// If conn is nil,
	// a new connection will be made and guaranteed to be closed before returning.
	// You should pre-dial and pass in the conn if you plan to call APIs on this
	// device repeatedly.
	//
	// If ack is false,
	// this function returns nil error after the API is sent successfully.
	// If ack is true,
	// this function will only return nil error after it received ack from the
	// device.
	SetColor(ctx context.Context, conn net.Conn, color *lifxlan.Color, transition time.Duration, ack bool) error

	// SetLightPower sets the power level of the device and specifies how long it
	// will take to transition to the new power state.
	//
	// If conn is nil,
	// a new connection will be made and guaranteed to be closed before returning.
	// You should pre-dial and pass in the conn if you plan to call APIs on this
	// device repeatedly.
	//
	// If ack is false,
	// this function returns nil error after the API is sent successfully.
	// If ack is true,
	// this function will only return nil error after it received ack from the
	// device.
	SetLightPower(ctx context.Context, conn net.Conn, power lifxlan.Power, transition time.Duration, ack bool) error

	// SetWaveform sends SetWaveformOptional message as defined in
	//
	// https://lan.developer.lifx.com/docs/changing-a-device#setwaveformoptional---packet-119
	SetWaveform(ctx context.Context, conn net.Conn, args *SetWaveformArgs, ack bool) error
}

type device struct {
	lifxlan.Device
}

var _ Device = (*device)(nil)

func (ld *device) String() string {
	if label := ld.Label().String(); label != lifxlan.EmptyLabel {
		return fmt.Sprintf("%s(%v)", label, ld.Target())
	}
	if parsed := ld.HardwareVersion().Parse(); parsed != nil {
		return fmt.Sprintf("%s(%v)", parsed.ProductName, ld.Target())
	}
	return fmt.Sprintf("LightDevice(%v)", ld.Target())
}
