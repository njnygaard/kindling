package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"

	owm "github.com/briandowns/openweathermap"
	"github.com/fogleman/gg"
)

const (
	API_KEY = "b40995f7e1911e427c0700778e542369"
)

// func printWeather() {
// 	w, err := owm.NewCurrent("F", "en", API_KEY)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	w.CurrentByName("Chandler")

// 	fmt.Printf("Sunrise: %d\n", w.Sys.Sunrise)
// 	fmt.Printf("Sunrise: %d\n", w.Sys.Sunset)
// 	fmt.Printf("Weather: %s\n", w.Weather[0].Main)
// 	fmt.Printf("Weather Description: %s\n", w.Weather[0].Description)
// 	fmt.Printf("Weather Icon: %s\n", w.Weather[0].Icon)
// 	fmt.Printf("Temp: %f\n", w.Main.Temp)
// 	fmt.Printf("Min: %f\n", w.Main.TempMin)
// 	fmt.Printf("Max: %f\n", w.Main.TempMax)
// 	fmt.Printf("Feels Like: %f\n", w.Main.FeelsLike)
// 	fmt.Printf("Humidity: %d\n", w.Main.Humidity)
// 	fmt.Printf("Wind Speed: %f\n", w.Wind.Speed)
// 	fmt.Printf("Humidity: %f\n", w.Wind.Deg)

// }

func main() {
	width := 758
	height := 1024

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	// printWeather()

	// DO NOT USE image.Gray16
	// The kindle will subtly break
	img := image.NewGray(image.Rectangle{upLeft, lowRight})

	drawTestPattern(img, width, height)
	//drawGrid(img, width, height)

	w, err := owm.NewCurrent("F", "en", API_KEY)
	if err != nil {
		log.Fatalln(err)
	}

	w.CurrentByName("Chandler")

	// const S = 1024
	dc := gg.NewContext(width, height)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	if err := dc.LoadFontFace("impact.ttf", 48); err != nil {
		panic(err)
	}
	dc.SetRGB(0, 0, 0)
	s := w.Weather[0].Main
	// n := 6 // "stroke" size
	// for dy := -n; dy <= n; dy++ {
	// 	for dx := -n; dx <= n; dx++ {
	// 		if dx*dx+dy*dy >= n*n {
	// 			// give it rounded corners
	// 			continue
	// 		}
	// 		x := float64(width)/2 + float64(dx)
	// 		y := float64(height)/2 + float64(dy)
	// 		dc.DrawStringAnchored(s, x, y, 0.5, 0.5)
	// 	}
	// }
	// dc.SetRGB(1, 1, 1)
	dc.DrawStringAnchored(s, float64(width)/2, float64(height)/4, 0.5, 0.5)

	s = fmt.Sprintf("%.0f°F", w.Main.Temp)
	dc.DrawStringAnchored(s, float64(width)/2, float64(3*height)/4, 0.5, 0.5)
	dc.SavePNG("out.png")

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

func isGrid(x int, y int, width int, height int, stroke int) bool {
	// if (x < stroke || x > (width-stroke)) || (y < stroke || y > (height-stroke)) {
	// 	return true
	// } else if ((x)%(106)) == 0 || ((x+stroke)%106) == 0 {
	// 	return true
	// }

	if ((y > 148) && (y < 876) && (x == 0 || x == 1 || x == 2 || x == 108 || x == 109 || x == 110 || x == 216 || x == 217 || x == 218 || x == 324 || x == 325 || x == 326 || x == 432 || x == 433 || x == 434 || x == 540 || x == 541 || x == 542 || x == 648 || x == 649 || x == 650 || x == 756 || x == 757 || x == 758)) ||
		(y == 0 || y == 1 || y == 2 || y == 146 || y == 147 || y == 148 || y == 292 || y == 293 || y == 294 || y == 438 || y == 439 || y == 440 || y == 584 || y == 585 || y == 586 || y == 730 || y == 731 || y == 732 || y == 876 || y == 877 || y == 878 || y == 1022 || y == 1023 || y == 1024) {
		return true
	}

	return false
}

func drawGrid(img *image.Gray, width int, height int) {

	stroke := 2
	// Set color for each pixel.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			switch {

			// Border
			case isGrid(x, y, width, height, stroke):
				img.Set(x, y, color.Black)

			default:
				img.Set(x, y, color.White)
			}
		}
	}

}

