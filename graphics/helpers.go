package graphics

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

type Light interface {
	Norm(*Vector3) *Vector3
	Intensity(*Vector3) *Color
	Transform(*Mat4) Light
}

type Color struct {
	R float64
	G float64
	B float64
	A float64
}

func (c *Color) ToRGBA() *color.RGBA {
	return &color.RGBA{
		R: uint8(max(min(c.R, 255), 0)),
		G: uint8(max(min(c.G, 255), 0)),
		B: uint8(max(min(c.B, 255), 0)),
		A: uint8(max(min(c.A, 255), 0)),
	}
}

type DirectionLight struct {
	Direction *Vector3
	Color     *Color
}

type PointLight struct {
	Location *Vector3
	R        float64
	G        float64
	B        float64
}

func (d *PointLight) Norm(v *Vector3) *Vector3 {
	return v.Sub(d.Location).Normalize()
}

func (d *PointLight) Intensity(v *Vector3) *Color {
	dist := v.Sub(d.Location).Norm()
	return &Color{
		R: d.R / (dist * dist),
		G: d.G / (dist * dist),
		B: d.B / (dist * dist),
		A: 255,
	}
}

func (d *PointLight) Transform(m *Mat4) Light {
	return &PointLight{
		Location: m.Dot(d.Location.Hom()).Dehom(),
		R:        d.R,
		G:        d.G,
		B:        d.B,
	}
}

func (d *DirectionLight) Norm(_ *Vector3) *Vector3 {
	return d.Direction
}

func (d *DirectionLight) Intensity(_ *Vector3) *Color {
	return d.Color
}

func (d *DirectionLight) Transform(m *Mat4) Light {
	return &DirectionLight{
		Direction: m.Dot(d.Direction.Ext()).Dehom(),
		Color:     d.Color,
	}
}

type Material interface {
	C(*Vector2) *Color
	SpecColor(*Vector2) *Color
	SpecCoeff(*Vector2) float64
	AmbientCoeff(vector2 *Vector2) float64
}

type SolidMaterial struct {
	Color         *Color
	SpecColor_    *Color
	SpecCoeff_    float64
	AmbientCoeff_ float64
}

func (s *SolidMaterial) C(_ *Vector2) *Color {
	return s.Color
}
func (s *SolidMaterial) SpecColor(_ *Vector2) *Color {
	return s.SpecColor_
}
func (s *SolidMaterial) SpecCoeff(_ *Vector2) float64 {
	return s.SpecCoeff_
}
func (s *SolidMaterial) AmbientCoeff(vector2 *Vector2) float64 {
	return s.AmbientCoeff_
}

type TextureMaterial struct {
	Im image.Image
	P1 *Vector2
	P2 *Vector2
	P3 *Vector2

	SpecColor_    *Color
	SpecCoeff_    float64
	AmbientCoeff_ float64
}

func (s *TextureMaterial) C(vec *Vector2) *Color {
	u := vec.X
	v := vec.Y
	w := 1 - u - v
	texNormalCoordinate := s.P1.Scale(u).Add(s.P2.Scale(v)).Add(s.P3.Scale(w))
	texCoordinateX := lin(texNormalCoordinate.X, 0, 1, float64(s.Im.Bounds().Min.X), float64(s.Im.Bounds().Max.X))
	texCoordinateY := lin(texNormalCoordinate.Y, 0, 1, float64(s.Im.Bounds().Min.Y), float64(s.Im.Bounds().Max.Y))
	bl := ToColor(s.Im.At(int(math.Floor(texCoordinateX)), int(math.Floor(texCoordinateY))))
	br := ToColor(s.Im.At(int(math.Ceil(texCoordinateX)), int(math.Floor(texCoordinateY))))
	tl := ToColor(s.Im.At(int(math.Floor(texCoordinateX)), int(math.Ceil(texCoordinateY))))
	tr := ToColor(s.Im.At(int(math.Ceil(texCoordinateX)), int(math.Ceil(texCoordinateY))))

	fracX := texCoordinateX - math.Floor(texCoordinateX)
	fracY := texCoordinateY - math.Floor(texCoordinateY)
	top := ColorAdd(ColorScale(tl, 1-fracX), ColorScale(tr, fracX))
	bottom := ColorAdd(ColorScale(bl, 1-fracX), ColorScale(br, fracX))

	return ColorAdd(ColorScale(bottom, 1-fracY), ColorScale(top, fracY))

}

