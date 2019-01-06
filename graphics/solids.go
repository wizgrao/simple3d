package graphics

import (
	"math"
	"image"
)

var (
	White = &Color{255, 255, 255, 255}
)

func Sphere(subdivisions int) []*Triangle {
	ret := make([]*Triangle, 0)
	pts := make([][]*Vector3, subdivisions-1)
	m := &SolidMaterial{
		Color:      White,
		SpecColor_: White,
		SpecCoeff_: 8,
		AmbientCoeff_: .05,
	}
	for i := range pts {
		pts[i] = make([]*Vector3, subdivisions*2)
		theta := math.Pi * (-0.5*float64(subdivisions) + 1.0 + float64(i)) / float64(subdivisions)
		y := math.Sin(theta)
		radius := math.Cos(theta)
		for j := range pts[i] {
			phi := math.Pi * float64(j) / float64(subdivisions)
			x := math.Cos(phi) * radius
			z := math.Sin(phi) * radius
			pts[i][j] = &Vector3{x, y, z}
		}
		if i == 0 {

			continue
		}
		for j := range pts[i] {
			l := (2*subdivisions + j - 1) % (2 * subdivisions)
			r := j
			t := i
			b := i - 1

			t1 := NewTriangle(pts[t][l], pts[t][r], pts[b][l], m)
			t1.N0 = pts[t][l]
			t1.N1 = pts[t][r]
			t1.N2 = pts[b][l]
			t2 := NewTriangle(pts[b][r], pts[t][r], pts[b][l], m)
			t2.N0 = pts[b][r]
			t2.N1 = pts[t][r]
			t2.N2 = pts[b][l]

			ret = append(ret, t1, t2)
		}
	}
	i := 0
	for j := range pts[i] {
		l := (2*subdivisions + j - 1) % (2 * subdivisions)
		r := j
		top := subdivisions - 2
		t1 := NewTriangle(pts[i][l], pts[i][r], &Vector3{0, -1, 0}, m)
		t1.N0 = pts[i][l]
		t1.N1 = pts[i][r]
		t1.N2 = &Vector3{0, -1, 0}
		t2 := NewTriangle(pts[top][l], pts[top][r], &Vector3{0, 1, 0}, m)
		t2.N0 = pts[top][l]
		t2.N1 = pts[top][r]
		t2.N2 = &Vector3{0, 1, 0}
		ret = append(ret, t1, t2)

	}
	return ret
}

func ImgSphere(subdivisions int, im image.Image) []*Triangle {
	ret := make([]*Triangle, 0)
	pts := make([][]*Vector3, subdivisions-1)
	for i := range pts {
		pts[i] = make([]*Vector3, subdivisions*2)
		theta := math.Pi * (-0.5*float64(subdivisions) + 1.0 + float64(i)) / float64(subdivisions)
		y := math.Sin(theta)
		radius := math.Cos(theta)
		for j := range pts[i] {
			phi := math.Pi * float64(j) / float64(subdivisions)
			x := math.Cos(phi) * radius
			z := math.Sin(phi) * radius
			pts[i][j] = &Vector3{x, y, z}
		}
		if i == 0 {

			continue
		}
		for j := range pts[i] {
			l := (2*subdivisions + j - 1) % (2 * subdivisions)
			r := j
			t := i
			b := i - 1
			tlTex := &Vector2{float64(l)/float64(2*subdivisions), float64(t+1)/float64(subdivisions)}
			blTex := &Vector2{float64(l)/float64(2*subdivisions), float64(b+1)/float64(subdivisions)}
			trTex := &Vector2{float64(r)/float64(2*subdivisions), float64(t+1)/float64(subdivisions)}
			brTex := &Vector2{float64(r)/float64(2*subdivisions), float64(b+1)/float64(subdivisions)}
			if j == 0 {
				trTex.X = 1
				brTex.X = 1
			}
			m1 := &TextureMaterial{
				P1: tlTex,
				P2: trTex,
				P3: blTex,
				Im: im,
				SpecColor_: ColorScale(White, .5),
				SpecCoeff_: 8,
				AmbientCoeff_: .05,
			}
			t1 := NewTriangle(pts[t][l], pts[t][r], pts[b][l], m1)
			t1.N0 = pts[t][l]
			t1.N1 = pts[t][r]
			t1.N2 = pts[b][l]

			m2 := &TextureMaterial{
				P1: brTex,
				P2: trTex,
				P3: blTex,
				Im: im,
				SpecColor_: ColorScale(White, .5),
				SpecCoeff_: 8,
				AmbientCoeff_: .05,
			}
			t2 := NewTriangle(pts[b][r], pts[t][r], pts[b][l], m2)
			t2.N0 = pts[b][r]
			t2.N1 = pts[t][r]
			t2.N2 = pts[b][l]

			ret = append(ret, t1, t2)
		}
	}
	i := 0
	for j := range pts[i] {
		l := (2*subdivisions + j - 1) % (2 * subdivisions)
		r := j

		m1 := &TextureMaterial{
			P1: &Vector2{float64(l)/float64(2*subdivisions), 1/float64(subdivisions)},
			P2: &Vector2{float64(r)/float64(2*subdivisions), 1/float64(subdivisions)},
			P3: &Vector2{.5, 0},
			Im: im,
			SpecColor_: ColorScale(White, .5),
			SpecCoeff_: 8,
			AmbientCoeff_: .05,
		}

		top := subdivisions - 2
		t1 := NewTriangle(pts[i][l], pts[i][r], &Vector3{0, -1, 0}, m1)
		t1.N0 = pts[i][l]
		t1.N1 = pts[i][r]
		t1.N2 = &Vector3{0, -1, 0}
		m2 := &TextureMaterial{
			P1: &Vector2{float64(l)/float64(2*subdivisions), 1-1/float64(subdivisions)},
			P2: &Vector2{float64(r)/float64(2*subdivisions), 1-1/float64(subdivisions)},
			P3: &Vector2{.5, 1},
			Im: im,
			SpecColor_: ColorScale(White, .5),
			SpecCoeff_: 8,
			AmbientCoeff_: .05,
		}
		t2 := NewTriangle(pts[top][l], pts[top][r], &Vector3{0, 1, 0}, m2)
		t2.N0 = pts[top][l]
		t2.N1 = pts[top][r]
		t2.N2 = &Vector3{0, 1, 0}
		ret = append(ret, t1, t2)

	}
	return ret
}