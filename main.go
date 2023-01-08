package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"strings"
	"time"

	owm "github.com/briandowns/openweathermap"
	"github.com/njnygaard/kindling/gg"
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
	// Basics
	width := 758
	height := 1024

	var fondamentoSize float64 = 60
	var winterSongSize float64 = 100

	var cityLabelVerticalOffset float64 = -2
	var cityLabelHorizontalOffset float64 = 0.5
	// upLeft := image.Point{0, 0}
	// lowRight := image.Point{width, height}
	// img := image.NewGray(image.Rectangle{upLeft, lowRight})

	// Draw Shapes
	// DO NOT USE image.Gray16
	// The kindle will subtly break
	// drawTestPattern(img, width, height)
	// drawGrid(img, width, height)

	// Draw Text
	dc := gg.NewContext(width, height)
	dc.SetRGB(1)
	dc.Clear()
	dc.SetRGB(0)

	// Time Tagging (to determine if it is successfully updating)
	if err := dc.LoadFontFace("SourceCodePro-Regular.ttf", 10); err != nil {
		panic(err)
	}
	current_time := time.Now()
	dc.DrawStringAnchored(current_time.Format("Monday, January 2, 2006"), 0, float64(0), 0, 1)
	dc.DrawStringAnchored(current_time.Format("15:04:05"), float64(width), float64(0), 1.05, 1)

	/****************/
	/*** Brisbane ***/
	/****************/
	// Weather
	w, err := owm.NewCurrent("F", "en", API_KEY)
	if err != nil {
		log.Fatalln(err)
	}
	w.CurrentByName("Brisbane")
	// City Label
	if err := dc.LoadFontFace("Fondamento-Regular.ttf", fondamentoSize); err != nil {
		panic(err)
	}
	dc.DrawStringAnchored(fmt.Sprintf("Brisbane %.0f°F", w.Main.Temp), float64(width)/2, float64(1*height)/6, cityLabelHorizontalOffset, cityLabelVerticalOffset)
	// ranged := w.Main.TempMax - w.Main.TempMin
	// scaledTemperature := w.Main.Temp - w.Main.TempMin
	// ratio := scaledTemperature / ranged
	// Description
	if err := dc.LoadFontFace("WinterSong-owRGB.ttf", winterSongSize); err != nil {
		panic(err)
	}
	s := strings.Title(w.Weather[0].Description)
	dc.DrawStringAnchored(s, float64(width)/2, float64(1*height)/6, 0.5, 1)
	// Temperature
	// if err := dc.LoadFontFace("SourceCodePro-Regular.ttf", 30); err != nil {
	// 	panic(err)
	// }

	// Temp Scale
	// dc.DrawStringAnchored(fmt.Sprintf("%.0f°F", w.Main.Temp), float64(width)*ratio, float64(1*height)/6, 1.1, 7)
	// if w.Main.TempMin+1 < w.Main.Temp {
	// 	dc.DrawStringAnchored(fmt.Sprintf("%.0f°F", w.Main.TempMin), 0, float64(1*height)/6, -0.1, 7)
	// }
	// if w.Main.TempMax-1 > w.Main.Temp {
	// 	dc.DrawStringAnchored(fmt.Sprintf("%.0f°F", w.Main.TempMax), float64(width), float64(1*height)/6, 1.1, 7)
	// }

	/*********************/
	/*** San Francisco ***/
	/*********************/
	// Weather
	w, err = owm.NewCurrent("F", "en", API_KEY)
	if err != nil {
		log.Fatalln(err)
	}
	w.CurrentByName("San Francisco")
	// City Label
	if err := dc.LoadFontFace("Fondamento-Regular.ttf", fondamentoSize); err != nil {
		panic(err)
	}
	dc.DrawStringAnchored(fmt.Sprintf("San Francisco %.0f°F", w.Main.Temp), float64(width)/2, float64(3*height)/6, cityLabelHorizontalOffset, cityLabelVerticalOffset)
	// ranged = w.Main.TempMax - w.Main.TempMin
	// scaledTemperature = w.Main.Temp - w.Main.TempMin
	// ratio = scaledTemperature / ranged
	// Description
	if err := dc.LoadFontFace("WinterSong-owRGB.ttf", winterSongSize); err != nil {
		panic(err)
	}
	s = strings.Title(w.Weather[0].Description)
	dc.DrawStringAnchored(s, float64(width)/2, float64(3*height)/6, 0.5, 1)
	// Temperature
	// if err := dc.LoadFontFace("SourceCodePro-Regular.ttf", 30); err != nil {
	// 	panic(err)
	// }

	// Temp Scale
	// dc.DrawStringAnchored(fmt.Sprintf("%.0f°F", w.Main.Temp), float64(width)*ratio, float64(3*height)/6, 1.1, 7)
	// if w.Main.TempMin+1 < w.Main.Temp {
	// 	dc.DrawStringAnchored(fmt.Sprintf("%.0f°F", w.Main.TempMin), 0, float64(3*height)/6, -0.1, 7)
	// }
	// if w.Main.TempMax-1 > w.Main.Temp {
	// 	dc.DrawStringAnchored(fmt.Sprintf("%.0f°F", w.Main.TempMax), float64(width), float64(3*height)/6, 1.1, 7)
	// }

	/****************/
	/*** Saint-Émilion ***/
	/****************/
	// Weather
	w, err = owm.NewCurrent("F", "en", API_KEY)
	if err != nil {
		log.Fatalln(err)
	}
	w.CurrentByName("Saint-Émilion")
	// City Label
	if err := dc.LoadFontFace("Fondamento-Regular.ttf", fondamentoSize); err != nil {
		panic(err)
	}
	dc.DrawStringAnchored(fmt.Sprintf("Saint-Émilion %.0f°F", w.Main.Temp), float64(width)/2, float64(5*height)/6, cityLabelHorizontalOffset, cityLabelVerticalOffset)
	// ranged = w.Main.TempMax - w.Main.TempMin
	// scaledTemperature = w.Main.Temp - w.Main.TempMin
	// ratio = scaledTemperature / ranged
	// Description
	if err := dc.LoadFontFace("WinterSong-owRGB.ttf", winterSongSize); err != nil {
		panic(err)
	}
	s = strings.Title(w.Weather[0].Description)
	dc.DrawStringAnchored(s, float64(width)/2, float64(5*height)/6, 0.5, 1)
	// Temperature
	// if err := dc.LoadFontFace("SourceCodePro-Regular.ttf", 30); err != nil {
	// 	panic(err)
	// }

	// Temp Scale
	// fmt.Println("w.Main.Temp", w.Main.Temp)
	// fmt.Println("w.Main.TempMin", w.Main.TempMin)
	// fmt.Println("w.Main.TempMax", w.Main.TempMax)
	// dc.DrawStringAnchored(fmt.Sprintf("%.0f°F", w.Main.Temp), float64(width)*ratio, float64(5*height)/6, 1.1, 7)
	// if w.Main.TempMin+1 < w.Main.Temp {
	// 	dc.DrawStringAnchored(fmt.Sprintf("%.0f°F", w.Main.TempMin), 0, float64(5*height)/6, -0.1, 7)
	// }
	// if w.Main.TempMax-1 > w.Main.Temp {
	// 	dc.DrawStringAnchored(fmt.Sprintf("%.0f°F", w.Main.TempMax), float64(width), float64(5*height)/6, 1.1, 7)
	// }

	// Export for Text
	dc.SavePNG("out.png")

	// Encode as PNG for Shapes
	// f, _ := os.Create("image.png")
	// png.Encode(f, img)
}