func ToColor(c color.Color) *Color {
	r, g, b, _ := c.RGBA()
	return &Color{
		R: float64(r / 256),
		G: float64(g / 256),
		B: float64(b / 256),
		A: 255,
	}
}
func (s *TextureMaterial) SpecColor(_ *Vector2) *Color {
	return s.SpecColor_
}
func (s *TextureMaterial) SpecCoeff(_ *Vector2) float64 {
	return s.SpecCoeff_
}
func (s *TextureMaterial) AmbientCoeff(vector2 *Vector2) float64 {
	return s.AmbientCoeff_
}

func Render(m Material, normal, camera *Vector3, lights []Light, v *Vector3, uv *Vector2) *Color {
	ret := ColorScale(m.C(uv), m.AmbientCoeff(uv))
	for _, l := range lights {
		lnorm := l.Norm(v)
		lintense := l.Intensity(v)

		inc := -normal.Dot(l.Norm(v))
		if inc < 0 {
			inc = 0
		}
		reflecc := normal.Scale(2 * normal.Dot(l.Norm(v))).Sub(lnorm)
		specCos := reflecc.Dot(camera)
		if specCos < 0 {
			specCos = 0
		}
		specColor := ColorMult(ColorScale(m.SpecColor(uv), math.Pow(specCos, m.SpecCoeff(uv))), lintense)
		diffColor := ColorMult(ColorScale(m.C(uv), inc), lintense)
		ret = ColorAdd(ret, ColorAdd(specColor, diffColor))
	}
	return ret
}

var zero = &Vector3{}

func GetSpecularShadow(env []*Triangle, m Material, normal *Vector3, lights []Light, uv *Vector2) *Color {
	ret := ColorScale(m.C(uv), m.AmbientCoeff(uv))
	for _, l := range lights {
		shouldrender := true
		lnorm := l.Norm(zero)
		for _, tri := range env {
			intersect := tri.DePerp(lnorm.Dehom()) //interesct is the vector from the surface to the surface in the way of the light
			if !tri.In(intersect) || intersect.Norm() < .001 {
				continue
			}
			lintersect := l.Norm(intersect) //lintersect is the vector from the light to the norm
			if lintersect.Dot(intersect) < 0 {
				shouldrender = false
				break
			}
		}
		if !shouldrender {
			continue
		}
		lintense := l.Intensity(zero)
		inc := -normal.Dot(l.Norm(zero))
		if inc < 0 {
			inc = 0
		}
		diffColor := ColorMult(ColorScale(m.C(uv), inc), lintense)


		ret = ColorAdd(ret, diffColor)
	}
	return ret
}

func RayCast(env []*Triangle, lights []Light, vec *Vector3, bounce int) *Color {
	var mindist float64
	var mintriangle *Triangle
	var minInteresction *Vector3
	for _, t := range env {
		intersection := t.RayIntersect(vec)
		if !t.In(intersection) {
			continue
		}
		dist := intersection.Dot(vec)
		if dist < 0.1 {
			continue
		}
		if mindist == 0 || dist < mindist {
			mindist = dist
			mintriangle = t
			minInteresction = intersection
		}
	}
	if mintriangle == nil {
		return &Color{
			A: 255,
		}
	}

	u, v, w := mintriangle.Bary(minInteresction)
	norm := mintriangle.N0.Scale(u).Add(mintriangle.N1.Scale(v)).Add(mintriangle.N2.Scale(w)).Normalize()
	reflect := minInteresction.Sub(norm.Scale(2 * norm.Dot(minInteresction))).Normalize()
	newLights := make([]Light, len(lights))
	transform := Translate(-minInteresction.X, -minInteresction.Y, -minInteresction.Z)
	for i, light := range lights {
		newLights[i] = light.Transform(transform)
	}
	uv := &Vector2{u, v}
	newEnv := ApplyTransform(env, transform)
	c := GetSpecularShadow(newEnv, mintriangle.Material, norm, newLights, uv)
	if bounce > 0 {
		c = ColorAdd(c, ColorMult(RayCast(newEnv, newLights, reflect, bounce-1), mintriangle.Material.SpecColor(uv)))
	}
	return c

}

