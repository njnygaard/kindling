package gg

import (
	"image"
	"image/color"
	"math"
)

type RepeatOp int

const (
// RepeatX RepeatOp = iota
// RepeatY
// RepeatNone
)

type Pattern interface {
	ColorAt(x, y int) color.Color
}

// Solid Pattern
type solidPattern struct {
	color color.Color
}

func (p *solidPattern) ColorAt(_, _ int) color.Color {
	return p.color
}

func NewSolidPattern(color color.Color) Pattern {
	return &solidPattern{color: color}
}

// Surface Pattern
//type surfacePattern struct {
//	im image.Image
//	op RepeatOp
//}

//func (p *surfacePattern) ColorAt(x, y int) color.Color {
//	b := p.im.Bounds()
//	switch p.op {
//	case RepeatX:
//		if y >= b.Dy() {
//			return color.Transparent
//		}
//	case RepeatY:
//		if x >= b.Dx() {
//			return color.Transparent
//		}
//	case RepeatNone:
//		if x >= b.Dx() || y >= b.Dy() {
//			return color.Transparent
//		}
//	}
//	x = x%b.Dx() + b.Min.X
//	y = y%b.Dy() + b.Min.Y
//	return p.im.At(x, y)
//}

//func NewSurfacePattern(im image.Image, op RepeatOp) Pattern {
//	return &surfacePattern{im: im, op: op}
//}

//type patternPainter struct {
//	im   *image.RGBA
//	mask *image.Alpha
//	p    Pattern
//}

// Paint satisfies the Painter interface.
//func (r *patternPainter) Paint(ss []raster.Span, _ bool) {
//	b := r.im.Bounds()
//	for _, s := range ss {
//		if s.Y < b.Min.Y {
//			continue
//		}
//		if s.Y >= b.Max.Y {
//			return
//		}
//		if s.X0 < b.Min.X {
//			s.X0 = b.Min.X
//		}
//		if s.X1 > b.Max.X {
//			s.X1 = b.Max.X
//		}
//		if s.X0 >= s.X1 {
//			continue
//		}
//		const m = 1<<16 - 1
//		y := s.Y - r.im.Rect.Min.Y
//		x0 := s.X0 - r.im.Rect.Min.X
//		// RGBAPainter.Paint() in $GOPATH/src/github.com/golang/freetype/raster/paint.go
//		i0 := (s.Y-r.im.Rect.Min.Y)*r.im.Stride + (s.X0-r.im.Rect.Min.X)*4
//		i1 := i0 + (s.X1-s.X0)*4
//		for i, x := i0, x0; i < i1; i, x = i+4, x+1 {
//			ma := s.Alpha
//			if r.mask != nil {
//				ma = ma * uint32(r.mask.AlphaAt(x, y).A) / 255
//				if ma == 0 {
//					continue
//				}
//			}
//			c := r.p.ColorAt(x, y)
//			cr, cg, cb, ca := c.RGBA()
//			dr := uint32(r.im.Pix[i+0])
//			dg := uint32(r.im.Pix[i+1])
//			db := uint32(r.im.Pix[i+2])
//			da := uint32(r.im.Pix[i+3])
//			a := (m - (ca * ma / m)) * 0x101
//			r.im.Pix[i+0] = uint8((dr*a + cr*ma) / m >> 8)
//			r.im.Pix[i+1] = uint8((dg*a + cg*ma) / m >> 8)
//			r.im.Pix[i+2] = uint8((db*a + cb*ma) / m >> 8)
//			r.im.Pix[i+3] = uint8((da*a + ca*ma) / m >> 8)
//		}
//	}
//}

//func newPatternPainter(im *image.RGBA, mask *image.Alpha, p Pattern) *patternPainter {
//	return &patternPainter{im, mask, p}
//}

//const (
//	width  = 640
//	height = 480
//)

//func main() {
//	// Create a new image
//	img := image.NewRGBA(image.Rect(0, 0, width, height))
//
//	// Draw the classic TV test pattern
//	DrawTestPattern(img)
//
//	// Save to file
//	file, err := os.Create("tv_test_pattern.png")
//	if err != nil {
//		panic(err)
//	}
//	defer file.Close()
//
//	err = png.Encode(file, img)
//	if err != nil {
//		panic(err)
//	}
//
//	println("TV test pattern saved as tv_test_pattern.png")
//}

func DrawTestPattern(img *image.RGBA, width, height int) {
	// Draw color bars at top (75% of height)
	drawColorBars(img, 0, 0, width, height*3/4)

	// Draw bottom section with different patterns
	bottomY := height * 3 / 4
	bottomHeight := height - bottomY

	// Left section: Pluge bars (Picture Line-Up Generation Equipment)
	drawPlugeBars(img, 0, bottomY, width/3, bottomHeight)

	// Middle section: Checkerboard pattern
	drawCheckerboard(img, width/3, bottomY, width/3, bottomHeight)

	// Right section: Circular pattern
	drawCircularPattern(img, 2*width/3, bottomY, width/3, bottomHeight)
}

