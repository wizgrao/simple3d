package main

import (
	"flag"
	"github.com/wizgrao/simple3d/graphics"
	"image"
	"image/jpeg"
	"image/png"
	"math"
	"os"
)

var (
	outputFile = flag.String("o", "out.png", "Output File (png)")
	inputFile  = flag.String("i", "in.png", "Input file (png)")
	size       = flag.Int("s", 2000, "Size of output image")
	xt         = flag.Float64("xt", 0, "Translation in X direction")
	yt         = flag.Float64("yt", 0, "Translation in Y direction")
	zt         = flag.Float64("zt", 1.8, "Translation in Z direction")
	xr         = flag.Float64("xr", 2*math.Pi, "Rotation in X direction")
	yr         = flag.Float64("yr", 1, "Rotation in Y direction")
	zr         = flag.Float64("zr", 1, "Rotation in Z direction")
)

func main() {
	flag.Parse()
	im := image.NewRGBA(image.Rect(0, 0, *size, *size))
	fg := &graphics.Color{255, 255, 255, 255}
	bg := &graphics.Color{0, 0, 0, 255}
	for i := 0; i < *size; i++ {
		for j := 0; j < *size; j++ {
			im.Set(i, j, bg.ToRGBA())
		}
	}
	imfile, _ := os.Open("asdf.jpg")
	defer imfile.Close()
	textureIm, _ := jpeg.Decode(imfile)

	p1 := &graphics.Vector3{-1, -1, 0}
	p2 := &graphics.Vector3{1, -1, 0}
	p3 := &graphics.Vector3{1, 1, 0}
	p4 := &graphics.Vector3{-1, 1, 0}
	m := &graphics.TextureMaterial{
		Im:         textureIm,
		P1:         &graphics.Vector2{0, 0},
		P2:         &graphics.Vector2{1, 0},
		P3:         &graphics.Vector2{1, 1},
		SpecColor_: fg,
		SpecCoeff_: 8,
	}
	m2 := &graphics.TextureMaterial{
		Im:         textureIm,
		P1:         &graphics.Vector2{0, 0},
		P2:         &graphics.Vector2{0, 1},
		P3:         &graphics.Vector2{1, 1},
		SpecColor_: fg,
		SpecCoeff_: 8,
	}
	t := graphics.NewTriangle(p1, p2, p3, m)
	t2 := graphics.NewTriangle(p1, p4, p3, m2)
	triangles := []*graphics.Triangle{t, t2}
	transform := graphics.Translate(*xt, *yt, *zt).
		Mult(graphics.RotZ(*zr)).
		Mult(graphics.RotY(*yr)).
		Mult(graphics.RotX(*xr))
	triangles = graphics.ApplyTransform(triangles, transform)
	lit1 := &graphics.PointLight{
		Location: &graphics.Vector3{2, 1, 0},
		R:        1000,
	}
	lit2 := &graphics.PointLight{
		Location: &graphics.Vector3{-2, 1, 0},
		B:        1000,
	}
	lit3 := &graphics.PointLight{
		Location: &graphics.Vector3{0, 1, -1},
		G:        1000,
	}
	graphics.DrawTrianglesParallel(im, triangles, []graphics.Light{lit1, lit2, lit3})
	f, _ := os.Create(*outputFile)
	png.Encode(f, im)
}
