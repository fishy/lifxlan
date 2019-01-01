package tile_test

import (
	"reflect"
	"testing"

	"github.com/fishy/lifxlan/tile"
)

// To grately reduce typing in test code.
type c = tile.Coordinate

func cc(x, y int) c {
	return c{
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
			expectedMin := cc(0, 0)
			expectedMax := cc(8, 8)
			expectedCs := [][]c{
				[]c{
					cc(0, 7),
					cc(1, 7),
					cc(2, 7),
					cc(3, 7),
					cc(4, 7),
					cc(5, 7),
					cc(6, 7),
					cc(7, 7),
				},
				[]c{
					cc(0, 6),
					cc(1, 6),
					cc(2, 6),
					cc(3, 6),
					cc(4, 6),
					cc(5, 6),
					cc(6, 6),
					cc(7, 6),
				},
				[]c{
					cc(0, 5),
					cc(1, 5),
					cc(2, 5),
					cc(3, 5),
					cc(4, 5),
					cc(5, 5),
					cc(6, 5),
					cc(7, 5),
				},
				[]c{
					cc(0, 4),
					cc(1, 4),
					cc(2, 4),
					cc(3, 4),
					cc(4, 4),
					cc(5, 4),
					cc(6, 4),
					cc(7, 4),
				},
				[]c{
					cc(0, 3),
					cc(1, 3),
					cc(2, 3),
					cc(3, 3),
					cc(4, 3),
					cc(5, 3),
					cc(6, 3),
					cc(7, 3),
				},
				[]c{
					cc(0, 2),
					cc(1, 2),
					cc(2, 2),
					cc(3, 2),
					cc(4, 2),
					cc(5, 2),
					cc(6, 2),
					cc(7, 2),
				},
				[]c{
					cc(0, 1),
					cc(1, 1),
					cc(2, 1),
					cc(3, 1),
					cc(4, 1),
					cc(5, 1),
					cc(6, 1),
					cc(7, 1),
				},
				[]c{
					cc(0, 0),
					cc(1, 0),
					cc(2, 0),
					cc(3, 0),
					cc(4, 0),
					cc(5, 0),
					cc(6, 0),
					cc(7, 0),
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
			expectedMin := cc(-1, 1)
			expectedMax := cc(1, 3)
			expectedCs := [][]c{
				[]c{
					cc(-1, 2),
					cc(0, 2),
				},
				[]c{
					cc(-1, 1),
					cc(0, 1),
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
			expectedMin := cc(4, -4)
			expectedMax := cc(6, -2)
			expectedCs := [][]c{
				[]c{
					cc(4, -3),
					cc(5, -3),
				},
				[]c{
					cc(4, -4),
					cc(5, -4),
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
