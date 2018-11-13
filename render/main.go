package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"github.com/wizgrao/raytracer/graphics"
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
	californiaGold := &color.RGBA{196, 130, 15, 255}
	berkeleyBlue := &color.RGBA{0, 50, 98, 255}
	for i := 0; i < *size; i++ {
		for j := 0; j < *size; j++ {
			im.Set(i, j, berkeleyBlue)
		}
	}
	/*lit := &graphics.Light{
		Norm: (&graphics.Vector3{0, 1, 1}).Normalize(),
		C: &color.RGBA{
			R:255,
			G:255,
			B:255,
			A:255,
		},
	}*/
	triangles, _ := graphics.OpenObj(*inputFile, californiaGold)
	transform := graphics.Translate(*xt, *yt, *zt).
		Mult(graphics.RotZ(*zr)).
		Mult(graphics.RotY(*yr)).
		Mult(graphics.RotX(*xr))
	triangles = graphics.ApplyTransform(triangles, transform)
	//graphics.DrawTriangles(im, triangles, lit)
	points2 := []*graphics.Vector2{
		{0, 30},
		{1000, 400},
		{500, 1200},
	}
	graphics.DrawTriangle2(im, points2, californiaGold)
	f, _ := os.Create(*outputFile)
	png.Encode(f, im)
}
