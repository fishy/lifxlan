package tile

import (
	"github.com/fishy/lifxlan"
)

// Tile related MessageType values.
const (
	GetDeviceChain   lifxlan.MessageType = 701
	StateDeviceChain                     = 702
	GetTileState64                       = 707
	StateTileState64                     = 711
	SetTileState64                       = 715
)
