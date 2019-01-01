// Package lifxlan provides API implemented in Go for LIFX LAN Protocol:
//
// https://lan.developer.lifx.com/v2.0/docs/
//
// This package focuses on the base stuff, device discovery,
// and capabilities shared by all types of devices.
// For more concreted capabilities like light control and tile control,
// please refer to the subpackages.
//
// Currently this package and its subpackages are not complete and implement all
// possible LIFX LAN Protocols is not the current goal of this package.
// The design choice for this package is that it exposes as much as possible,
// so another package can implement missing device APIs by wrapping a device
// returned by discovery using only exported functions.
// Please refer to the subpackages for an example of extending device
// capabilities.
//
// The API is unstable right now,
// but the maintainers try very hard not to break them.
package lifxlan
