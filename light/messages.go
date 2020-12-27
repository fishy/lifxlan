package light

import (
	"go.yhsif.com/lifxlan"
)

// Light related MessageType values.
const (
	Get                 lifxlan.MessageType = 101
	SetColor            lifxlan.MessageType = 102
	State               lifxlan.MessageType = 107
	SetWaveformOptional lifxlan.MessageType = 119
)
