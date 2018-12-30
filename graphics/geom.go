package graphics

import (
	"math"
)

type Vector4 struct {
	X []float64
}

type Mat4 struct {
	X []float64
}

func (m *Mat4)At(i, j int) float64{
	return m.X[i*4 + j]
}

func (m *Mat4)Set(i, j int, val float64) {
	m.X[i*4 + j] = val
}

func NewMat4() *Mat4{
	m := Mat4{}
	m.X = make([]float64, 16, 16)
	for i:=0; i < 16; i+= 5 {
		m.X[i] = 1
	}
	return &m
}

func (m *Mat4) Mult (n *Mat4)  *Mat4{
	res := NewMat4()
	for i:=0; i < 4; i++ {
		for j:= 0; j < 4; j++ {
			var sum float64
			for k :=0; k < 4; k++ {
				sum += m.At(i, k)*n.At(k, j)
			}
			res.Set(i, j, sum)
		}
	}
	return res
}

func RotZ(theta float64) *Mat4 {
	c := math.Cos(theta)
	s := math.Sin(theta)

	res := NewMat4()

	res.X = []float64{
		c, -s, 0, 0,
		s, c, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
	return res
}

func RotY(theta float64) *Mat4 {
	c := math.Cos(theta)
	s := math.Sin(theta)

	res := NewMat4()

	res.X = []float64{
		c, 0, -s, 0,
		0, 1, 0, 0,
		s, 0, c, 0,
		0, 0, 0, 1,
	}
	return res
}

func RotX(theta float64) *Mat4 {
	c := math.Cos(theta)
	s := math.Sin(theta)

	res := NewMat4()

	res.X = []float64{
		1, 0, 0, 0,
		0, c, -s, 0,
		0, s, c, 0,
		0, 0, 0, 1,
	}
	return res
}

func Translate(x, y, z float64) *Mat4 {
	res := NewMat4()

	res.X = []float64{
		1, 0, 0, x,
		0, 1, 0, y,
		0, 0, 1, z,
		0, 0, 0, 1,
	}
	return res
}



type Vector3 struct {
	X float64
	Y float64
	Z float64
}

type Vector2 struct {
	X float64
	Y float64
}

func (v *Vector2) Hom() *Vector3{
	return &Vector3{
		X:v.X,
		Y:v.Y,
		Z:1,
	}
}

type Triangle struct {
	P0   *Vector3
	P1   *Vector3
	P2   *Vector3
	N0   *Vector3
	N1   *Vector3
	N2   *Vector3

	Norm *Vector3
	Material

}

func (t *Triangle) DePerp(v *Vector2) *Vector3 {
	d := t.Norm.Dot(t.P0) //ax + by + cz = d
	vhom := v.Hom()
	coeff := vhom.Dot(t.Norm)
	z := d/coeff
	return vhom.Scale(z)
}


func (v *Vector3) Dehom() *Vector2 {
	return &Vector2{
		X: v.X / v.Z,
		Y: v.Y / v.Z,
	}
}



func (v *Vector3) Hom() *Vector4 {
	return &Vector4{
		X: []float64{v.X, v.Y, v.Z, 1},
	}
}

func (v *Vector3) Ext() *Vector4 {
	return &Vector4{
		X: []float64{v.X, v.Y, v.Z, 0},
	}
}

func (v *Vector4) Dehom() *Vector3 {
	return &Vector3{
		X: v.X[0]/v.X[3],
		Y: v.X[1]/v.X[3],
		Z: v.X[2]/v.X[3],
	}
}


func (v *Vector4) Unex() *Vector3 {
	return &Vector3{
		X: v.X[0],
		Y: v.X[1],
		Z: v.X[2],
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

func (u *Vector2) Add(v *Vector2) *Vector2 {
	return &Vector2{
		X: u.X + v.X,
		Y: u.Y + v.Y,
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
		Y: -(v1.X*v2.Z - v1.Z*v2.X),
		Z: v1.X*v2.Y - v1.Y*v2.X,
	}
}

func (m *Mat4) Dot(v *Vector4)  *Vector4 {
	x := make([]float64 , 4, 4)
	for i:=0; i < 4; i++ {
		var sum float64
		for j:=0; j<4; j++ {
			sum += m.At(i, j) * v.X[j]
		}
		x[i] = sum
	}
	return &Vector4{X:x}
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

func (v *Vector2) Scale(s float64) *Vector2 {
	return &Vector2{
		X: v.X * s,
		Y: v.Y * s,
	}
}



func CalcNorm(p0, p1, p2 *Vector3) *Vector3 {
	res := Cross(p0.Sub(p1), p0.Sub(p2)).Normalize()
	return res
}

func (t *Triangle) Bary(v *Vector3)  (float64, float64, float64) {
	t1 := (&Triangle{
		P0: v,
		P1: t.P1,
		P2: t.P2,
	}).Area()
	t2 := (&Triangle{
		P0: t.P0,
		P1: v,
		P2: t.P2,
	}).Area()
	t3 := (&Triangle{
		P0: t.P0,
		P1: t.P1,
		P2: v,
	}).Area()
	a := t.Area()
	return t1/a, t2/a, t3/a
}

func (t *Triangle) Area() float64 {
	return Cross(t.P0.Sub(t.P1), t.P0.Sub(t.P2)).Norm()*0.5
}

func (t *Triangle) In(vx *Vector3) bool{
	u, v, w := t.Bary(vx)
	return u + v + w <= 1.0001
}
func NewTriangle(p0, p1, p2 *Vector3, m Material) *Triangle {
	norm := CalcNorm(p0, p1, p2)
	return &Triangle{
		P0:   p0,
		P1:   p1,
		P2:   p2,
		Norm: norm,
		Material:   m,
		N0: norm,
		N1: norm,
		N2: norm,
	}
}


func ApplyTransform(triangles []*Triangle, mat *Mat4) []*Triangle{
	res := make([]*Triangle, len(triangles), len(triangles))
	for i, t := range triangles {
		res[i] = NewTriangle(
			mat.Dot(t.P0.Hom()).Dehom(),
			mat.Dot(t.P1.Hom()).Dehom(),
			mat.Dot(t.P2.Hom()).Dehom(),
			t.Material,
		)
		res[i].N0 = mat.Dot(t.N0.Ext()).Unex()
		res[i].N1 = mat.Dot(t.N1.Ext()).Unex()
		res[i].N2 = mat.Dot(t.N2.Ext()).Unex()



	}
	return res
}