func drawThird(position int) {

}

// func formatTemp(min float64, temp float64, max float64) (f string) {
// 	minInt := int(min)
// 	maxInt := int(max)
// 	tempInt := int(temp)

// 	// tempInt := 51

// 	for i := minInt; i < maxInt; i++ {
// 		switch {
// 		case i == minInt:
// 			f = f + fmt.Sprintf("%.0f°F", min)
// 		case i == maxInt-1:
// 			f = f + fmt.Sprintf("%.0f°F", max)
// 		case i == tempInt:
// 			if tempInt == minInt || tempInt == maxInt {
// 				// continue
// 				f = f + "                "
// 			} else {
// 				log.Printf("tempInt: %d", tempInt)
// 				log.Printf("minInt: %d", minInt)
// 				log.Printf("maxInt: %d", maxInt)
// 				log.Printf("equal condition")
// 				f = f + fmt.Sprintf("%.0f°F", temp)
// 			}
// 		default:
// 			f = f + "                "
// 		}
// 	}

// 	return f
// }

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

	// for x := 0; x < width; x++ {
	// 	for y := 0; y < height; y++ {
	// 		switch {

	// 		case x == (width/2)+2 || x == (width/2)-2:
	// 			img.Set(x, y, color.Gray{0xff})

	// 		case y > (height/2)+2 && x < (height/2)-2:
	// 			img.Set(x, y, color.Gray{0x00})
	// 		}
	// 	}
	// }

	// Set color for each pixel.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			switch {

			// Center Circle
			case isCircle((1*width)/2, (1*height)/2, 50, x, y):
				img.Set(x, y, color.Gray{0xff})
			case isCircle((1*width)/2, (1*height)/2, 75, x, y):
				img.Set(x, y, color.Gray{0x77})
			case isCircle((1*width)/2, (1*height)/2, 100, x, y):
				img.Set(x, y, color.Gray{0x00})

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

			// Upper Right Circles
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

			// Lower Left Circles
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

			// Lower Right Circles
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
