package relay

import (
	"context"
	"fmt"
	"net"

	"go.yhsif.com/lifxlan"
)

// Device is a wrapped lifxlan.Device that provides relay related APIs.
type Device interface {
	lifxlan.Device

	// GetRPower returns the current power level on the relay at index.
	//
	// If conn is nil,
	// a new connection will be made and guaranteed to be closed before returning.
	// You should pre-dial and pass in the conn if you plan to call APIs on this
	// device repeatedly.
	GetRPower(ctx context.Context, conn net.Conn, index uint8) (lifxlan.Power, error)
	// SetRPower sets the power level at the relay index. (Turn it on or off.)
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
	SetRPower(ctx context.Context, conn net.Conn, index uint8, power lifxlan.Power, ack bool) error
}

type device struct {
	lifxlan.Device
}

var _ Device = (*device)(nil)

func (rd *device) String() string {
	if label := rd.Label().String(); label != lifxlan.EmptyLabel {
		return fmt.Sprintf("%s(%v)", label, rd.Target())
	}
	if parsed := rd.HardwareVersion().Parse(); parsed != nil {
		return fmt.Sprintf("%s(%v)", parsed.ProductName, rd.Target())
	}
	return fmt.Sprintf("RelayDevice(%v)", rd.Target())
}
