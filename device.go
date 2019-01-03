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
type Device interface {
	// Target returns the target of this device, usually it's the MAC address.
	Target() Target

	// Dial tries to establish a connection to this device.
	Dial() (net.Conn, error)

	// Source returns a consistent random source to be used with API calls.
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
	Send(ctx context.Context, conn net.Conn, tagged TaggedHeader, flags AckResFlag, message MessageType, payload []byte) (seq uint8, err error)
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
	return fmt.Sprintf("Device(%v)", d.target)
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