func RenderShadow(t *Triangle, env []*Triangle, m Material, normal, camera, v *Vector3, lights []Light, uv *Vector2) *Color {
	newenv := ApplyTransform(env, Translate(-v.X, -v.Y, -v.Z))
	ret := ColorScale(m.C(uv), m.AmbientCoeff(uv))

	for _, l := range lights {
		shouldrender := true
		lnorm := l.Norm(v)
		for i, tri := range newenv {
			if env[i] == t {
				continue
			}
			intersect := tri.DePerp(lnorm.Dehom()) //interesct is the vector from the surface to the surface in the way of the light
			if !tri.In(intersect) {
				continue
			}
			lintersect := l.Norm(intersect) //lintersect is the vector from the light to the norm
			if lintersect.Dot(intersect) < 0 {
				shouldrender = false
				break
			}
		}
		if !shouldrender {
			continue
		}
		lintense := l.Intensity(v)
		inc := -normal.Dot(l.Norm(v))
		if inc < 0 {
			inc = 0
		}
		reflecc := normal.Scale(2 * normal.Dot(l.Norm(v))).Sub(lnorm)
		specCos := reflecc.Dot(camera)
		if specCos < 0 {
			specCos = 0
		}
		specColor := ColorMult(ColorScale(m.SpecColor(uv), math.Pow(specCos, m.SpecCoeff(uv))), lintense)
		diffColor := ColorMult(ColorScale(m.C(uv), inc), lintense)
		ret = ColorAdd(ret, ColorAdd(specColor, diffColor))
	}
	return ret
}

func DrawTrianglesParallel(im *image.RGBA, t []*Triangle, l []Light) {
	width := im.Rect.Max.X - im.Rect.Min.X
	height := im.Rect.Max.Y - im.Rect.Min.Y
	zbuf := make([][]float64, width, width)
	zbuflock := make([][]sync.Mutex, width, width)
	wg := sync.WaitGroup{}
	for row := range zbuf {
		zbuf[row] = make([]float64, height, height)
		zbuflock[row] = make([]sync.Mutex, height, height)
	}
	for _, tri := range t {
		wg.Add(1)
		go func(tri *Triangle) {
			p0 := tri.P0.Dehom()
			p1 := tri.P1.Dehom()
			p2 := tri.P2.Dehom()

			minx := int(lin(min3(p0.X, p1.X, p2.X), -1, 1, 0, float64(width)))
			miny := int(lin(min3(p0.Y, p1.Y, p2.Y), -1, 1, 0, float64(height)))
			maxx := int(lin(max3(p0.X, p1.X, p2.X), -1, 1, 0, float64(width)))
			maxy := int(lin(max3(p0.Y, p1.Y, p2.Y), -1, 1, 0, float64(height)))

			for i := maxi(minx-1, 0); i < mini(maxx+1, width); i++ {
				for j := maxi(miny-1, 0); j < mini(maxy+1, height); j++ {
					coordx := lin(float64(i), 0, float64(width), -1, 1)
					coordy := lin(float64(j), 0, float64(height), -1, 1)

					screenCoord := &Vector2{coordx, coordy}
					if !In(p0, p1, p2, screenCoord) {
						continue
					}
					dePerp := tri.DePerp(&Vector2{coordx, coordy})
					zbuflock[i][j].Lock()
					if dePerp.Z <= 0 || (zbuf[i][j] != 0 && zbuf[i][j] <= dePerp.Z) {
						zbuflock[i][j].Unlock()
						continue
					}

					zbuf[i][j] = dePerp.Z
					u, v, w := tri.Bary(dePerp)
					norm := tri.N0.Scale(u).Add(tri.N1.Scale(v)).Add(tri.N2.Scale(w)).Normalize()
					im.Set(i, j, Render(tri.Material, norm, screenCoord.Hom().Normalize(), l, dePerp, &Vector2{u, v}).ToRGBA())

					zbuflock[i][j].Unlock()

				}
			}
			wg.Done()
		}(tri)
	}
	wg.Wait()
}

