package lifxlan

// MessageType is the 16-bit header indicates the type of the message.
type MessageType uint16

// MessageType values.
const (
	Acknowledgement MessageType = 45
	StateUnhandled  MessageType = 223

	GetService        MessageType = 2
	StateService      MessageType = 3
	GetHostFirmware   MessageType = 14
	StateHostFirmware MessageType = 15
	GetPower          MessageType = 20
	StatePower        MessageType = 22
	SetPower          MessageType = 21
	GetLabel          MessageType = 23
	StateLabel        MessageType = 25
	GetVersion        MessageType = 32
	StateVersion      MessageType = 33
	EchoRequest       MessageType = 58
	EchoResponse      MessageType = 59
)
