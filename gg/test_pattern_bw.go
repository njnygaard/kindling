package gg

import (
	"image"
	"image/color"
	"math"
)

func DrawTestPatternBW(img *image.Gray, width, height int) {
	// Draw color bars at top (75% of height)
	drawColorBarsBW(img, 0, 0, width, height*3/4)

	// Draw bottom section with different patterns
	bottomY := height * 3 / 4
	bottomHeight := height - bottomY

	// Left section: Pluge bars (Picture Line-Up Generation Equipment)
	drawPlugeBarsBW(img, 0, bottomY, width/3, bottomHeight)

	// Middle section: Checkerboard pattern (already black/white)
	drawCheckerboard(img, width/3, bottomY, width/3, bottomHeight)

	// Right section: Circular pattern
	drawCircularPatternBW(img, 2*width/3, bottomY, width/3, bottomHeight)
}

func drawColorBarsBW(img *image.Gray, x, y, w, h int) {
	// Gray levels for the bars (from white to black)
	grayLevels := []int{
		255, // White
		223, // Light gray
		191, // Medium light gray
		159, // Medium gray
		127, // Medium dark gray
		95,  // Dark gray
		63,  // Very dark gray
		31,  // Almost black
		0,   // Black
	}

	barWidth := w / len(grayLevels)

	for i, grayLevel := range grayLevels {
		startX := x + i*barWidth
		endX := startX + barWidth
		if i == len(grayLevels)-1 {
			endX = x + w // Ensure last bar fills remaining space
		}

		fillRectBW(img, startX, y, endX-startX, h, grayLevel)
	}
}

func drawPlugeBarsBW(img *image.Gray, x, y, w, h int) {
	// PLUGE bars for brightness/contrast adjustment
	grayLevels := []int{
		0,   // Super black
		32,  // Black
		64,  // Dark gray
		128, // Medium gray
	}

	barWidth := w / len(grayLevels)

	for i, grayLevel := range grayLevels {
		startX := x + i*barWidth
		endX := startX + barWidth
		if i == len(grayLevels)-1 {
			endX = x + w
		}

		fillRectBW(img, startX, y, endX-startX, h, grayLevel)
	}
}

func drawCircularPatternBW(img *image.Gray, x, y, w, h int) {
	centerX := x + w/2
	centerY := y + h/2
	maxRadius := minimum(w, h) / 2

	// Draw concentric circles with dithering
	for radius := 10; radius < maxRadius; radius += 20 {
		drawCircleBW(img, centerX, centerY, radius, 255) // White circles
	}

	// Draw crosshairs (white lines)
	// Horizontal line
	for i := x; i < x+w; i++ {
		img.Set(i, centerY, color.White)
	}
	// Vertical line
	for i := y; i < y+h; i++ {
		img.Set(centerX, i, color.White)
	}

	// Draw center dot (white)
	fillCircleBW(img, centerX, centerY, 5, 255)
}

func fillRectBW(img *image.Gray, x, y, w, h int, grayLevel int) {
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			pixelColor := DitherPixel(grayLevel, x+dx, y+dy)
			img.Set(x+dx, y+dy, pixelColor)
		}
	}
}

func drawCircleBW(img *image.Gray, centerX, centerY, radius int, grayLevel int) {
	for angle := 0.0; angle < 2*math.Pi; angle += 0.01 {
		x := centerX + int(float64(radius)*math.Cos(angle))
		y := centerY + int(float64(radius)*math.Sin(angle))

		if x >= 0 && x < img.Bounds().Dx() && y >= 0 && y < img.Bounds().Dy() {
			pixelColor := DitherPixel(grayLevel, x, y)
			img.Set(x, y, pixelColor)
		}
	}
}

func fillCircleBW(img *image.Gray, centerX, centerY, radius int, grayLevel int) {
	for y := centerY - radius; y <= centerY+radius; y++ {
		for x := centerX - radius; x <= centerX+radius; x++ {
			dx := x - centerX
			dy := y - centerY
			if dx*dx+dy*dy <= radius*radius {
				if x >= 0 && x < img.Bounds().Dx() && y >= 0 && y < img.Bounds().Dy() {
					pixelColor := DitherPixel(grayLevel, x, y)
					img.Set(x, y, pixelColor)
				}
			}
		}
	}
}
