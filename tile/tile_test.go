package tile_test

import (
	"reflect"
	"testing"

	"go.yhsif.com/lifxlan/tile"
)

// To grately reduce typing in test code.
func c(x, y int) tile.Coordinate {
	return tile.Coordinate{
		X: x,
		Y: y,
	}
}

const rotation = tile.RotationRightSideUp

func TestTileBoardCooridates(t *testing.T) {
	t.Run(
		"0,0",
		func(t *testing.T) {
			ti := tile.Tile{
				UserX:    0,
				UserY:    0,
				Width:    8,
				Height:   8,
				Rotation: rotation,
			}
			expectedMin := c(0, 0)
			expectedMax := c(8, 8)
			expectedCs := [][]tile.Coordinate{
				{
					c(0, 7),
					c(1, 7),
					c(2, 7),
					c(3, 7),
					c(4, 7),
					c(5, 7),
					c(6, 7),
					c(7, 7),
				},
				{
					c(0, 6),
					c(1, 6),
					c(2, 6),
					c(3, 6),
					c(4, 6),
					c(5, 6),
					c(6, 6),
					c(7, 6),
				},
				{
					c(0, 5),
					c(1, 5),
					c(2, 5),
					c(3, 5),
					c(4, 5),
					c(5, 5),
					c(6, 5),
					c(7, 5),
				},
				{
					c(0, 4),
					c(1, 4),
					c(2, 4),
					c(3, 4),
					c(4, 4),
					c(5, 4),
					c(6, 4),
					c(7, 4),
				},
				{
					c(0, 3),
					c(1, 3),
					c(2, 3),
					c(3, 3),
					c(4, 3),
					c(5, 3),
					c(6, 3),
					c(7, 3),
				},
				{
					c(0, 2),
					c(1, 2),
					c(2, 2),
					c(3, 2),
					c(4, 2),
					c(5, 2),
					c(6, 2),
					c(7, 2),
				},
				{
					c(0, 1),
					c(1, 1),
					c(2, 1),
					c(3, 1),
					c(4, 1),
					c(5, 1),
					c(6, 1),
					c(7, 1),
				},
				{
					c(0, 0),
					c(1, 0),
					c(2, 0),
					c(3, 0),
					c(4, 0),
					c(5, 0),
					c(6, 0),
					c(7, 0),
				},
			}
			cs, min, max := ti.BoardCoordinates()
			if !reflect.DeepEqual(cs, expectedCs) {
				t.Errorf("coordinates expected %+v, got %+v", expectedCs, cs)
			}
			if !reflect.DeepEqual(min, expectedMin) {
				t.Errorf("min expected %+v, got %+v", expectedMin, min)
			}
			if !reflect.DeepEqual(max, expectedMax) {
				t.Errorf("max expected %+v, got %+v", expectedMax, max)
			}
		},
	)

	t.Run(
		"-0.5,0.5",
		func(t *testing.T) {
			ti := tile.Tile{
				UserX:    -0.5,
				UserY:    0.5,
				Width:    2,
				Height:   2,
				Rotation: rotation,
			}
			expectedMin := c(-1, 1)
			expectedMax := c(1, 3)
			expectedCs := [][]tile.Coordinate{
				{
					c(-1, 2),
					c(0, 2),
				},
				{
					c(-1, 1),
					c(0, 1),
				},
			}
			cs, min, max := ti.BoardCoordinates()
			if !reflect.DeepEqual(cs, expectedCs) {
				t.Errorf("coordinates expected %+v, got %+v", expectedCs, cs)
			}
			if !reflect.DeepEqual(min, expectedMin) {
				t.Errorf("min expected %+v, got %+v", expectedMin, min)
			}
			if !reflect.DeepEqual(max, expectedMax) {
				t.Errorf("max expected %+v, got %+v", expectedMax, max)
			}
		},
	)

	t.Run(
		"2,-2",
		func(t *testing.T) {
			ti := tile.Tile{
				UserX:    2,
				UserY:    -2,
				Width:    2,
				Height:   2,
				Rotation: rotation,
			}
			expectedMin := c(4, -4)
			expectedMax := c(6, -2)
			expectedCs := [][]tile.Coordinate{
				{
					c(4, -3),
					c(5, -3),
				},
				{
					c(4, -4),
					c(5, -4),
				},
			}
			cs, min, max := ti.BoardCoordinates()
			if !reflect.DeepEqual(cs, expectedCs) {
				t.Errorf("coordinates expected %+v, got %+v", expectedCs, cs)
			}
			if !reflect.DeepEqual(min, expectedMin) {
				t.Errorf("min expected %+v, got %+v", expectedMin, min)
			}
			if !reflect.DeepEqual(max, expectedMax) {
				t.Errorf("max expected %+v, got %+v", expectedMax, max)
			}
		},
	)
}

func TestTileRotate(t *testing.T) {
	t.Run(
		"RightSideUp",
		func(t *testing.T) {
			ti := tile.Tile{
				Width:    4,
				Height:   4,
				Rotation: tile.RotationRightSideUp,
			}

			t.Run(
				"(0, 0)",
				func(t *testing.T) {
					x, y := ti.Rotate(0, 0)
					if x != 0 || y != 3 {
						t.Errorf("Expected (0, 3), got (%d, %d)", x, y)
					}
				},
			)

			t.Run(
				"(0, 3)",
				func(t *testing.T) {
					x, y := ti.Rotate(0, 3)
					if x != 3 || y != 3 {
						t.Errorf("Expected (3, 3), got (%d, %d)", x, y)
					}
				},
			)

			t.Run(
				"(3, 3)",
				func(t *testing.T) {
					x, y := ti.Rotate(3, 3)
					if x != 3 || y != 0 {
						t.Errorf("Expected (3, 0), got (%d, %d)", x, y)
					}
				},
			)
		},
	)
}
