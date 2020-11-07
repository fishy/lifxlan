package lifxlan

import (
	"fmt"
)

// RawStateUnhandledPayload defines the struct to be used for encoding and decoding.
//
// https://lan.developer.lifx.com/v2.0/docs/the-lifx-switch#section-stateunhandled-223
type RawStateUnhandledPayload struct {
	UnhandledType MessageType
}

// GenerateError generates an error regarding this unhandled message.
func (p RawStateUnhandledPayload) GenerateError(expected MessageType) error {
	if expected == p.UnhandledType {
		return fmt.Errorf(
			"lifxlan.RawStateUnhandledPayload: unhandled message: %v",
			p.UnhandledType,
		)
	}
	return fmt.Errorf(
		"lifxlan.RawStateUnhandledPayload: unexpected UnhandledType %v, was expecting %v",
		p.UnhandledType,
		expected,
	)
}
