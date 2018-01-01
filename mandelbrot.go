package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"math/cmplx"
	"os"
)

var (
	maxIterations int
	size          int
)

func init() {
	flag.IntVar(&maxIterations, "i", 30, "max iterations")
	flag.IntVar(&size, "s", 2400, "Size of the image")
}

func main() {
	flag.Parse()
	drawMandelbrot(os.Stdout)
}

func smoothColor(n int, z complex128) *color.RGBA {

	var hue float64
	hue = (float64(n) + 1.0) - (math.Log(math.Log(cmplx.Abs(z))) / math.Log(2.0))
	hue = 0.95 + 20.0*hue // adjust to make it prettier
	// the hsv function expects values from 0 to 360
	for hue > 360.0 {
		hue -= 360.0
	}
	for hue < 0.0 {
		hue += 360.0
	}
	hsv := HSV{hue, float64(n) / (float64(n) + 8.0), 1.0}
	return hsv.RGBA()
}

// Mandelbrot equation
func mandelbrot(complexCoords complex128, maxIter int, radius float64) int {
	z := complexCoords
	for i := 0; i < maxIter-1; i++ {
		if cmplx.Abs(z) > radius {
			return i
		}
		// z = z^2 + C
		z = cmplx.Pow(z, 2) + complexCoords
	}
	return maxIter
}

func drawMandelbrot(out io.Writer) {
	imageHeight := size
	imageWidth := size

	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{size, size}})

	// values for calculation zoom
	imgZoomCenter := complex(0.0, 0.0)

	radius := 2.0

	zoomWidth := radius * 2
	pixelWidth := zoomWidth / float64(imageWidth)
	pixelHeight := pixelWidth
	viewHeight := (float64(imageHeight) / float64(imageWidth)) * zoomWidth
	left := (real(imgZoomCenter) - (zoomWidth / 2)) + pixelWidth/2
	top := (imag(imgZoomCenter) - (viewHeight / 2)) + pixelHeight/2

	// iterate and create image
	for xPos := 0; xPos < imageWidth; xPos++ {
		for yPos := 0; yPos < imageHeight; yPos++ {
			coord := complex(left+float64(xPos)*pixelWidth, top+float64(yPos)*pixelHeight)
			iteration := mandelbrot(coord, maxIterations, radius)
			if iteration < maxIterations {
				pointColor := smoothColor(iteration, coord)
				img.Set(xPos, yPos, pointColor)
			} else {
				img.Set(xPos, yPos, color.Black)
			}
		}

	}

	// Output image
	png.Encode(out, img)
}

type HSV struct {
	H, S, V float64
}

func (c *HSV) RGBA() *color.RGBA {
	alpha := 255 // opacity
	var r, g, b float64
	if c.S == 0 { //HSV from 0 to 1
		r = c.V * 255
		g = c.V * 255
		b = c.V * 255
	} else {
		h := c.H * 6
		if h == 6 {
			h = 0
		} //H must be < 1
		i := math.Floor(h) //Or ... var_i = floor( var_h )
		v1 := c.V * (1 - c.S)
		v2 := c.V * (1 - c.S*(h-i))
		v3 := c.V * (1 - c.S*(1-(h-i)))

		if i == 0 {
			r = c.V
			g = v3
			b = v1
		} else if i == 1 {
			r = v2
			g = c.V
			b = v1
		} else if i == 2 {
			r = v1
			g = c.V
			b = v3
		} else if i == 3 {
			r = v1
			g = v2
			b = c.V
		} else if i == 4 {
			r = v3
			g = v1
			b = c.V
		} else {
			r = c.V
			g = v1
			b = v2
		}

		r = r * 255 //RGB results from 0 to 255
		g = g * 255
		b = b * 255
	}
	rgb := &color.RGBA{uint8(r), uint8(g), uint8(b), uint8(alpha)}
	return rgb

}
