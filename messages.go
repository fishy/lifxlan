package lifxlan

// MessageType is the 16-bit header indicates the type of the message.
type MessageType uint16

// MessageType values.
const (
	// Discover
	GetService   MessageType = 2
	StateService             = 3

	// Tile
	GetDeviceChain   = 701
	StateDeviceChain = 702
)
