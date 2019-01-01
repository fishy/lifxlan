package main

import (
	"context"
	"crypto/rand"
	"log"
	"time"

	"github.com/fishy/lifxlan"
	"github.com/fishy/lifxlan/tile"
)

var neighbours = []tile.Coordinate{
	tile.Coordinate{
		X: -1,
		Y: -1,
	},
	tile.Coordinate{
		X: -1,
		Y: 0,
	},
	tile.Coordinate{
		X: 1,
		Y: 0,
	},
	tile.Coordinate{
		X: 0,
		Y: -1,
	},
	tile.Coordinate{
		X: 0,
		Y: 1,
	},
	tile.Coordinate{
		X: 1,
		Y: -1,
	},
	tile.Coordinate{
		X: 1,
		Y: 0,
	},
	tile.Coordinate{
		X: 1,
		Y: -1,
	},
}

func draw(td tile.Device) {
	color := lifxlan.FromColor(origColor, uint16(*kelvin))
	width := td.Width()
	height := td.Height()
	var step int

	initBoard := func() [][]bool {
		step = 0
		board := make([][]bool, width)
		for i := range board {
			board[i] = make([]bool, height)
		}

		// Generate random initial board
		buf := make([]byte, 1)
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				if !td.OnTile(x, y) {
					continue
				}
				if _, err := rand.Read(buf); err != nil {
					log.Fatal(err)
				}
				// 50% light would usually be too bright, so make it 25% instead.
				board[x][y] = buf[0]%4 == 1
			}
		}

		return board
	}

	board := initBoard()

	var empty bool
	drawBoard := func() {
		colors := tile.MakeColorBoard(width, height)
		empty = true
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				if board[x][y] {
					empty = false
					colors[x][y] = color
				}
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), *drawTimeout)
		defer cancel()
		start := time.Now()
		if err := td.SetColors(ctx, nil, colors, *ack); err != nil {
			log.Printf("Failed to set colors: %v", err)
		} else {
			log.Printf("SetColors took %v", time.Since(start))
		}
	}

	drawBoard()

	countNeighbours := func(x, y int) (n int) {
		for _, nei := range neighbours {
			newX := x + nei.X
			if newX < 0 || newX >= width {
				continue
			}
			newY := y + nei.Y
			if newY < 0 || newY >= height {
				continue
			}
			if board[newX][newY] {
				n++
			}
		}
		return
	}

	counts := make([][]int, width)
	for i := range counts {
		counts[i] = make([]int, height)
	}

	evolve := func() {
		start := time.Now()
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				if td.OnTile(x, y) {
					counts[x][y] = countNeighbours(x, y)
				}
			}
		}
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				if td.OnTile(x, y) {
					switch counts[x][y] {
					case 0, 1, 4, 5, 6, 7, 8:
						board[x][y] = false
					case 3:
						board[x][y] = true
					}
				}
			}
		}
		log.Printf("evolve took %v", time.Since(start))
	}

	for range time.Tick(*interval) {
		if empty {
			log.Print("Board empty, resetting board...")
			board = initBoard()
		} else {
			step++
			if *reset > 0 && step >= *reset {
				log.Print("Resetting board...")
				board = initBoard()
			} else {
				log.Printf("Step %d...", step)
				evolve()
			}
		}
		drawBoard()
	}
}
