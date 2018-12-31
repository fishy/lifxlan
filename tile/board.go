package tile

import (
	"fmt"
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
// Then the width of the board will be 24, and the height would be 12.
// The coordinate of the topleft corner would be (0, 11) and the coordinate
// of the bottomright corner of the middle tile would be (15, 0).
type Board interface {
	Width() int
	Height() int

	// OnTile returns true if coordinate (x, y) is on a tile.
	OnTile(c Coordinate) bool
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
// Empty value represents an empty board of size 0x0.
type BoardData struct {
	Coordinate

	Data [][]*IndexData
}

func (td *device) parseBoard() {
	td.board = ParseBoard(td.tiles)
}

// ParseBoard parses tiles into BoardData.
func ParseBoard(tiles []*Tile) BoardData {
	var board BoardData
	var offset Coordinate
	for i, tile := range tiles {
		b, min, max := tile.BoardCoordinates()

		if min.X+offset.X < 0 {
			// Need to pad x to the head.
			newOffset := 0 - min.X - offset.X
			newData := make([][]*IndexData, len(board.Data)+newOffset)
			for i := 0; i < newOffset; i++ {
				newData[i] = make([]*IndexData, board.Y)
			}
			copy(newData[newOffset:], board.Data)
			board.Data = newData
			offset.X += newOffset
			board.X += newOffset
		}
		if max.X+offset.X > board.X {
			// Need to pad x to the tail.
			sizeDiff := max.X + offset.X - board.X
			for i := 0; i < sizeDiff; i++ {
				board.Data = append(board.Data, make([]*IndexData, board.Y))
			}
			board.X += sizeDiff
		}

		if min.Y+offset.Y < 0 {
			// Need to pad y to the head.
			newOffset := 0 - min.Y - offset.Y
			for i := range board.Data {
				newRow := make([]*IndexData, board.Y+newOffset)
				copy(newRow[newOffset:], board.Data[i])
				board.Data[i] = newRow
			}
			offset.Y += newOffset
			board.Y += newOffset
		}
		if max.Y+offset.Y > board.Y {
			// Need to pad y to the tail.
			sizeDiff := max.Y + offset.Y - board.Y
			for i := range board.Data {
				newRow := make([]*IndexData, board.Y+sizeDiff)
				copy(newRow, board.Data[i])
				board.Data[i] = newRow
			}
			board.Y += sizeDiff
		}

		for x := range b {
			for y, c := range b[x] {
				data := &IndexData{
					Coordinate: Coordinate{
						X: x,
						Y: y,
					},
					Index: i,
				}
				board.Data[c.X+offset.X][c.Y+offset.Y] = data
			}
		}
	}

	// Trim down the board
	var max Coordinate
	min := Coordinate{
		X: board.X,
		Y: board.Y,
	}
	for x := 0; x < board.X; x++ {
		for y := 0; y < board.Y; y++ {
			if board.Data[x][y] != nil {
				if x < min.X {
					min.X = x
				}
				if x+1 > max.X {
					max.X = x + 1
				}
				if y < min.Y {
					min.Y = y
				}
				if y+1 > max.Y {
					max.Y = y + 1
				}
			}
		}
	}
	if min.X > 0 || max.X < board.X {
		newData := make([][]*IndexData, max.X-min.X)
		copy(newData, board.Data[min.X:])
		board.Data = newData
	}
	board.X = max.X - min.X
	if min.Y > 0 || max.Y < board.Y {
		for i := range board.Data {
			newRow := make([]*IndexData, max.Y-min.Y)
			copy(newRow, board.Data[i][min.Y:])
			board.Data[i] = newRow
		}
	}
	board.Y = max.Y - min.Y

	return board
}

func (td *device) Width() int {
	return td.board.X
}

func (td *device) Height() int {
	return td.board.Y
}

func (td *device) OnTile(c Coordinate) bool {
	if c.X < 0 || c.X >= td.Width() || c.Y < 0 || c.Y >= td.Height() {
		return false
	}
	return td.board.Data[c.X][c.Y] != nil
}
