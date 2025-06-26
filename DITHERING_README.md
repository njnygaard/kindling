# Black & White Dithering Test Pattern

This implementation provides a new version of the test pattern function that simulates gray colors using only black and white pixels through dithering.

## Overview

The `DrawTestPatternBW` function creates a test pattern similar to the original `DrawTestPattern`, but uses Bayer dithering to convert gray levels into black and white pixel patterns. This is useful for displays that only support binary (black/white) output.

## Key Features

### 1. Bayer Dithering
- Uses an 8x8 Bayer matrix for consistent, predictable dithering patterns
- Converts gray levels (0-255) to black/white pixels based on position and threshold
- Maintains visual hierarchy and contrast relationships

### 2. Test Pattern Components
- **Color Bars**: 9 bars ranging from white to black, dithered to black/white
- **PLUGE Bars**: 4 bars for brightness/contrast adjustment (super black to medium gray)
- **Checkerboard**: Already black/white, no changes needed
- **Circular Pattern**: White circles and crosshairs on dithered background

### 3. Functions

#### Main Function
```go
func DrawTestPatternBW(img *image.Gray, width, height int)
```
Creates the complete black & white dithered test pattern.

#### Dithering Function
```go
func DitherPixel(grayLevel int, x, y int) color.Color
```
Converts a gray level to black or white using Bayer dithering. Can be used independently for custom dithering applications.

#### Helper Functions
- `drawColorBarsBW()` - Dithered color bars
- `drawPlugeBarsBW()` - Dithered PLUGE bars  
- `drawCircularPatternBW()` - Dithered circular patterns
- `fillRectBW()` - Fill rectangle with dithered gray level
- `drawCircleBW()` - Draw circle with dithered gray level
- `fillCircleBW()` - Fill circle with dithered gray level

## Usage Examples

### Basic Test Pattern
```go
img := image.NewGray(image.Rect(0, 0, 800, 480))
gg.DrawTestPatternBW(img, 800, 480)
```

### Custom Dithering
```go
// Create a gradient using dithering
for y := 0; y < height; y++ {
    grayLevel := int(float64(y) / float64(height) * 255)
    for x := 0; x < width; x++ {
        pixelColor := gg.DitherPixel(grayLevel, x, y)
        img.Set(x, y, pixelColor)
    }
}
```

## Generated Files

When running the main program, three test pattern files are generated:

1. `trmnl/test_pattern.png` - Original gray test pattern
2. `trmnl/test_pattern_bw.png` - Black & white dithered test pattern
3. `trmnl/dither_demo.png` - Simple gradient demonstration of dithering

## Technical Details

### Bayer Matrix
The 8x8 Bayer matrix provides 64 different threshold levels (0-63), which when scaled to 0-255 provides good dithering quality for most applications.

### Dithering Algorithm
```go
threshold := bayerMatrix[x%8][y%8]
if grayLevel > threshold {
    return color.White
}
return color.Black
```

### Performance
- Dithering adds computational overhead but is still efficient for typical image sizes
- The Bayer matrix lookup is O(1) per pixel
- Memory usage remains the same as the original implementation

## Applications

This implementation is particularly useful for:
- E-paper displays (Kindle, etc.)
- Monochrome LCD screens
- Print applications requiring binary output
- Legacy display systems
- Testing display contrast and brightness capabilities 