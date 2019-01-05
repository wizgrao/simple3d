package main

import (
	"flag"
	"fmt"
	"github.com/wizgrao/simple3d/graphics"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math"
	"os"
)

var (
	outputFile = flag.String("o", "out.png", "Output File (png)")
	inputSize  = flag.Int("r", 10, "subdivisions of sphere")
	inputFile  = flag.String("i", "asdf.jpg", "input file texture")
	size       = flag.Int("s", 2000, "Size of output image")
	xt         = flag.Float64("xt", 0, "Translation in X direction")
	yt         = flag.Float64("yt", 0, "Translation in Y direction")
	zt         = flag.Float64("zt", 1.8, "Translation in Z direction")
	xr         = flag.Float64("xr", 0, "Rotation in X direction")
	yr         = flag.Float64("yr", math.Pi, "Rotation in Y direction")
	zr         = flag.Float64("zr", math.Pi, "Rotation in Z direction")
	frames     = flag.Int("f", 10, "number of frames to render")
)

func main() {
	flag.Parse()

	imfile, err := os.Open(*inputFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	textureIm, err := jpeg.Decode(imfile)
	if err != nil {
		imfile.Close()
		imfile, err = os.Open(*inputFile)
		if err != nil {
			fmt.Println(err)
			return
		}
		textureIm, err = png.Decode(imfile)
		if err != nil {
			imfile.Close()
			imfile, err = os.Open(*inputFile)
			if err != nil {
				fmt.Println(err)
				return
			}
			textureIm, err = gif.Decode(imfile)
			if err != nil {
				fmt.Println("oop")
				return
			}
		}
	}
	imfile.Close()

	im := image.NewRGBA(image.Rect(0, 0, *size, *size))
	bg := &color.RGBA{0, 0, 0, 255}
	for i := 0; i < *size; i++ {
		for j := 0; j < *size; j++ {
			im.Set(i, j, bg)
		}
	}
	triangles := graphics.ImgSphere(*inputSize, textureIm)

	lit1 := &graphics.DirectionLight{
		Direction: (&graphics.Vector3{1, 1, 1}).Normalize(),
		Color:     graphics.White,
	}
	transform := graphics.Translate(*xt, *yt, *zt).
		Mult(graphics.RotZ(*zr)).
		Mult(graphics.RotY(*yr)).
		Mult(graphics.RotX(*xr))
	triangles = graphics.ApplyTransform(triangles, transform)
	rot := graphics.Translate(*xt, *yt, *zt).
		Mult(graphics.RotY(2 * math.Pi / float64(*frames))).
		Mult(graphics.Translate(-*xt, -*yt, -*zt))
	for i := 0; i < *frames; i++ {
		fmt.Println("starting image", i, "of ", *frames)
		graphics.DrawTrianglesParallel(im, triangles, []graphics.Light{lit1})

		f, _ := os.Create(fmt.Sprintf("%3d", i) + *outputFile)
		png.Encode(f, im)
		f.Close()
		triangles = graphics.ApplyTransform(triangles, rot)

	}
}
