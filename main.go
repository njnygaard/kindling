package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

func main() {
	width := 758
	height := 1024

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewGray16(image.Rectangle{upLeft, lowRight})

	// Colors are defined by Red, Green, Blue, Alpha uint8 values.
	// g := color.Gray16{100}

	// Set color for each pixel.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			switch {
			// Center Circle
			case isCircle(width/2, height/2, 50, x, y):
				img.Set(x, y, color.Gray16{0x7777})
			// Upper Left Circle
			case isCircle(width/4, height/4, 50, x, y):
				img.Set(x, y, color.Gray16{0x7777})
			// Upper Right Circle
			case isCircle((3*width)/4, height/4, 50, x, y):
				img.Set(x, y, color.Gray16{0x0000})
			// Lower Left Circle
			case isCircle(width/4, (3*height)/4, 50, x, y):
				img.Set(x, y, color.Gray16{0xffff})
			// Lower Right Circle
			case isCircle((3*width)/4, (3*height)/4, 50, x, y):
				img.Set(x, y, color.Gray16{0x7777})
			// First Quadrant
			case x < width/2 && y < height/2 && x%2 == 0:
				img.Set(x, y, color.Black)
			case x < width/2 && y < height/2 && x%2 == 1:
				img.Set(x, y, color.White)
			// Second Quadrant
			case x > width/2 && y < height/2 && x%2 == 0:
				img.Set(x, y, color.White)
			case x > width/2 && y < height/2 && x%2 == 1:
				img.Set(x, y, color.White)
			// Third Quadrant
			case x < width/2 && y > height/2 && x%2 == 0:
				img.Set(x, y, color.Black)
			case x < width/2 && y > height/2 && x%2 == 1:
				img.Set(x, y, color.Black)
			// Fourth Quadrant
			case x > width/2 && y > height/2 && x%2 == 0:
				img.Set(x, y, color.Black)
			case x > width/2 && y > height/2 && x%2 == 1:
				img.Set(x, y, color.White)
			default:
				img.Set(x, y, color.White)
			}
		}
	}

	// Encode as PNG.
	f, _ := os.Create("image.png")
	png.Encode(f, img)
}

func isCircle(Xo int, Yo int, r int, x int, y int) (inside bool) {

	xc := math.Pow(float64(x-Xo), 2)
	yc := math.Pow(float64(y-Yo), 2)
	rc := math.Pow(float64(r), 2)

	return (xc + yc) < rc

}
