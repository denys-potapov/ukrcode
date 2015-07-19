package main

import (
	// "fmt"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/datamatrix"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

const (
	any int = iota
	set
	skp
	emp
)

var (
	white = color.White
	black = color.RGBA{0x21, 0x21, 0x21, 0xff}
	red   = color.RGBA{0xf4, 0x43, 0x36, 0xff}
)

type Bitmap struct {
	data   []bool
	width  int
	height int
}

func (b *Bitmap) valid(x, y int) bool {
	return (x >= 0) && (y >= 0) && (x < b.width) && (y < b.height)
}

func (b *Bitmap) Set(x, y int, value bool) {
	if b.valid(x, y) {
		b.data[y*b.width+x] = value
	}
}

func (b *Bitmap) Get(x, y int) bool {
	if b.valid(x, y) {
		return b.data[y*b.width+x]
	}
	return false
}

func Match(b *Bitmap, p [][]int, sx, sy int) bool {
	for y := 0; y < len(p); y++ {
		for x := 0; x < len(p[y]); x++ {
			if p[y][x] == any {
				continue
			}
			if (p[y][x] == set || p[y][x] == skp) && (!b.Get(sx + x, sy + y)) {
				return false
			}
			if p[y][x] == emp && (b.Get(sx + x, sy + y)) {
				return false
			}
		}
	}

	return true
}

func Draw(i *image.RGBA, b *Bitmap, p [][]int, sx, sy, scale int) {
	for y := 0; y < len(p); y++ {
		for x := 0; x < len(p[y]); x++ {
			if p[y][x] != set && p[y][x] != skp {
				continue
			}
			b.Set(sx + x, sy + y, false)

			color := red;
			if p[y][x] == skp {
				color = black
			}
			pixel := image.Rect((sx + x)*scale, (sy + y)*scale, (sx + x+1)*scale, (sy + y+1)*scale)
			draw.Draw(i, pixel, &image.Uniform{color}, image.ZP, draw.Src)
		}
	}
}

func NewBitmap(code barcode.Barcode) *Bitmap {
	dim := code.Bounds().Max

	b := &Bitmap{data: make([]bool, dim.X*dim.Y), height: dim.Y, width: dim.X}
	for y := 0; y < dim.Y; y++ {
		for x := 0; x < dim.X; x++ {
			if code.At(x, y) == color.Black {
				b.Set(x, y, true)
			}
		}
	}

	return b
}

func hor(code barcode.Barcode, scale int) image.Image {
	b := NewBitmap(code)
	img := image.NewRGBA(image.Rect(0, 0, b.height*scale, b.width*scale))
	draw.Draw(img, img.Bounds(), &image.Uniform{white}, image.ZP, draw.Src)

	patterns := [][][]int {
		{	{any, set, any},
			{set, skp, set},
			{any, set, any}}}

	for _, pattern := range patterns {
		for y := 0; y < b.height; y++ {
			for x := 0; x < b.width; x++ {
				if Match(b, pattern, x, y) {
					Draw(img, b, pattern, x, y, scale)
				}
			}
		}	
	}
	

	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			if b.Get(x, y) {
				pixel := image.Rect(x*scale, y*scale, (x+1)*scale, (y+1)*scale)
				draw.Draw(img, pixel, &image.Uniform{black}, image.ZP, draw.Src)
			}
		}
	}

	return img
}

func main() {
	text := "http://denyspotapov.com/"
	code, _ := datamatrix.Encode(text)
	png.Encode(os.Stdout, hor(code, 10))
}
