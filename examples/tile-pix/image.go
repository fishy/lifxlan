package main

import (
	"image"
	"image/color"
	"io"
	"math"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/fishy/lifxlan"
	"github.com/fishy/lifxlan/tile"
)

func readImage(r io.Reader) (img image.Image, err error) {
	img, _, err = image.Decode(r)
	return
}

func resizeImage(img image.Image, board tile.Board) (cb tile.ColorBoard, horizontal bool) {
	bounds := img.Bounds()
	origWidth := bounds.Max.X - bounds.Min.X
	origHeight := bounds.Max.Y - bounds.Min.Y

	width := board.Width()
	height := origHeight * width / origWidth
	if height < board.Height() {
		horizontal = true
		height = board.Height()
		width = origWidth * height / origHeight
	}
	ratio := float64(origWidth) / float64(width)
	n := int(math.Ceil(ratio))
	colors := make([]color.Color, n*n)

	cb = tile.MakeColorBoard(width, height)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			baseX := int(math.Round(float64(x) * ratio))
			if baseX < 0 {
				baseX = 0
			}
			if baseX+n > origWidth {
				baseX -= (baseX + n - origWidth)
			}

			baseY := int(math.Round(float64(y) * ratio))
			if baseY < 0 {
				baseY = 0
			}
			if baseY+n > origHeight {
				baseY -= (baseY + n - origHeight)
			}
			for i := 0; i < n; i++ {
				for j := 0; j < n; j++ {
					index := i*n + j
					xx := baseX + bounds.Min.X + i
					yy := baseY + bounds.Min.Y + j
					colors[index] = img.At(xx, yy)
				}
			}

			// Flip the y axis to fit our coordinates system.
			cb[x][height-1-y] = averageColor(colors)
		}
	}
	return
}

type rgb struct {
	r, g, b uint32
}

var _ color.Color = (*rgb)(nil)

func (r rgb) RGBA() (uint32, uint32, uint32, uint32) {
	return r.r, r.g, r.b, 0
}

func averageColor(colors []color.Color) *lifxlan.Color {
	var rr, gg, bb float64
	for _, c := range colors {
		r, g, b, _ := c.RGBA()
		rr += float64(r) * float64(r)
		gg += float64(g) * float64(g)
		bb += float64(b) * float64(b)
	}
	rr = math.Sqrt(rr / float64(len(colors)))
	gg = math.Sqrt(gg / float64(len(colors)))
	bb = math.Sqrt(bb / float64(len(colors)))
	return lifxlan.FromColor(
		&rgb{
			r: uint32(rr),
			g: uint32(gg),
			b: uint32(bb),
		},
		uint16(*kelvin),
	)
}