func DrawTrianglesParallelFaster(im *image.RGBA, t []*Triangle, l []Light) {
	width := im.Rect.Max.X - im.Rect.Min.X
	height := im.Rect.Max.Y - im.Rect.Min.Y
	zbuf := make([][]float64, width, width)
	zbuflock := make([][]sync.Mutex, width, width)
	wg := sync.WaitGroup{}
	for row := range zbuf {
		zbuf[row] = make([]float64, height, height)
		zbuflock[row] = make([]sync.Mutex, height, height)
	}
	for _, tri := range t {
		wg.Add(1)
		go func(tri *Triangle) {
			p0 := tri.P0.Dehom()
			p1 := tri.P1.Dehom()
			p2 := tri.P2.Dehom()

			minx := int(lin(min3(p0.X, p1.X, p2.X), -1, 1, 0, float64(width)))
			miny := int(lin(min3(p0.Y, p1.Y, p2.Y), -1, 1, 0, float64(height)))
			maxx := int(lin(max3(p0.X, p1.X, p2.X), -1, 1, 0, float64(width)))
			maxy := int(lin(max3(p0.Y, p1.Y, p2.Y), -1, 1, 0, float64(height)))

			for i := maxi(minx-1, 0); i < mini(maxx+1, width); i++ {
				for j := maxi(miny-1, 0); j < mini(maxy+1, height); j++ {
					coordx := lin(float64(i), 0, float64(width), -1, 1)
					coordy := lin(float64(j), 0, float64(height), -1, 1)

					screenCoord := &Vector2{coordx, coordy}
					if !In(p0, p1, p2, screenCoord) {
						continue
					}
					dePerp := tri.DePerp(&Vector2{coordx, coordy})
					zbuflock[i][j].Lock()
					if dePerp.Z <= 0 || (zbuf[i][j] != 0 && zbuf[i][j] <= dePerp.Z) {
						zbuflock[i][j].Unlock()
						continue
					}

					zbuf[i][j] = dePerp.Z
					u, v, w := tri.Bary(dePerp)
					norm := tri.N0.Scale(u).Add(tri.N1.Scale(v)).Add(tri.N2.Scale(w)).FastNormalize()
					im.Set(i, j, Render(tri.Material, norm, screenCoord.Hom().FastNormalize(), l, dePerp, &Vector2{u, v}).ToRGBA())

					zbuflock[i][j].Unlock()
				}
			}
			wg.Done()
		}(tri)
	}
	wg.Wait()
}

func DrawTrianglesParallelShadow(im *image.RGBA, t []*Triangle, l []Light) {
	width := im.Rect.Max.X - im.Rect.Min.X
	height := im.Rect.Max.Y - im.Rect.Min.Y
	zbuf := make([][]float64, width, width)
	zbuflock := make([][]sync.Mutex, width, width)
	wg := sync.WaitGroup{}
	for row := range zbuf {
		zbuf[row] = make([]float64, height, height)
		zbuflock[row] = make([]sync.Mutex, height, height)
	}
	var ct int32
	for _, tri := range t {
		wg.Add(1)
		go func(tri *Triangle) {
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
					dePerp := tri.DePerp(&Vector2{coordx, coordy})
					if !tri.In(dePerp) {
						continue
					}

					zbuflock[i][j].Lock()
					if dePerp.Z <= 0 || (zbuf[i][j] != 0 && zbuf[i][j] <= dePerp.Z) {
						zbuflock[i][j].Unlock()
						continue
					}

					zbuf[i][j] = dePerp.Z
					u, v, w := tri.Bary(dePerp)
					norm := tri.N0.Scale(u).Add(tri.N1.Scale(v)).Add(tri.N2.Scale(w)).Normalize()
					im.Set(i, j, RenderShadow(tri, t, tri.Material, norm, screenCoord.Hom().Normalize(), dePerp, l, &Vector2{u, v}).ToRGBA())
					zbuflock[i][j].Unlock()

				}
			}
			fmt.Println(atomic.AddInt32(&ct, 1), "of", len(t))
			wg.Done()
		}(tri)
	}
	wg.Wait()
}

