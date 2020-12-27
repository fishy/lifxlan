package tile

import (
	"go.yhsif.com/lifxlan"
)

// Tile related MessageType values.
const (
	GetDeviceChain   lifxlan.MessageType = 701
	StateDeviceChain lifxlan.MessageType = 702
	GetTileState64   lifxlan.MessageType = 707
	StateTileState64 lifxlan.MessageType = 711
	SetTileState64   lifxlan.MessageType = 715
)
