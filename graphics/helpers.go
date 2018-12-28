package graphics

import (
	"bufio"
	"image"
	"image/color"
	"os"
	"strconv"
	"strings"
	"sync"
	"math"
)


type Light struct {
	Norm *Vector3 //direction of light ray
	C    *color.RGBA
}


type Material struct {
	C *color.RGBA
	SpecColor *color.RGBA
	SpecCoeff float64
	AmbientCoeff float64
}

func (m *Material) Render(normal, camera *Vector3, l *Light) *color.RGBA{
	if normal.Dot(camera) >0 {
		normal = normal.Scale(-1)
	}
	inc := -normal.Dot(l.Norm)
	if inc < 0 {
		inc = 0
	}
	reflecc := normal.Scale(2*normal.Dot(l.Norm)).Sub(l.Norm)
	specCos := reflecc.Dot(camera)
	if specCos < 0 {
		specCos = 0
	}
	specColor := ColorMult(ColorScale(m.SpecColor, math.Pow(specCos, m.SpecCoeff)), l.C)
	diffColor := ColorMult(ColorScale(m.C, inc), l.C)
	ambColor := ColorScale(m.C, m.AmbientCoeff)
	return ColorAdd(ColorAdd(specColor, diffColor), ambColor)

}

func DrawTrianglesParallel(im *image.RGBA, t []*Triangle, l *Light) {
	width := im.Rect.Max.X - im.Rect.Min.X
	height := im.Rect.Max.Y - im.Rect.Min.Y
	zbuf := make([][]float64, width, width)
	zbuflock :=make([][]sync.Mutex, width ,width)
	wg := sync.WaitGroup{}
	for row := range zbuf {
		zbuf[row] = make([]float64, height, height)
		zbuflock[row] = make([]sync.Mutex, height, height)
	}
	for _,tri := range t {
		wg.Add(1)
		go func(tri *Triangle) {
			inc := -tri.Norm.Dot(l.Norm)
			if inc < 0 {
				inc *= -1
			}

			p0 := tri.P0.Dehom()
			p1 := tri.P1.Dehom()
			p2 := tri.P2.Dehom()

			minx := int(lin(min3(p0.X, p1.X, p2.X), -1, 1, 0, float64(width)))
			miny := int(lin(min3(p0.Y, p1.Y, p2.Y), -1, 1, 0, float64(height)))
			maxx := int(lin(max3(p0.X, p1.X, p2.X), -1, 1, 0, float64(width)))
			maxy := int(lin(max3(p0.Y, p1.Y, p2.Y), -1, 1, 0, float64(height)))

			for i := maxi(minx, 0); i < mini(maxx+1, width); i++ {
				for j := maxi(miny, 0); j < mini(maxy+1, height); j++ {
					coordx := lin(float64(i), 0, float64(width), -1, 1)
					coordy := lin(float64(j), 0, float64(height), -1, 1)

					screenCoord := &Vector2{coordx, coordy}
					if In(p0, p1, p2, screenCoord) {
						dePerp := tri.DePerp(&Vector2{coordx, coordy})
						zbuflock[i][j].Lock()
						if zbuf[i][j] == 0 || zbuf[i][j] > dePerp.Z {
							zbuf[i][j] = dePerp.Z
							im.Set(i, j, tri.Render(tri.Norm, screenCoord.Hom().Normalize(), l))
						}
						zbuflock[i][j].Unlock()
					}
				}
			}
			wg.Done()
		}(tri)
	}
	wg.Wait()
}

func ColorInterp(c0, c1, c2 *color.RGBA, s0, s1, s2 float64) *color.RGBA{
	return &color.RGBA{
		R: uint8(float64(c0.R) * s0 + float64(c1.R) * s1 + float64(c2.R) * s2),
		G: uint8(float64(c0.G) * s0 + float64(c1.G) * s1 + float64(c2.G) * s2),
		B: uint8(float64(c0.B) * s0 + float64(c1.B) * s1 + float64(c2.B) * s2),
		A: uint8(float64(c0.A) * s0 + float64(c1.A) * s1 + float64(c2.A) * s2),
	}
}

func ColorAdd(c0, c1 *color.RGBA) *color.RGBA{
	return &color.RGBA{
		R: uint8(mini(255, int(c0.R) + int(c1.R))),
		G: uint8(mini(255, int(c0.G) + int(c1.G))),
		B: uint8(mini(255, int(c0.B) + int(c1.B))),
		A: uint8(mini(255, int(c0.A) + int(c1.A))),
	}
}

func ColorMult(c0, c1 *color.RGBA) *color.RGBA{
	return &color.RGBA{
		R:  uint8(int(c0.R) * int(c1.R)/255),
		G: uint8(int(c0.G) * int(c1.G)/255),
		B:  uint8(int(c0.B) * int(c1.B)/255),
		A:  uint8(int(c0.A) * int(c1.A)/255),
	}
}


func ColorScale(c *color.RGBA, s float64) *color.RGBA{
	return &color.RGBA{
		R: uint8(float64(c.R)*s),
		G:uint8(float64(c.G)*s),
		B: uint8(float64(c.B)*s),
		A: c.A,
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

			points = append(points, &Vector3{x, y, z})
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

func min(x, y float64) float64 {
	if x < y {
		return x
	}
	return y
}


func mini(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y float64) float64 {
	if x > y {
		return x
	}
	return y
}

func maxi(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func maxu(x, y uint8) uint8 {
	if x > y {
		return x
	}
	return y
}

func minu(x, y uint8) uint8 {
	if x < y {
		return x
	}
	return y
}

func min3(x, y, z float64) float64 {
	return min(min(x, y), z)
}

func max3(x, y, z float64) float64 {
	return max(max(x, y), z)
}