func DrawTrianglesRayTracer(im *image.RGBA, t []*Triangle, l []Light) {
	width := im.Rect.Max.X - im.Rect.Min.X
	height := im.Rect.Max.Y - im.Rect.Min.Y
	wg := sync.WaitGroup{}
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			wg.Add(1)
			go func(i, j int) {
				coordx := lin(float64(i), 0, float64(width), -1, 1)
				coordy := lin(float64(j), 0, float64(height), -1, 1)

				screenCoord := &Vector2{coordx, coordy}
				im.Set(i, j, RayCast(t, l, screenCoord.Hom().Normalize(), 1).ToRGBA())
				wg.Done()
				fmt.Println(i, j)
			}(i, j)

		}
	}
	wg.Wait()
}
func ColorInterp(c0, c1, c2 *Color, s0, s1, s2 float64) *Color {
	return &Color{
		R: c0.R*s0 + c1.R*s1 + c2.R*s2,
		G: c0.G*s0 + c1.G*s1 + c2.G*s2,
		B: c0.B*s0 + c1.B*s1 + c2.B*s2,
		A: c0.A*s0 + c1.A*s1 + c2.A*s2,
	}
}

func ColorAdd(c0, c1 *Color) *Color {
	return &Color{
		R: c0.R + c1.R,
		G: c0.G + c1.G,
		B: c0.B + c1.B,
		A: c0.A + c1.A,
	}
}

func ColorMult(c0, c1 *Color) *Color {
	return &Color{
		R: c0.R * c1.R / 255,
		G: c0.G * c1.G / 255,
		B: c0.B * c1.B / 255,
		A: c0.A * c1.A / 255,
	}
}

func ColorScale(c *Color, s float64) *Color {
	return &Color{
		R: c.R * s,
		G: c.G * s,
		B: c.B * s,
		A: c.A,
	}
}

func OpenObj(filename string, rgba *Color) ([]*Triangle, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewScanner(f)
	reader.Split(bufio.ScanWords)
	var points []*Vector3
	var normals []*Vector3
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
		if t == "vn" {
			reader.Scan()
			xs := reader.Text()
			reader.Scan()
			ys := reader.Text()
			reader.Scan()
			zs := reader.Text()

			x, _ := strconv.ParseFloat(xs, 64)
			y, _ := strconv.ParseFloat(ys, 64)
			z, _ := strconv.ParseFloat(zs, 64)

			normals = append(normals, &Vector3{x, y, z})
		}
		if t == "f" {
			var xn, yn, zn string
			reader.Scan()
			val := strings.Split(reader.Text(), "/")
			xs := val[0]
			if len(val) == 3 {
				xn = val[2]
			} else {
				xn = "0"
			}
			reader.Scan()
			val = strings.Split(reader.Text(), "/")
			ys := val[0]
			if len(val) == 3 {
				yn = val[2]
			} else {
				yn = "0"
			}
			reader.Scan()
			val = strings.Split(reader.Text(), "/")
			zs := val[0]
			if len(val) == 3 {
				zn = val[2]
			} else {
				zn = "0"
			}
			x, _ := strconv.ParseInt(xs, 10, 32)
			y, _ := strconv.ParseInt(ys, 10, 32)
			z, _ := strconv.ParseInt(zs, 10, 32)
			xnn, _ := strconv.ParseInt(xn, 10, 32)
			ynn, _ := strconv.ParseInt(yn, 10, 32)
			znn, _ := strconv.ParseInt(zn, 10, 32)
			t := NewTriangle(points[x-1], points[y-1], points[z-1], &SolidMaterial{
				Color:         rgba,
				SpecColor_:    &Color{10, 10, 10, 255},
				SpecCoeff_:    8,
				AmbientCoeff_: .01,
			})
			if xnn != 0 {
				t.N0 = normals[xnn-1]
				t.N1 = normals[ynn-1]
				t.N2 = normals[znn-1]
			} else {
				t.N0 = t.Norm
				t.N1 = t.Norm
				t.N2 = t.Norm
			}
			triangles = append(triangles, t)

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
