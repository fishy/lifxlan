package main

import (
	"context"
	"image"
	"log"
	"time"

	"github.com/fishy/lifxlan/tile"
)

func getBoard(
	full tile.ColorBoard,
	board tile.Board,
	horizontal bool,
	step int,
) (cb tile.ColorBoard, last bool) {
	cb = tile.MakeColorBoard(board.Width(), board.Height())

	var fullSize, boardSize int
	if horizontal {
		fullSize = len(full)
		boardSize = board.Width()
	} else {
		fullSize = len(full[0])
		boardSize = board.Height()
	}

	offset := 0
	start := 0
	end := step - start
	if end > boardSize {
		start += (end - boardSize)
	}
	if end > fullSize {
		offset = end - fullSize
		end = fullSize
	}
	last = (end <= start)

	for i := start; i < end; i++ {
		ci := i - start + offset
		fi := fullSize - end + i - start
		if horizontal {
			for j := 0; j < board.Height(); j++ {
				cb[ci][j] = full[fi][j]
			}
		} else {
			for j := 0; j < board.Width(); j++ {
				cb[j][ci] = full[j][fi]
			}
		}
	}

	return
}

func draw(td tile.Device, img image.Image) {
	var step int

	full, horizontal := resizeImage(img, td)

	conn, err := td.Dial()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	var origCB tile.ColorBoard
	if !*loop {
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), *drawTimeout)
			defer cancel()
			var err error
			origCB, err = td.GetColors(ctx, conn)
			if err != nil {
				log.Fatalf("Cannot get the current colors on %v: %v", td, err)
			}
		}()
	}

	for range time.Tick(*interval) {
		step++
		log.Printf("Step %d...", step)
		cb, last := getBoard(full, td, horizontal, step)
		start := time.Now()
		if err := func() error {
			ctx, cancel := context.WithTimeout(context.Background(), *drawTimeout)
			defer cancel()
			return td.SetColors(ctx, conn, cb, 0, !*noack)
		}(); err != nil {
			log.Printf("Failed to set colors: %v", err)
			if *noskip {
				step--
				continue
			}
		} else {
			log.Printf("SetColors took %v", time.Since(start))
		}
		if last {
			if *loop {
				step = 0
				log.Print("Finished. Resetting...")
			} else {
				break
			}
		}
	}

	for {
		if err := func() error {
			ctx, cancel := context.WithTimeout(context.Background(), *drawTimeout)
			defer cancel()
			return td.SetColors(ctx, conn, origCB, 0, true)
		}(); err != nil {
			log.Printf("Failed to set original colors, retrying... %v", err)
		} else {
			return
		}
	}
}
