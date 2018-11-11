package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"math"
	"sort"
	"bufio"
	"strconv"
	"fmt"
	"strings"
)

func main() {


	image := image.NewRGBA(image.Rect(0, 0, 2000, 2000))
	bg := &color.RGBA{196, 130, 15, 255}
	fg := &color.RGBA{0, 50, 98, 255}
	for i := 0; i < 2000; i++ {
		for j := 0; j < 2000; j++ {
			image.Set(i, j, fg)
		}
	}
	lit := &Light{
		norm: (&Vector3{0, 1, 1}).Normalize(),
		c: &color.RGBA{
			R:255,
			G:255,
			B:255,
			A:255,
		},
	}


	triangles, _ := OpenObj(os.Args[1], bg)
	fmt.Println(len(triangles))
	DrawTriangles(image, triangles, lit)

	f, _ := os.Create("out.png")
	png.Encode(f, image)
}

type Vector3 struct {
	x float64
	y float64
	z float64
}

type Vector2 struct {
	x float64
	y float64
}

type Triangle struct {
	p0   *Vector3
	p1   *Vector3
	p2   *Vector3
	norm *Vector3
	c    *color.RGBA
}

type Light struct {
	norm *Vector3
	c    *color.RGBA
}

func (v *Vector3) Dehom() *Vector2 {
	return &Vector2{
		x: v.x/v.z,
		y: v.y/v.z,
	}
}

func DrawTriangle(im *image.RGBA, t *Triangle, l *Light) {
	inc := -t.norm.Dot(l.norm)
	if inc < 0 {
		inc *= -1
	}
	drawColor := color.RGBA{
		R: uint8(float64(t.c.R) * float64(l.c.R) * inc / 255.0),
		G: uint8(float64(t.c.G) * float64(l.c.G) * inc / 255.0),
		B: uint8(float64(t.c.B) * float64(l.c.B) * inc / 255.0),
		A: 255,
	}
	minPoint := im.Bounds().Min
	maxPoint := im.Bounds().Max
	imDx := maxPoint.X - minPoint.X
	imDy := maxPoint.Y - minPoint.Y

	p1 := t.p0.Dehom()
	p2 := t.p1.Dehom()
	p3 := t.p2.Dehom()

	maxX := math.Max(math.Max(p1.x, p2.x), p3.x)
	maxY := math.Max(math.Max(p1.y, p2.y), p3.y)
	minX := math.Min(math.Min(p1.x, p2.x), p3.x)
	minY := math.Min(math.Min(p1.y, p2.y), p3.y)

	pminx := int(lin(minX, -1, 1, float64(minPoint.X),float64(maxPoint.X)))
	pmaxx := int(lin(maxX, -1, 1, float64(minPoint.X),float64(maxPoint.X)))
	pminy := int(lin(minY, -1, 1, float64(minPoint.Y),float64(maxPoint.Y)))
	pmaxy := int(lin(maxY, -1, 1, float64(minPoint.Y),float64(maxPoint.Y)))


	pminx = max(pminx, minPoint.X)
	pmaxx = min(pmaxx, maxPoint.X)
	pminy = max(pminy, minPoint.Y)
	pmaxy = min(pmaxy, maxPoint.Y)
	for i:=pminx - 1; i < pmaxx + 1; i++ {
		for j:=pminy - 1; j < pmaxy + 1; j++ {
			unproj := &Vector2{ 2*float64(i-minPoint.X)/float64(imDx) - 1, 2*float64(j-minPoint.Y)/float64(imDy) - 1}

			if In(p1, p2, p3, unproj) {
				im.Set(i, maxPoint.Y - j + minPoint.Y, drawColor)
			}
		}
	}

}

func DrawTriangles(im *image.RGBA, t []*Triangle, l *Light) {
	sort.Slice(t, func(i, j int) bool {
		return t[i].Centroid().z > t[j].Centroid().z
	})

	for _, tri := range t {
		DrawTriangle(im, tri, l)
	}
}
func (t *Triangle) Centroid() *Vector3 {
	return t.p0.Add(t.p1).Add(t.p2).Scale(1.0/3.0)
}
func (u *Vector3) Dot(v *Vector3) float64 {
	return u.x*v.x + u.y*v.y + u.z*v.z
}


