package lifxlan

// MessageType is the 16-bit header indicates the type of the message.
type MessageType uint16

// MessageType values.
const (
	Acknowledgement MessageType = 45
	GetService                  = 2
	StateService                = 3
	GetPower                    = 20
	StatePower                  = 22
	SetPower                    = 21
	GetLabel                    = 23
	StateLabel                  = 25
	GetVersion                  = 32
	StateVersion                = 33
)
