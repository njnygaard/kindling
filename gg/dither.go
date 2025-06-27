package gg

import "image/color"

// Bayer dithering matrix (8x8)
var bayerMatrix = [8][8]int{
	{0, 48, 12, 60, 3, 51, 15, 63},
	{32, 16, 44, 28, 35, 19, 47, 31},
	{8, 56, 4, 52, 11, 59, 7, 55},
	{40, 24, 36, 20, 43, 27, 39, 23},
	{2, 50, 14, 62, 1, 49, 13, 61},
	{34, 18, 46, 30, 33, 17, 45, 29},
	{10, 58, 6, 54, 9, 57, 5, 53},
	{42, 26, 38, 22, 41, 25, 37, 21},
}

// DitherPixel converts a gray level to black or white using Bayer dithering
func DitherPixel(grayLevel int, x, y int) color.Color {
	// Get threshold from Bayer matrix
	threshold := bayerMatrix[x%8][y%8]

	// Convert gray level (0-255) to threshold comparison
	// If gray level is higher than threshold, make it white, otherwise black
	if grayLevel > threshold {
		return color.White
	}
	return color.Black
}
