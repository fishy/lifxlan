package tile_test

import (
	"context"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
	"time"

	"github.com/fishy/lifxlan"
	"github.com/fishy/lifxlan/mock"
	"github.com/fishy/lifxlan/tile"
)

func TestColorBoard(t *testing.T) {
	t.Run(
		"Empty",
		func(t *testing.T) {
			var n int
			var cb tile.ColorBoard

			// Seed the pseudo-random generator
			now := time.Now()
			rander := rand.New(rand.NewSource(now.Unix() + int64(now.Nanosecond())))

			f := func() bool {
				n++
				x := rander.Int()
				y := rander.Int()
				color := cb.GetColor(x, y)
				if color != nil {
					t.Logf("GetColor(%d, %d) returned non-nil color: %v", x, y, *color)
					return false
				}
				return true
			}
			if err := quick.Check(f, nil); err != nil {
				t.Error(err)
			}
			t.Logf("quick did %d checks", n)
		},
	)

	t.Run(
		"NonEmpty",
		func(t *testing.T) {
			makeBoard := func(t *testing.T, width, height int) tile.ColorBoard {
				cb := tile.MakeColorBoard(width, height)
				if len(cb) != width {
					t.Fatalf("MakeColorBoard returned with width %d", len(cb))
				}
				for i, row := range cb {
					if len(row) != height {
						t.Errorf(
							"MakeColorBoard returned row %d with height %d",
							i,
							len(row),
						)
					}
				}
				if t.Failed() {
					t.FailNow()
				}
				return cb
			}

			inPoints := func(x, y int, points []tile.Coordinate) bool {
				for _, p := range points {
					if p.X == x && p.Y == y {
						return true
					}
				}
				return false
			}

			t.Run(
				"4x6",
				func(t *testing.T) {
					const (
						width  = 4
						height = 6

						minX = -10
						maxX = 10
						minY = -10
						maxY = 10
					)
					color := &lifxlan.ColorBlack
					cb := makeBoard(t, width, height)

					points := []tile.Coordinate{
						{
							X: 1,
							Y: 2,
						},
						{
							X: 3,
							Y: 0,
						},
					}
					for _, p := range points {
						cb[p.X][p.Y] = color
					}

					for x := minX; x < maxX; x++ {
						for y := minY; y < maxY; y++ {
							var expected *lifxlan.Color
							if inPoints(x, y, points) {
								expected = color
							}
							got := cb.GetColor(x, y)
							if got != expected {
								t.Errorf(
									"GetColor(%d, %d) expected %+v, got %+v",
									x,
									y,
									expected,
									got,
								)
							}
						}
					}
				},
			)

			t.Run(
				"100x100",
				func(t *testing.T) {
					const (
						width  = 100
						height = 100

						minX = -100
						maxX = 150
						minY = -100
						maxY = 150
					)
					color := &lifxlan.ColorBlack
					cb := makeBoard(t, width, height)

					// Some random coordinates
					points := []tile.Coordinate{
						{
							X: 12,
							Y: 17,
						},
						{
							X: 54,
							Y: 23,
						},
						{
							X: 12,
							Y: 87,
						},
						{
							X: 24,
							Y: 23,
						},
						{
							X: 35,
							Y: 18,
						},
					}
					for _, p := range points {
						cb[p.X][p.Y] = color
					}

					for x := minX; x < maxX; x++ {
						for y := minY; y < maxY; y++ {
							var expected *lifxlan.Color
							if inPoints(x, y, points) {
								expected = color
							}
							got := cb.GetColor(x, y)
							if got != expected {
								t.Errorf(
									"GetColor(%d, %d) expected %+v, got %+v",
									x,
									y,
									expected,
									got,
								)
							}
						}
					}
				},
			)
		},
	)
}

