package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	owm "github.com/briandowns/openweathermap"
	"github.com/njnygaard/kindling/gg"
)

const (
	ApiKey = "b40995f7e1911e427c0700778e542369"
)

func generateWeather(width, height int) {
	var fondamentoSize float64 = 60
	var winterSongSize float64 = 90

	var cityLabelVerticalOffset float64 = -2
	var cityLabelHorizontalOffset = 0.5

	// Draw Text
	dc := gg.NewContext(width, height)
	dc.SetRGB(1)
	dc.Clear()
	dc.SetRGB(0)

	// Time Tagging (to determine if it is successfully updating)
	if err := dc.LoadFontFace("SourceCodePro-Regular.ttf", 14); err != nil {
		panic(err)
	}
	currentTime := time.Now()
	dc.DrawStringAnchored(currentTime.Format("Monday, January 2, 2006"), 0, float64(0), 0, 1)
	dc.DrawStringAnchored(currentTime.Format("15:04:05"), float64(width), float64(0), 1.05, 1)

	/****************/
	/*** Brisbane ***/
	/****************/

	// Weather
	w, err := owm.NewCurrent("F", "en", ApiKey)
	if err != nil {
		log.Fatalln(err)
	}
	err = w.CurrentByName("Brisbane")
	if err != nil {
		return
	}
	// City Label
	if err := dc.LoadFontFace("Fondamento-Regular.ttf", fondamentoSize); err != nil {
		panic(err)
	}
	dc.DrawStringAnchored(fmt.Sprintf("Brisbane %.0f°F", w.Main.Temp), float64(width)/2, float64(1*height+120)/4, cityLabelHorizontalOffset, cityLabelVerticalOffset)

	// Description
	if err := dc.LoadFontFace("WinterSong-owRGB.ttf", winterSongSize); err != nil {
		panic(err)
	}
	c := cases.Title(language.English)
	s := c.String(w.Weather[0].Description)
	dc.DrawStringAnchored(s, float64(width)/2, float64(1*height)/4, 0.5, 1)

	/*********************/
	/*** Saint-Émilion ***/
	/*********************/

	// Weather
	w, err = owm.NewCurrent("F", "en", ApiKey)
	if err != nil {
		log.Fatalln(err)
	}
	err = w.CurrentByName("Saint-Émilion")
	if err != nil {
		return
	}
	// City Label
	if err := dc.LoadFontFace("Fondamento-Regular.ttf", fondamentoSize); err != nil {
		panic(err)
	}
	dc.DrawStringAnchored(fmt.Sprintf("Saint-Émilion %.0f°F", w.Main.Temp), float64(width)/2, float64(3*height+120)/4, cityLabelHorizontalOffset, cityLabelVerticalOffset)

	// Description
	if err := dc.LoadFontFace("WinterSong-owRGB.ttf", winterSongSize); err != nil {
		panic(err)
	}
	c = cases.Title(language.English)
	s = c.String(w.Weather[0].Description)
	dc.DrawStringAnchored(s, float64(width)/2, float64(3*height)/4, 0.5, 1)

	// Export for Text
	err = dc.SavePNG("trmnl/weather.png")
	if err != nil {
		return
	}

	generateTestPattern(width, height)
}

func generateTestPattern(width, height int) {
	// Create a new image
	img := image.NewGray(image.Rect(0, 0, width, height))

	// Draw the classic TV test pattern
	gg.DrawTestPattern(img, width, height)

	// Save to file
	file, err := os.Create("trmnl/test_pattern.png")
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}

	//println("TV test pattern saved as tv_test_pattern.png")

}

func generateTestPatternBW(width, height int) {
	// Create a new image
	img := image.NewGray(image.Rect(0, 0, width, height))

	// Draw the black & white dithered test pattern
	gg.DrawTestPatternBW(img, width, height)

	// Save to file
	file, err := os.Create("trmnl/test_pattern_bw.png")
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}

	//println("Black & white test pattern saved as test_pattern_bw.png")
}

func generateCircularGradient(x, y, width, height int, invert bool) (grayLevel int) {
	// Calculate distance from center
	centerX := width / 2
	centerY := height / 2
	dx := float64(x - centerX)
	dy := float64(y - centerY)

	// Normalize distance to create gradient
	maxDist := math.Sqrt(float64(centerX*centerX + centerY*centerY))
	dist := math.Sqrt(dx*dx + dy*dy)

	// Calculate gradient value (white center by default)
	gradient := (1.0 - dist/maxDist) * 64
	if invert {
		gradient = 64 - gradient // Invert to make black center
	}

	grayLevel = int(gradient)
	if grayLevel < 0 {
		grayLevel = 0
	}
	return
}

func generateSquareGradient(x, y, width, height int, invert bool) (grayLevel int) {
	// Calculate normalized distances from edges
	normalizedX := math.Abs(float64(x)-float64(width)/2) / (float64(width) / 2)
	normalizedY := math.Abs(float64(y)-float64(height)/2) / (float64(height) / 2)

	// Take the maximum of the normalized distances to create square gradient
	maxNormalizedDist := math.Max(normalizedX, normalizedY)

	// Create sharp falloff to black at edges by squaring the distance
	// and using a smaller scale factor
	gradient := (1.0 - maxNormalizedDist*maxNormalizedDist) * 64
	if invert {
		gradient = 64 - gradient
	}
	grayLevel = int(gradient)
	if grayLevel < 0 {
		grayLevel = 0
	}
	return
}

func generateDitherDemo(width, height int, shape string, invert bool, filename string) {
	// Create a simple demonstration of dithering
	img := image.NewGray(image.Rect(0, 0, width, height))

	// Create a gradient from black to white using dithering
	for y := 0; y < height; y++ {
		// Calculate gray level based on Y position (0 to 64 (8x8))
		// grayLevel := int(float64(y) / float64(height) * 64)

		for x := 0; x < width; x++ {
			// Use the dithering function to convert gray level to black/white
			var pixelColor color.Color
			switch shape {
			case "square":
				pixelColor = gg.DitherPixel(generateSquareGradient(x, y, width, height, invert), x, y)
			case "circle":
				pixelColor = gg.DitherPixel(generateCircularGradient(x, y, width, height, invert), x, y)
			default:
				panic("Invalid shape")
			}
			img.Set(x, y, pixelColor)
		}
	}

	// Save to file
	file, err := os.Create(fmt.Sprintf("trmnl/%s.png", filename))
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}
}

func main() {
	// Basics
	width := 800
	height := 480

	generateWeather(width, height)
	// I don't really like these test patterns.
	// The dithering demo is better.
	// generateTestPattern(width, height)
	// generateTestPatternBW(width, height)
	generateDitherDemo(width, height, "square", false, "dither_square")
	generateDitherDemo(width, height, "circle", false, "dither_circle")
	generateDitherDemo(width, height, "square", true, "dither_square_invert")
	generateDitherDemo(width, height, "circle", true, "dither_circle_invert")
}
