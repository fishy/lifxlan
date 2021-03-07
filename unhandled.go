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

func (p RawStateUnhandledPayload) Error() string {
	return fmt.Sprintf(
		"lifxlan: unhandled message: %v",
		p.UnhandledType,
	)
}
