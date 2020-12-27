package tile_test

import (
	"reflect"
	"testing"

	"go.yhsif.com/lifxlan/tile"
)

func id(i, x, y int) *tile.IndexData {
	t := &tile.IndexData{}
	t.Index = i
	t.X = x
	t.Y = y
	return t
}

func bd(
	x, y int,
	data [][]*tile.IndexData,
	reverse [][][]tile.Coordinate,
) tile.BoardData {
	b := tile.BoardData{
		Data:        data,
		ReverseData: reverse,
	}
	b.X = x
	b.Y = y
	return b
}

func TestParse(t *testing.T) {
	t.Run(
		"Empty",
		func(t *testing.T) {
			var b tile.BoardData
			if b.X != 0 || b.Y != 0 {
				t.Errorf("BoardData %+v is not empty.", b)
			}
		},
	)

	t.Run(
		"Single",
		func(t *testing.T) {
			ti := &tile.Tile{
				UserX:    2,
				UserY:    -2,
				Width:    2,
				Height:   2,
				Rotation: rotation,
			}
			expected := bd(
				2,
				2,
				[][]*tile.IndexData{
					{id(0, 1, 0), id(0, 0, 0)},
					{id(0, 1, 1), id(0, 0, 1)},
				},
				[][][]tile.Coordinate{
					{
						{c(0, 1), c(1, 1)},
						{c(0, 0), c(1, 0)},
					},
				},
			)
			b := tile.ParseBoard([]*tile.Tile{ti})
			if !reflect.DeepEqual(b, expected) {
				t.Errorf("BoardData %+v expected to be %+v", b, expected)
			}
		},
	)

	t.Run(
		"Double",
		func(t *testing.T) {
			t1 := &tile.Tile{
				UserX:    1,
				UserY:    1,
				Width:    1,
				Height:   1,
				Rotation: rotation,
			}
			t2 := &tile.Tile{
				UserX:    -1,
				UserY:    -1,
				Width:    1,
				Height:   1,
				Rotation: rotation,
			}
			expected := bd(
				3,
				3,
				[][]*tile.IndexData{
					{id(1, 0, 0), nil, nil},
					{nil, nil, nil},
					{nil, nil, id(0, 0, 0)},
				},
				[][][]tile.Coordinate{
					{
						{c(2, 2)},
					},
					{
						{c(0, 0)},
					},
				},
			)
			b := tile.ParseBoard([]*tile.Tile{t1, t2})
			if !reflect.DeepEqual(b, expected) {
				t.Errorf("BoardData %+v expected to be %+v", b, expected)
			}
		},
	)
}
