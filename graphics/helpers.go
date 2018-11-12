package graphics

import (
	"bufio"
	"image"
	"image/color"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Vector3 struct {
	X float64
	Y float64
	Z float64
}

type Vector2 struct {
	X float64
	Y float64
}

type Triangle struct {
	P0   *Vector3
	P1   *Vector3
	P2   *Vector3
	Norm *Vector3
	C    *color.RGBA
}

type Light struct {
	Norm *Vector3
	C    *color.RGBA
}

func (v *Vector3) Dehom() *Vector2 {
	return &Vector2{
		X: v.X / v.Z,
		Y: v.Y / v.Z,
	}
}

func DrawTriangle(im *image.RGBA, t *Triangle, l *Light) {
	inc := -t.Norm.Dot(l.Norm)
	if inc < 0 {
		inc *= -1
	}
	drawColor := color.RGBA{
		R: uint8(float64(t.C.R) * float64(l.C.R) * inc / 255.0),
		G: uint8(float64(t.C.G) * float64(l.C.G) * inc / 255.0),
		B: uint8(float64(t.C.B) * float64(l.C.B) * inc / 255.0),
		A: 255,
	}
	minPoint := im.Bounds().Min
	maxPoint := im.Bounds().Max
	imDx := maxPoint.X - minPoint.X
	imDy := maxPoint.Y - minPoint.Y

	p1 := t.P0.Dehom()
	p2 := t.P1.Dehom()
	p3 := t.P2.Dehom()

	maxX := math.Max(math.Max(p1.X, p2.X), p3.X)
	maxY := math.Max(math.Max(p1.Y, p2.Y), p3.Y)
	minX := math.Min(math.Min(p1.X, p2.X), p3.X)
	minY := math.Min(math.Min(p1.Y, p2.Y), p3.Y)

	pminx := int(lin(minX, -1, 1, float64(minPoint.X), float64(maxPoint.X)))
	pmaxx := int(lin(maxX, -1, 1, float64(minPoint.X), float64(maxPoint.X)))
	pminy := int(lin(minY, -1, 1, float64(minPoint.Y), float64(maxPoint.Y)))
	pmaxy := int(lin(maxY, -1, 1, float64(minPoint.Y), float64(maxPoint.Y)))

	pminx = max(pminx, minPoint.X)
	pmaxx = min(pmaxx, maxPoint.X)
	pminy = max(pminy, minPoint.Y)
	pmaxy = min(pmaxy, maxPoint.Y)
	for i := pminx - 1; i < pmaxx+1; i++ {
		for j := pminy - 1; j < pmaxy+1; j++ {
			unproj := &Vector2{2*float64(i-minPoint.X)/float64(imDx) - 1, 2*float64(j-minPoint.Y)/float64(imDy) - 1}

			if In(p1, p2, p3, unproj) {
				im.Set(i, maxPoint.Y-j+minPoint.Y, drawColor)
			}
		}
	}

}

func DrawTriangles(im *image.RGBA, t []*Triangle, l *Light) {
	sort.Slice(t, func(i, j int) bool {
		return t[i].Centroid().Z > t[j].Centroid().Z
	})

	for _, tri := range t {
		DrawTriangle(im, tri, l)
	}
}
func (t *Triangle) Centroid() *Vector3 {
	return t.P0.Add(t.P1).Add(t.P2).Scale(1.0 / 3.0)
}
func (u *Vector3) Dot(v *Vector3) float64 {
	return u.X*v.X + u.Y*v.Y + u.Z*v.Z
}

func (u *Vector2) Sub(v *Vector2) *Vector2 {
	return &Vector2{
		X: u.X - v.X,
		Y: u.Y - v.Y,
	}
}

func (u *Vector3) Sub(v *Vector3) *Vector3 {
	return &Vector3{
		X: u.X - v.X,
		Y: u.Y - v.Y,
		Z: u.Z - v.Z,
	}
}

func (u *Vector3) Add(v *Vector3) *Vector3 {
	return &Vector3{
		X: u.X + v.X,
		Y: u.Y + v.Y,
		Z: u.Z + v.Z,
	}
}

func In(p0, p1, p2, p *Vector2) bool {
	area := 0.5 * (-p1.Y*p2.X + p0.Y*(-p1.X+p2.X) + p0.X*(p1.Y-p2.Y) + p1.X*p2.Y)
	s := 1 / (2 * area) * (p0.Y*p2.X - p0.X*p2.Y + (p2.Y-p0.Y)*p.X + (p0.X-p2.X)*p.Y)
	t := 1 / (2 * area) * (p0.X*p1.Y - p0.Y*p1.X + (p0.Y-p1.Y)*p.X + (p1.X-p0.X)*p.Y)
	return s > 0 && t > 0 && 1-s-t > 0

}

func Cross(v1, v2 *Vector3) *Vector3 {
	return &Vector3{
		X: v1.Y*v2.Z - v1.Z*v2.Y,
		Y: -(v1.X*v2.Z - v1.Z*v2.Z),
		Z: v1.X*v2.Y - v1.Y*v2.X,
	}
}

func (v *Vector3) Norm() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v *Vector3) Normalize() *Vector3 {
	norm := v.Norm()
	return &Vector3{
		X: v.X / norm,
		Y: v.Y / norm,
		Z: v.Z / norm,
	}
}

func (v *Vector3) Scale(s float64) *Vector3 {
	return &Vector3{
		X: v.X * s,
		Y: v.Y * s,
		Z: v.Z * s,
	}
}

func CalcNorm(p0, p1, p2 *Vector3) *Vector3 {
	res := Cross(p0.Sub(p1), p0.Sub(p2)).Normalize()
	return res
}

func NewTriangle(p0, p1, p2 *Vector3, c *color.RGBA) *Triangle {
	return &Triangle{
		P0:   p0,
		P1:   p1,
		P2:   p2,
		Norm: CalcNorm(p0, p1, p2),
		C:    c,
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

			points = append(points, &Vector3{x, y, -z + 1.8})
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

func lin(p, mini, maxi, mino, maxo float64) float64 {
	return (maxo-mino)*(p-mini)/(maxi-mini) + mino
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