package main

import (
	"github.com/wizgrao/blow/maps"
	"github.com/wizgrao/blow/wasmsocket"
	"github.com/wizgrao/simple3d/graphics"
	"math"
)

func main() {
	sock := wasmsocket.GetSocket("socket")
	fg := &graphics.Color{100, 100, 100, 255}
	gc := &graphics.Color{253, 181, 21, 255}
	bc := &graphics.Color{0, 58, 98, 255}
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

	lit1 := &graphics.PointLight{
		Location: &graphics.Vector3{1.5, -1, -0},
		R:        500,
		G:        500,
		B:        500,
	}
	lit2 := &graphics.PointLight{
		Location: &graphics.Vector3{-1.5, -1, -0},
		R:        500,
		B:        500,
		G:        500,
	}
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
	r := .5
	c1 := graphics.ApplyTransform(graphics.SphereMat(15, gm),graphics.Translate(-r, 0, 0).Mult(graphics.Scale(r)))
	c2 := graphics.ApplyTransform(graphics.SphereMat(15, bm),graphics.Translate(r, 0, 0).Mult(graphics.Scale(r)))
	mesh := append(c1, c2...)
	mesh = graphics.ApplyTransform(mesh, graphics.Translate(0, r, 1.5))
	triangles := append(mesh, t1, t2)
	asdf := graphics.RotX(math.Pi/8)
	triangles = graphics.ApplyTransform(triangles, asdf)
	mapper := &graphics.RayTraceMapper{
		Bounces: 3,
		Width: 128,
		Height: 128,
		XMin: -1,
		XMax: 1,
		YMin:-1,
		YMax: 1,
		Mesh: triangles,
		Lights:[]graphics.Light{lit1, lit2},
	}
	h := maps.NewHost(sock)
	h.Register(mapper)
	h.Start()
	select {}
}
