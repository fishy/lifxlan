package lifxlan

import (
	"context"
	"fmt"
	"net"
	"sync/atomic"
)

// ServiceType define the type of the service this device provides.
type ServiceType uint8

// Documented values for ServiceType.
const (
	ServiceUDP ServiceType = 1
)

func (s ServiceType) String() string {
	switch s {
	default:
		return fmt.Sprintf("<UNKNOWN> (%d)", uint8(s))
	case ServiceUDP:
		return "UDP"
	}
}

// Device defines the common interface between lifxlan devices.
//
// For the Foo() and GetFoo() function pairs (e.g. Label() and GetLabel()),
// the Foo() one will return an pointer to the cached property,
// guaranteed to be non-nil but could be the zero value,
// while the GetFoo() one will use an API call to update the cached property.
//
// There will also be an EmptyFoo string constant defined,
// so that you can compare against Device.Foo().String() to determine if a
// GetFoo() call is needed.
// Here is an example code snippet to get a device's label:
//
//     func GetLabel(ctx context.Context, d lifxlanDevice) (string, error) {
//         if d.Label().String() != lifxlan.EmptyLabel {
//             return d.Label().String(), nil
//         }
//         if err := d.GetLabel(ctx, nil); err = nil {
//             return "", nil
//         }
//         return d.Label().String(), nil
//     }
//
// If you are extending a device code and you got the property as part of
// another API's return payload,
// you can also use the Foo() function to update the cached value.
// Here is an example code snippet to update a device's cached label:
//
//     func UpdateLabel(d lifxlanDevice, newLabel *lifxlan.RawLabel) {
//         *d.Label() = *newLabel
//     }
//
// The conn arg in GetFoo() functions can be nil.
// In such cases,
// a new connection will be made and guaranteed to be closed before returning.
// You should pre-dial and pass in the conn if you plan to call APIs on this
// device repeatedly.
//
// In case of network error (e.g. response packet loss),
// the GetFoo() functions might block until the context is cancelled,
// as a result, it's a good idea to set a timeout to the context.
type Device interface {
	// Target returns the target of this device, usually it's the MAC address.
	Target() Target

	// Dial tries to establish a connection to this device.
	Dial() (net.Conn, error)

	// Source returns a consistent random source to be used with API calls.
	// It's guaranteed to be non-zero.
	Source() uint32

	// NextSequence returns the next sequence value to be used with API calls.
	NextSequence() uint8

	// Send generates and sends a message to the device.
	//
	// conn must be pre-dialed or this function will fail.
	//
	// It calls the device's Target(), Source(), and NextSequence() functions to
	// fill the appropriate headers.
	//
	// The sequence used in this message will be returned.
	Send(ctx context.Context, conn net.Conn, flags AckResFlag, message MessageType, payload []byte) (seq uint8, err error)

	// The label of the device.
	Label() *RawLabel
	GetLabel(ctx context.Context, conn net.Conn) error

	// The hardware version info of the device.
	HardwareVersion() *RawHardwareVersion
	GetHardwareVersion(ctx context.Context, conn net.Conn) error
}

var _ Device = (*device)(nil)

// device defines the base type of a lifxlan device.
type device struct {
	// The network address, in "ip:port" format.
	addr string
	// The type of service this device provides.
	service ServiceType
	// The target of this device, usually it's the MAC address.
	target Target

	source   uint32
	sequence uint32

	// Cached properties.
	label   RawLabel
	version RawHardwareVersion
}

// NewDevice creates a new Device.
//
// addr must be in "host:port" format and service must be a known service type,
// otherwise the later Dial funcion will fail.
func NewDevice(addr string, service ServiceType, target Target) Device {
	return &device{
		addr:    addr,
		service: service,
		target:  target,
		source:  RandomSource(),
	}
}

func (d *device) String() string {
	if label := d.Label().String(); label != EmptyLabel {
		return fmt.Sprintf("%s(%v)", label, d.Target())
	}
	if name := d.HardwareVersion().Parse().ProductName; name != "" {
		return fmt.Sprintf("%s(%v)", name, d.Target())
	}
	return fmt.Sprintf("Device(%v)", d.Target())
}

func (d *device) Target() Target {
	return d.target
}

func (d *device) Dial() (net.Conn, error) {
	var network string
	switch d.service {
	default:
		return nil, fmt.Errorf(
			"lifxlan.Device.Dial: unknown device service type: %v",
			d.service,
		)
	case ServiceUDP:
		network = "udp"
	}
	return net.Dial(network, d.addr)
}

func (d *device) Source() uint32 {
	return d.source
}

const uint8mask = uint32(0xff)

func (d *device) NextSequence() uint8 {
	return uint8(atomic.AddUint32(&d.sequence, 1) & uint8mask)
}
