package gg

import (
	"image"
	"image/color"
	"math"
)

func DrawTestPattern(img *image.Gray, width, height int) {
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

func drawColorBars(img *image.Gray, x, y, w, h int) {
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

	for i, _ := range colors {
		startX := x + i*barWidth
		endX := startX + barWidth
		if i == len(colors)-1 {
			endX = x + w // Ensure last bar fills remaining space
		}

		fillRectGray(img, startX, y, endX-startX, h, i+1)
	}
}

func drawPlugeBars(img *image.Gray, x, y, w, h int) {
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

func drawCheckerboard(img *image.Gray, x, y, w, h int) {
	squareSize := 16
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	black := color.RGBA{A: 255}

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

func drawCircularPattern(img *image.Gray, x, y, w, h int) {
	centerX := x + w/2
	centerY := y + h/2
	maxRadius := minimum(w, h) / 2

	// Draw concentric circles
	for radius := 10; radius < maxRadius; radius += 20 {
		drawCircle(img, centerX, centerY, radius, color.RGBA{R: 255, G: 255, B: 255, A: 255})
	}

	// Draw crosshairs
	// Horizontal line
	for i := x; i < x+w; i++ {
		img.Set(i, centerY, color.RGBA{R: 255, G: 255, B: 255, A: 255})
	}
	// Vertical line
	for i := y; i < y+h; i++ {
		img.Set(centerX, i, color.RGBA{R: 255, G: 255, B: 255, A: 255})
	}

	// Draw center dot
	fillCircle(img, centerX, centerY, 5, color.RGBA{R: 255, G: 255, B: 255, A: 255})
}

func fillRect(img *image.Gray, x, y, w, h int, col color.RGBA) {
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			img.Set(x+dx, y+dy, col)
		}
	}
}

func fillRectGray(img *image.Gray, x, y, w, h, level int) {
	count := 0
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			if dx%level == dy%level {
				img.Set(x+dx, y+dy, color.Black)
			} else {
				img.Set(x+dx, y+dy, color.White)
			}
			count++
		}
	}
}

func drawCircle(img *image.Gray, centerX, centerY, radius int, col color.RGBA) {
	for angle := 0.0; angle < 2*math.Pi; angle += 0.01 {
		x := centerX + int(float64(radius)*math.Cos(angle))
		y := centerY + int(float64(radius)*math.Sin(angle))

		if x >= 0 && x < img.Bounds().Dx() && y >= 0 && y < img.Bounds().Dy() {
			img.Set(x, y, col)
		}
	}
}

func fillCircle(img *image.Gray, centerX, centerY, radius int, col color.RGBA) {
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