func (u *Vector2) Sub(v *Vector2) *Vector2 {
	return &Vector2{
		x: u.x - v.x,
		y: u.y - v.y,
	}
}

func (u *Vector3) Sub(v *Vector3) *Vector3{
	return &Vector3{
		x: u.x - v.x,
		y: u.y - v.y,
		z: u.z - v.z,
	}
}

func (u *Vector3) Add(v *Vector3) *Vector3{
	return &Vector3{
		x: u.x + v.x,
		y: u.y + v.y,
		z: u.z + v.z,
	}
}

func In(p0, p1, p2, p *Vector2) bool {
    area := 0.5 *(-p1.y*p2.x + p0.y*(-p1.x + p2.x) + p0.x*(p1.y - p2.y) + p1.x*p2.y)
	s := 1/(2*area)*(p0.y*p2.x - p0.x*p2.y + (p2.y - p0.y)*p.x + (p0.x - p2.x)*p.y)
	t := 1/(2*area)*(p0.x*p1.y - p0.y*p1.x + (p0.y - p1.y)*p.x + (p1.x - p0.x)*p.y)
	return s > 0 && t > 0 && 1-s-t > 0

}

func Cross(v1, v2 *Vector3) *Vector3 {
	return &Vector3{
		x: v1.y*v2.z - v1.z*v2.y,
		y: -(v1.x*v2.z - v1.z*v2.z),
		z: v1.x*v2.y - v1.y*v2.x,
	}
}

func (v*Vector3) Norm() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z)
}

func (v* Vector3) Normalize() *Vector3 {
	norm := v.Norm()
	return &Vector3{
		x: v.x/norm,
		y: v.y/norm,
		z: v.z/norm,
	}
}

func (v*Vector3) Scale(s float64) *Vector3 {
	return &Vector3{
		x: v.x * s,
		y: v.y * s,
		z: v.z * s,
	}
}

func CalcNorm(p0, p1, p2 *Vector3) *Vector3 {
	res := Cross(p0.Sub(p1), p0.Sub(p2)).Normalize()
	return res
}

func NewTriangle(p0, p1, p2 *Vector3, c  *color.RGBA) *Triangle {
	return &Triangle{
		p0:   p0,
		p1:   p1,
		p2:   p2,
		norm: CalcNorm(p0, p1, p2),
		c:    c,
	}
}

func OpenObj(filename string, rgba *color.RGBA) ([]*Triangle, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewScanner(f)
	reader.Split(bufio.ScanWords)
	var points []*Vector3
	var triangles []*Triangle

	for reader.Scan() {
		t := reader.Text()
		if t == "v" {
			reader.Scan()
			xs := reader.Text()
			reader.Scan()
			ys := reader.Text()
			reader.Scan()
			zs := reader.Text()

			x, _ := strconv.ParseFloat(xs, 64)
			y, _ := strconv.ParseFloat(ys, 64)
			z, _ := strconv.ParseFloat(zs, 64)

			points = append(points, &Vector3{x, y, -z+1.8})
		}
		if t == "f" {
			reader.Scan()
			xs := strings.Split(reader.Text(), "/")[0]
			reader.Scan()
			ys := strings.Split(reader.Text(), "/")[0]
			reader.Scan()
			zs := strings.Split(reader.Text(), "/")[0]

			x, _ := strconv.ParseInt(xs, 10, 32)
			y, _ := strconv.ParseInt(ys, 10, 32)
			z, _ := strconv.ParseInt(zs, 10, 32)

			triangles = append(triangles, NewTriangle(points[x-1], points[y-1], points[z-1], rgba))

		}
	}
	return triangles, nil

}

func lin(p, mini, maxi, mino, maxo float64) float64{
	return (maxo - mino)*(p-mini)/(maxi-mini) + mino
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}