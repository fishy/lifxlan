package tile

// Board defines a board for a tile device (collection of tiles).
//
// Assuming every tile is 8x8, and we have 3 tiles arranged as:
//
//     +--+    +--+
//     |  |    |  |
//     |  |+--+|  |
//     +--+|  |+--+
//         |  |
//         +--+
//
// Then the width of the board will be 24, and the height would be 12.
// The coordination of the topleft corner would be (0, 0) and the coordination
// of the bottomright corner of the middle tile would be (15, 11).
//
// If the board only contains a single tile,
// the Width and Height would both be 8 and for (x, y) with 0 <= x,y < 8 Shown
// should always return true.
type Board interface {
	Width() int
	Height() int

	// OnTile returns true if coordination (x, y) is on a tile.
	OnTile(x, y int) bool
}

func (td *device) initBoard() {
	// TODO
}

func (td *device) Width() int {
	// TODO
	return 0
}

func (td *device) Height() int {
	// TODO
	return 0
}

func (td *device) OnTile(x, y int) bool {
	// TODO
	return false
}