func drawColorBars(img *image.RGBA, x, y, w, h int) {
	// Standard SMPTE color bars
	colors := []color.RGBA{
		//{191, 191, 191, 255}, // 75% White
		//{191, 191, 0, 255},   // 75% Yellow
		//{0, 191, 191, 255},   // 75% Cyan
		//{0, 191, 0, 255},     // 75% Green
		//{191, 0, 191, 255},   // 75% Magenta
		//{191, 0, 0, 255},     // 75% Red
		//{0, 0, 191, 255},     // 75% Blue
		//{0, 0, 0, 255},       // Black
		{255, 255, 255, 255},
		{223, 223, 223, 255},
		{191, 191, 191, 255},
		{159, 159, 159, 255},
		{127, 127, 127, 255},
		{95, 95, 95, 255},
		{63, 63, 63, 255},
		{31, 31, 31, 255},
		{0, 0, 0, 255},
	}

	barWidth := w / len(colors)

	for i, col := range colors {
		startX := x + i*barWidth
		endX := startX + barWidth
		if i == len(colors)-1 {
			endX = x + w // Ensure last bar fills remaining space
		}

		fillRect(img, startX, y, endX-startX, h, col)
	}
}

func drawPlugeBars(img *image.RGBA, x, y, w, h int) {
	// PLUGE bars for brightness/contrast adjustment
	colors := []color.RGBA{
		{0, 0, 0, 255},       // Super black
		{32, 32, 32, 255},    // Black
		{64, 64, 64, 255},    // Dark gray
		{128, 128, 128, 255}, // Medium gray
	}

	barWidth := w / len(colors)

	for i, col := range colors {
		startX := x + i*barWidth
		endX := startX + barWidth
		if i == len(colors)-1 {
			endX = x + w
		}

		fillRect(img, startX, y, endX-startX, h, col)
	}
}

func drawCheckerboard(img *image.RGBA, x, y, w, h int) {
	squareSize := 16
	white := color.RGBA{255, 255, 255, 255}
	black := color.RGBA{0, 0, 0, 255}

	for row := 0; row < h/squareSize+1; row++ {
		for col := 0; col < w/squareSize+1; col++ {
			startX := x + col*squareSize
			startY := y + row*squareSize

			// Clamp to bounds
			endX := startX + squareSize
			endY := startY + squareSize
			if endX > x+w {
				endX = x + w
			}
			if endY > y+h {
				endY = y + h
			}

			var fillColor color.RGBA
			if (row+col)%2 == 0 {
				fillColor = white
			} else {
				fillColor = black
			}

			fillRect(img, startX, startY, endX-startX, endY-startY, fillColor)
		}
	}
}

func drawCircularPattern(img *image.RGBA, x, y, w, h int) {
	centerX := x + w/2
	centerY := y + h/2
	maxRadius := minimum(w, h) / 2

	// Draw concentric circles
	for radius := 10; radius < maxRadius; radius += 20 {
		drawCircle(img, centerX, centerY, radius, color.RGBA{255, 255, 255, 255})
	}

	// Draw crosshairs
	// Horizontal line
	for i := x; i < x+w; i++ {
		img.Set(i, centerY, color.RGBA{255, 255, 255, 255})
	}
	// Vertical line
	for i := y; i < y+h; i++ {
		img.Set(centerX, i, color.RGBA{255, 255, 255, 255})
	}

	// Draw center dot
	fillCircle(img, centerX, centerY, 5, color.RGBA{255, 255, 255, 255})
}

func fillRect(img *image.RGBA, x, y, w, h int, col color.RGBA) {
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			img.Set(x+dx, y+dy, col)
		}
	}
}

func drawCircle(img *image.RGBA, centerX, centerY, radius int, col color.RGBA) {
	for angle := 0.0; angle < 2*math.Pi; angle += 0.01 {
		x := centerX + int(float64(radius)*math.Cos(angle))
		y := centerY + int(float64(radius)*math.Sin(angle))

		if x >= 0 && x < img.Bounds().Dx() && y >= 0 && y < img.Bounds().Dy() {
			img.Set(x, y, col)
		}
	}
}

func fillCircle(img *image.RGBA, centerX, centerY, radius int, col color.RGBA) {
	for y := centerY - radius; y <= centerY+radius; y++ {
		for x := centerX - radius; x <= centerX+radius; x++ {
			dx := x - centerX
			dy := y - centerY
			if dx*dx+dy*dy <= radius*radius {
				if x >= 0 && x < img.Bounds().Dx() && y >= 0 && y < img.Bounds().Dy() {
					img.Set(x, y, col)
				}
			}
		}
	}
}

func minimum(a, b int) int {
	if a < b {
		return a
	}
	return b
}
