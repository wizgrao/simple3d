package graphics

import (
	"testing"
	"image"
	"image/color"
	"fmt"
	"os"
	"image/png"
)

func init() {

}


func BenchmarkDrawTrianglesParallel(b *testing.B) {
	californiaGold := &color.RGBA{196, 130, 15, 255}
	lit := &Light{
		Norm: (&Vector3{0, 1, 1}).Normalize(),
		C: &color.RGBA{
			R:255,
			G:255,
			B:255,
			A:255,
		},
	}
	triangles, _ := OpenObj("../render/teapot.obj", californiaGold)
	triangles = ApplyTransform(triangles, Translate(0,0,1.8))
	fmt.Println(len(triangles))
	im := image.NewRGBA(image.Rect(0, 0, 1000, 1000))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		DrawTrianglesParallel(im, triangles, lit)
	}
}

func TestTriangle_Bary(t *testing.T) {
	p0 := &Vector3{0,0,0}
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

func TestDrawTriangle(t *testing.T) {
	im := image.NewRGBA(image.Rect(0, 0, 500, 500))
	berkeleyBlue := &color.RGBA{0, 50, 98, 255}
	tri := NewTriangle( &Vector3{0,0,1},&Vector3{250,200,1}, &Vector3{500, 0, 1}, berkeleyBlue)
	c0 := &color.RGBA{255, 0, 0, 255}
	c1 := &color.RGBA{0, 255, 0, 255}
	c2 := &color.RGBA{0, 0, 255, 255}

	for i := 0; i < 500; i++ {
		for j := 0; j < 500; j++ {
			dp := tri.DePerp(&Vector2{float64(i), float64(j)})
			if tri.In(dp) {
				u, v, w := tri.Bary(dp)
				im.Set(i, j, ColorInterp(c0, c1, c2, u, v, w))
			}
		}
	}
	f, _ := os.Create("test.png")
	png.Encode(f, im)

}

func TestThings(T *testing.T) {
	var a, b, c uint8
	a = 255
	b = a
	c = a*b
	fmt.Println(c)


}