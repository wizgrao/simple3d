package main

import (
	"image"
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
	shadow = flag.Bool("h", false, "whether to draw shadows")

)

func main() {
	flag.Parse()
	im := image.NewRGBA(image.Rect(0, 0, *size, *size))
	fg := &graphics.Color{255, 255, 255, 255}
	bg := &graphics.Color{0, 0, 0, 255}

	m := &graphics.SolidMaterial{
		Color:         fg,
		SpecColor_:    &graphics.Color{255, 255, 255, 255},
		SpecCoeff_:    8,
		AmbientCoeff_: .01,
	}

	for i := 0; i < *size; i++ {
		for j := 0; j < *size; j++ {
			im.Set(i, j, bg.ToRGBA())
		}
	}
	triangles, _ := graphics.OpenObj(*inputFile, fg)

	lit1 := &graphics.PointLight{
		Location: &graphics.Vector3{1.5, -1, -0},
		R: 500,
	}
	lit2 := &graphics.PointLight{
		Location: &graphics.Vector3{-1.5, -1, -0},
		B: 500,
	}
	lit3 := &graphics.PointLight{
		Location: &graphics.Vector3{0, -1, -1},
		G: 500,
	}
	/*lit4 := &graphics.DirectionLight{
		Direction: &graphics.Vector3{0,1, 0},
		Color: fg,
	}*/
	transform := graphics.Translate(*xt, *yt, *zt).
		Mult(graphics.RotZ(*zr)).
		Mult(graphics.RotY(*yr)).
		Mult(graphics.RotX(*xr))
	triangles = graphics.ApplyTransform(triangles, transform)
	f1 := &graphics.Vector3{-10, 1, .01}
	f2 := &graphics.Vector3{10, 1, .01}
	f3 := &graphics.Vector3{10, 1, 10}
	f4 := &graphics.Vector3{-10, 1, 10}
	n := &graphics.Vector3{0, -1, 0}
	t1 := graphics.NewTriangle(f1, f2, f3, m)
	t1.N0 = n
	t1.N1 = n
	t1.N2 = n

	t2 := graphics.NewTriangle(f1, f3, f4, m)
	t2.N0 = n
	t2.N1 = n
	t2.N2 = n
	triangles = append(triangles, t1, t2)
	fn := graphics.DrawTrianglesParallel
	if *shadow {
		fn = graphics.DrawTrianglesParallelShadow
	}
	fn(im, triangles, []graphics.Light{lit1, lit2, lit3})
	f, _ := os.Create(*outputFile)
	png.Encode(f, im)
}
