package main

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	owm "github.com/briandowns/openweathermap"
	"github.com/njnygaard/kindling/gg"
)

const (
	ApiKey = "b40995f7e1911e427c0700778e542369"
)

func main() {
	// Basics
	width := 800
	height := 480

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
}
