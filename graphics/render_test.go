package graphics

import (
	"fmt"
	"image"
	"testing"
)

func init() {

}

func BenchmarkDrawTrianglesParallel(b *testing.B) { //414826807
	californiaGold := &Color{196, 130, 15, 255}
	lit := &DirectionLight{
		Direction: (&Vector3{0, 1, 1}).Normalize(),
		Color: &Color{
			R: 255,
			G: 255,
			B: 255,
			A: 255,
		},
	}
	triangles, _ := OpenObj("../render/teapot.obj", californiaGold)
	triangles = ApplyTransform(triangles, Translate(0, 0, 1.8))
	fmt.Println(len(triangles))
	im := image.NewRGBA(image.Rect(0, 0, 1000, 1000))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		DrawTrianglesParallel(im, triangles, []Light{lit})
	}
}

func BenchmarkDrawTrianglesParallelFaster(b *testing.B) { // 470799477
	californiaGold := &Color{196, 130, 15, 255}
	lit := &DirectionLight{
		Direction: (&Vector3{0, 1, 1}).Normalize(),
		Color: &Color{
			R: 255,
			G: 255,
			B: 255,
			A: 255,
		},
	}
	triangles, _ := OpenObj("../render/teapot.obj", californiaGold)
	triangles = ApplyTransform(triangles, Translate(0, 0, 1.8))
	fmt.Println(len(triangles))
	im := image.NewRGBA(image.Rect(0, 0, 1000, 1000))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		DrawTrianglesParallelFaster(im, triangles, []Light{lit})
	}
}

func BenchmarkVector3_FastNormalize(b *testing.B) {
	v := &Vector3{12.131, 14.132, -1.42}
	for n := 0; n < b.N; n++ {
		v.FastNormalize()
	}
}

func BenchmarkVector3_Normalize(b *testing.B) {
	v := &Vector3{12.131, 14.132, -1.42}
	for n := 0; n < b.N; n++ {
		v.Normalize()
	}
}

func BenchmarkInvSqrt(b *testing.B) {
	v := float64(1)
	for n := 0; n < b.N; n++ {
		fastInvSqrt(v)
		v++
	}
}

func BenchmarkSlowInvSqrt(b *testing.B) {
	v := float64(1)
	for n := 0; n < b.N; n++ {
		slowInvSqrt(v)
		v++
	}
}

func TestTriangle_Bary(t *testing.T) {
	p0 := &Vector3{0, 0, 0}
	p1 := &Vector3{-1, 1, 0}
	p2 := &Vector3{1, 1, 0}
	p3 := &Vector3{.2, .5, 0}
	p4 := &Vector3{.7, .9, 0}
	tri := NewTriangle(p0, p1, p2, nil)
	u, v, w := tri.Bary(p3)
	fmt.Println(u, v, w)
	u, v, w = tri.Bary(p4)
	fmt.Println(u, v, w, u+v+w)
}

func TestThings(T *testing.T) {
	var a, b, c uint8
	a = 255
	b = a
	c = a * b
	fmt.Println(c)

}
