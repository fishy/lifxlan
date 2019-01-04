package tile

import (
	"fmt"
	"math"
)

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
// Then the width of the board would be 24, and the height would be 12.
// The coordinate of the topleft corner would be (0, 11) and the coordinate
// of the bottomright corner of the middle tile would be (15, 0).
//
// The size is also guaranteed to be trimmed.
// Which means that on the above example,
// if either the left or the right tile is removed from the device,
// the size would change to 16x12.
// But if the middle tile is removed from the device,
// the size would change to 24x8 with a 8x8 blackhole in the middle.
type Board interface {
	Width() int
	Height() int

	// OnTile returns true if coordinate (x, y) is on a tile.
	OnTile(x, y int) bool
}

// IndexData stores the data linked to a tile for a Board coordinate.
type IndexData struct {
	// The coordinate inside the tile.
	Coordinate

	// The index of the tile.
	Index int
}

func (id IndexData) String() string {
	return fmt.Sprintf("%d:(%d,%d)", id.Index, id.X, id.Y)
}

// BoardData is the parsed, normalized board data.
//
// The zero value represents an empty board of size 0x0.
type BoardData struct {
	// The size of the board.
	Coordinate

	// Parsed index data, with size X*Y.
	Data [][]*IndexData
	// Parsed reverse coordinate data,
	// with size len(tiles)*tileWidth*tileHeight,
	// The coordinate is the coordinate of this tile pixel on the board.
	ReverseData [][][]Coordinate
}

func (td *device) parseBoard() {
	td.board = ParseBoard(td.tiles)
}

// ParseBoard parses tiles into BoardData.
func ParseBoard(tiles []*Tile) BoardData {
	var board BoardData
	var base Coordinate
	bs := make([][][]Coordinate, len(tiles))
	board.ReverseData = make([][][]Coordinate, len(tiles))

	base.X = int(math.MaxInt32)
	base.Y = int(math.MaxInt32)
	board.X = int(math.MinInt32)
	board.Y = int(math.MinInt32)

	for i, tile := range tiles {
		var min, max Coordinate
		bs[i], min, max = tile.BoardCoordinates()

		if min.X < base.X {
			base.X = min.X
		}
		if min.Y < base.Y {
			base.Y = min.Y
		}
		if max.X > board.X {
			board.X = max.X
		}
		if max.Y > board.Y {
			board.Y = max.Y
		}
	}

	board.X -= base.X
	board.Y -= base.Y
	board.Data = make([][]*IndexData, board.X)
	for i := range board.Data {
		board.Data[i] = make([]*IndexData, board.Y)
	}

	for i, b := range bs {
		board.ReverseData[i] = make([][]Coordinate, len(b))
		for x := range b {
			board.ReverseData[i][x] = make([]Coordinate, len(b[x]))
			for y, c := range b[x] {
				data := &IndexData{
					Coordinate: Coordinate{
						X: x,
						Y: y,
					},
					Index: i,
				}
				bx := c.X - base.X
				by := c.Y - base.Y

				board.Data[bx][by] = data
				board.ReverseData[i][x][y] = Coordinate{
					X: bx,
					Y: by,
				}
			}
		}
	}

	return board
}

func (td *device) Width() int {
	return td.board.X
}

func (td *device) Height() int {
	return td.board.Y
}

func (td *device) OnTile(x, y int) bool {
	if x < 0 || x >= td.Width() || y < 0 || y >= td.Height() {
		return false
	}
	return td.board.Data[x][y] != nil
}
