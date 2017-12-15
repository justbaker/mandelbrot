package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"io"
	"math/cmplx"
	"os"
)

var maxIterations = flag.Int("i", 30, "max iterations")
var size = flag.Int("s", 2400, "Size of the image")
var zoomWidthRatio = flag.Int("z", 1, "Size of the image")
var xStart = flag.Float64("x", 0.0, "x start position")
var yStart = flag.Float64("y", 0.0, "y start position")

func main() {
	flag.Parse()
	drawMandelbrot(os.Stdout, 2.0)
}

func colorForPoint(iteration int, maxIt int) color.RGBA {
	if iteration <= maxIt {
		return color.RGBA{uint8(iteration), uint8(iteration) * uint8(255/maxIt), uint8(0), 255}
	} else {
		return color.RGBA{uint8(0), uint8(0), uint8(0), 255}
	}
}

// Mandelbrot equation
func mandelbrot(complexCoords complex128) int {
	z := complexCoords
	for i := 0; i < *maxIterations-1; i++ {
		if cmplx.Abs(z) > 2 {
			return i
		}
		// z = z^2 + C
		z = cmplx.Pow(z, 2) + complexCoords
	}
	return *maxIterations - 1
}

func drawMandelbrot(out io.Writer, radius float64) {
	// set dimensions
	imageHeight := *size
	imageWidth := *size

	// init image
	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{*size, *size}})

	// values for calculation zoom
	imgZoomCenter := complex(*xStart, *yStart)

	zoomWidth := radius * 2 * float64(*zoomWidthRatio)
	pixelWidth := zoomWidth / float64(imageWidth)
	pixelHeight := pixelWidth
	viewHeight := (float64(imageHeight) / float64(imageWidth)) * zoomWidth
	left := (real(imgZoomCenter) - (zoomWidth / 2)) + pixelWidth/2
	top := (imag(imgZoomCenter) - (viewHeight / 2)) + pixelHeight/2

	// iterate and create image
	for xPos := 0; xPos < imageWidth; xPos++ {
		for yPos := 0; yPos < imageHeight; yPos++ {
			coord := complex(left+float64(xPos)*pixelWidth, top+float64(yPos)*pixelHeight)
			iteration := mandelbrot(coord)

			img.Set(xPos, yPos, colorForPoint(iteration, *maxIterations))
		}

	}

	// Output image
	png.Encode(out, img)
}