func drawTestPattern(img *image.Gray, width int, height int) {

	// Set color for each pixel.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			switch {

			// Center Circle
			case isCircle((1*width)/2, (1*height)/2, 50, x, y):
				img.Set(x, y, color.Gray{0x77})

			// Upper Left Circles
			case isCircle((1*width)/8, (1*height)/8, 50, x, y):
				img.Set(x, y, color.Black)
			case isCircle((3*width)/8, (1*height)/8, 50, x, y):
				img.Set(x, y, color.Gray{0x33})
			case isCircle((1*width)/8, (3*height)/8, 50, x, y):
				img.Set(x, y, color.Gray{0xAA})
			case isCircle((3*width)/8, (3*height)/8, 50, x, y):
				img.Set(x, y, color.White)
			case isCircle((2*width)/8, (2*height)/8, 50, x, y):
				img.Set(x, y, color.Gray{0x77})

			// Upper Right Circle
			case isCircle((5*width)/8, (1*height)/8, 50, x, y):
				img.Set(x, y, color.Black)
			case isCircle((7*width)/8, (1*height)/8, 50, x, y):
				img.Set(x, y, color.Gray{0x33})
			case isCircle((5*width)/8, (3*height)/8, 50, x, y):
				img.Set(x, y, color.Gray{0xAA})
			case isCircle((7*width)/8, (3*height)/8, 50, x, y):
				img.Set(x, y, color.White)
			case isCircle((6*width)/8, (2*height)/8, 50, x, y):
				img.Set(x, y, color.Gray{0x77})

			// Lower Left Circle
			case isCircle((1*width)/8, (5*height)/8, 50, x, y):
				img.Set(x, y, color.Black)
			case isCircle((3*width)/8, (5*height)/8, 50, x, y):
				img.Set(x, y, color.Gray{0x33})
			case isCircle((1*width)/8, (7*height)/8, 50, x, y):
				img.Set(x, y, color.Gray{0xAA})
			case isCircle((3*width)/8, (7*height)/8, 50, x, y):
				img.Set(x, y, color.White)
			case isCircle((2*width)/8, (6*height)/8, 50, x, y):
				img.Set(x, y, color.Gray{0x77})

			// Lower Right Circle
			case isCircle((5*width)/8, (5*height)/8, 50, x, y):
				img.Set(x, y, color.Black)
			case isCircle((7*width)/8, (5*height)/8, 50, x, y):
				img.Set(x, y, color.Gray{0x33})
			case isCircle((5*width)/8, (7*height)/8, 50, x, y):
				img.Set(x, y, color.Gray{0xAA})
			case isCircle((7*width)/8, (7*height)/8, 50, x, y):
				img.Set(x, y, color.White)
			case isCircle((6*width)/8, (6*height)/8, 50, x, y):
				img.Set(x, y, color.Gray{0x77})

			// First Quadrant
			case x < width/2 && y < height/2 && x%2 == 0:
				img.Set(x, y, color.Black)
			case x < width/2 && y < height/2 && x%2 == 1:
				img.Set(x, y, color.White)

			// Second Quadrant
			case x > width/2 && y < height/2 && x%2 == 0:
				img.Set(x, y, color.Black)
			case x > width/2 && y < height/2 && x%2 == 1:
				img.Set(x, y, color.Black)

			// Third Quadrant
			case x < width/2 && y > height/2 && x%2 == 0:
				img.Set(x, y, color.White)
			case x < width/2 && y > height/2 && x%2 == 1:
				img.Set(x, y, color.White)

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

}
