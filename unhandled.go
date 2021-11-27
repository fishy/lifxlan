package lifxlan

import (
	"fmt"
)

// RawStateUnhandledPayload defines the struct to be used for encoding and decoding.
//
// https://lan.developer.lifx.com/docs/information-messages#stateunhandled---packet-223
type RawStateUnhandledPayload struct {
	UnhandledType MessageType
}

func (p RawStateUnhandledPayload) Error() string {
	return fmt.Sprintf(
		"lifxlan: unhandled message: %v",
		p.UnhandledType,
	)
}
