package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"github.com/wizgrao/simple3d/graphics"
	"flag"
	"math"
)

var (
	outputFile = flag.String("o", "out.png", "Output File (png)")
	inputFile = flag.String("i", "in.obj", "Input file (png)")
	size = flag.Int("s", 2000, "Size of output image")
	xt = flag.Float64("xt", 0, "Translation in X direction")
	yt = flag.Float64("yt", 0, "Translation in Y direction")
	zt = flag.Float64("zt", 1.8, "Translation in Z direction")
	xr = flag.Float64("xr", 0, "Rotation in X direction")
	yr = flag.Float64("yr", math.Pi, "Rotation in Y direction")
	zr = flag.Float64("zr", math.Pi, "Rotation in Z direction")

)

func main() {
	flag.Parse()
	im := image.NewRGBA(image.Rect(0, 0, *size, *size))
	fg := &color.RGBA{255, 255, 255, 255}
	bg := &color.RGBA{0, 0, 0, 255}
	for i := 0; i < *size; i++ {
		for j := 0; j < *size; j++ {
			im.Set(i, j, bg)
		}
	}
	triangles, _ := graphics.OpenObj(*inputFile, fg)

	lit1 := &graphics.PointLight{
		Location: &graphics.Vector3{2, 1, 0},
		R: 500,
	}
	lit2 := &graphics.PointLight{
		Location: &graphics.Vector3{-2, 1, 0},
		B: 500,
	}
	lit3 := &graphics.PointLight{
		Location: &graphics.Vector3{0, 1, -1},
		G: 500,
	}
	transform := graphics.Translate(*xt, *yt, *zt).
		Mult(graphics.RotZ(*zr)).
		Mult(graphics.RotY(*yr)).
		Mult(graphics.RotX(*xr))
	triangles = graphics.ApplyTransform(triangles, transform)
	graphics.DrawTrianglesParallel(im, triangles, []graphics.Light{lit1, lit2, lit3})
	f, _ := os.Create(*outputFile)
	png.Encode(f, im)
}
