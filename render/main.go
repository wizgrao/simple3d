package main

import (
	"flag"
	"github.com/wizgrao/simple3d/graphics"
	"image"
	"image/png"
	"math"
	"os"
	"github.com/wizgrao/blow/maps"
)

var (
	outputFile = flag.String("o", "out.png", "Output File (png)")
	inputFile  = flag.String("i", "in.obj", "Input file (png)")
	size       = flag.Int("s", 2000, "Size of output image")
	xt         = flag.Float64("xt", 0, "Translation in X direction")
	yt         = flag.Float64("yt", 0, "Translation in Y direction")
	zt         = flag.Float64("zt", 1.8, "Translation in Z direction")
	xr         = flag.Float64("xr", 0, "Rotation in X direction")
	yr         = flag.Float64("yr", math.Pi, "Rotation in Y direction")
	zr         = flag.Float64("zr", math.Pi, "Rotation in Z direction")
	shadow     = flag.Bool("h", false, "whether to draw shadows")
	trace     = flag.Bool("t", false, "whether to raytrace")
	bounces = flag.Int("b", 3, "number of bounces on the ray tracer")
	circles = flag.Bool("circles", false, "draw alternate scene")
	grey = flag.Bool("g", false, "use white lighting instead of colored")
	parallel = flag.Int("p", 16, "number of parallel goroutines to use for rendering")
)

func main() {
	flag.Parse()
	im := image.NewRGBA(image.Rect(0, 0, *size, *size))
	fg := &graphics.Color{100, 100, 100, 255}
	gc := &graphics.Color{253, 181, 21, 255}
	bc := &graphics.Color{0,58,98,255}
	bg := &graphics.Color{0, 0, 0, 255}

	m := &graphics.SolidMaterial{
		Color:         fg,
		SpecColor_:    &graphics.Color{150, 150, 150, 255},
		SpecCoeff_:    8,
		AmbientCoeff_: .01,
	}
	gm := &graphics.SolidMaterial{
		Color:         graphics.ColorScale(gc, .3),
		SpecColor_:    graphics.ColorScale(gc, .7),
		SpecCoeff_:    8,
		AmbientCoeff_: .01,
	}
	bm := &graphics.SolidMaterial{
		Color:         graphics.ColorScale(bc, .6),
		SpecColor_:    graphics.ColorScale(bc, 1.4),
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
		R:        500,
	}
	lit2 := &graphics.PointLight{
		Location: &graphics.Vector3{-1.5, -1, -0},
		B:        500,
	}
	lit3 := &graphics.PointLight{
		Location: &graphics.Vector3{0, -1, -1},
		G:        500,
	}

	if *grey {
		lit1.G =  500
		lit1.B = 500
		lit2.G =  500
		lit2.R = 500
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

	if *circles {
		r := .5
		c1 := graphics.ApplyTransform(graphics.SphereMat(50, gm),graphics.Translate(-r, 0, 0).Mult(graphics.Scale(r)))
		c2 := graphics.ApplyTransform(graphics.SphereMat(50, bm),graphics.Translate(r, 0, 0).Mult(graphics.Scale(r)))
		mesh := append(c1, c2...)
		mesh = graphics.ApplyTransform(mesh, graphics.Translate(0, r, 1.5))
		triangles = append(mesh, t1, t2)
		asdf := graphics.RotX(math.Pi/8)
		_ = asdf
		triangles = graphics.ApplyTransform(triangles, asdf)

	}else {
		triangles = append(triangles, t1, t2)
	}


	if *shadow {
		_ = lit3
		graphics.DrawTrianglesParallelShadow(im, triangles, []graphics.Light{lit1, lit2 /*, lit3*/})
	}
	if *trace {
		source := &graphics.PixelSource{13, im}
		mapper := &graphics.RayTraceMapper{
			Bounces: *bounces,
			Width: *size,
			Height: *size,
			XMin: -1,
			XMax: 1,
			YMin:-1,
			YMax: 1,
			Mesh: triangles,
			Lights:[]graphics.Light{lit1, lit2},
		}
		writer := &graphics.WriterMapper{im, 0, *size * *size}
		maps.GeneratorSource(source, nil).MapLocalParallel(mapper, *parallel).MapLocal(writer).Sink()
	}else {
		graphics.DrawTrianglesParallel(im, triangles, []graphics.Light{lit1, lit2 /*, lit3*/})
	}
	f, _ := os.Create(*outputFile)
	png.Encode(f, im)
}