// This example demonstrates how to make a ColorBoard of random colors for a board.
func ExampleMakeColorBoard() {
	// Variables should be initialized properly in real code.
	var (
		board tile.Board

		// Should return a random color.
		colorGenerator func() *lifxlan.Color
	)

	// This makes a full size ColorBoard.
	// If you only need to draw partially and leave the rest of the board black,
	// you can use smaller width/height values that's enough to cover the area you
	// want to draw.
	cb := tile.MakeColorBoard(board.Width(), board.Height())
	for x := 0; x < board.Width(); x++ {
		for y := 0; y < board.Height(); y++ {
			if !board.OnTile(x, y) {
				// This coordinate is not on any tile so there's no need to draw it.
				continue
			}
			cb[x][y] = colorGenerator()
		}
	}
}

func TestColorsAPIs(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	const timeout = time.Millisecond * 200

	service, device := mock.StartService(t)
	defer service.Stop()
	rawTile1 := tile.RawTileDevice{
		Width:  8,
		Height: 8,
	}
	rawTile2 := tile.RawTileDevice{
		UserX:  1,
		Width:  8,
		Height: 8,
	}
	rawChain := &tile.RawStateDeviceChainPayload{
		TotalCount: 2,
	}
	rawChain.TileDevices[0] = rawTile1
	rawChain.TileDevices[1] = rawTile2
	service.RawStateDeviceChainPayload = rawChain

	td, err := func() (tile.Device, error) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		return tile.Wrap(ctx, device, false)
	}()
	if err != nil {
		t.Fatal(err)
	}
	if td == nil {
		t.Fatal("Can't mock tile device.")
	}

	stateColor1 := tile.RawStateTileState64Payload{
		TileIndex: 0,
		Width:     8,
	}
	for x, row := range stateColor1.Colors {
		for y := range row {
			stateColor1.Colors[x][y] = lifxlan.ColorBlack
		}
	}
	stateColor2 := stateColor1
	stateColor2.TileIndex = 1

	t.Run(
		"GetColors",
		func(t *testing.T) {
			t.Run(
				"NotEnough",
				func(t *testing.T) {
					service.RawStateTileState64Payloads = []*tile.RawStateTileState64Payload{
						&stateColor1,
					}

					ctx, cancel := context.WithTimeout(context.Background(), timeout)
					defer cancel()

					if _, err := td.GetColors(ctx, nil); err == nil {
						t.Error("Expected error when not enough tiles returned, got nil")
					}
				},
			)

			t.Run(
				"Normal",
				func(t *testing.T) {
					service.RawStateTileState64Payloads = []*tile.RawStateTileState64Payload{
						&stateColor1,
						&stateColor2,
					}

					ctx, cancel := context.WithTimeout(context.Background(), timeout)
					defer cancel()

					cb, err := td.GetColors(ctx, nil)
					if err != nil {
						t.Fatal(err)
					}

					for x := 0; x < 16; x++ {
						for y := 0; y < 8; y++ {
							if !reflect.DeepEqual(*cb.GetColor(x, y), lifxlan.ColorBlack) {
								t.Errorf("Got color %+v at (%d, %d)", cb.GetColor(x, y), x, y)
							}
						}
					}
				},
			)
		},
	)

	t.Run(
		"SetColors",
		func(t *testing.T) {
			t.Run(
				"NotEnoughAcks",
				func(t *testing.T) {
					service.AcksToDrop = 1

					ctx, cancel := context.WithTimeout(context.Background(), timeout)
					defer cancel()

					if err := td.SetColors(ctx, nil, nil, 0, true); err == nil {
						t.Error("Expected error when not enough acks returned, got nil")
					}
				},
			)

			t.Run(
				"Normal",
				func(t *testing.T) {
					service.AcksToDrop = 0

					ctx, cancel := context.WithTimeout(context.Background(), timeout)
					defer cancel()

					if err := td.SetColors(ctx, nil, nil, 0, true); err != nil {
						t.Error(err)
					}
				},
			)
		},
	)
}
