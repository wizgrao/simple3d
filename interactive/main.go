package main

import (
	"fmt"
	"net/http"
	"flag"
	"math"
	"github.com/wizgrao/raytracer/graphics"
	"image"
	"image/color"
	"bytes"
	"log"
	"strconv"
	"image/png"
)

type helloHandler struct {T []*graphics.Triangle}
var (
	inputFile = flag.String("i", "in.obj", "Input file (png)")
	size = flag.Int("s", 2000, "Size of output image")
	xt = flag.Float64("xt", 0, "Translation in X direction")
	yt = flag.Float64("yt", 0, "Translation in Y direction")
	zt = flag.Float64("zt", 1.8, "Translation in Z direction")
	xr = flag.Float64("xr", 0, "Rotation in X direction")
	yr = flag.Float64("yr", math.Pi, "Rotation in Y direction")
	zr = flag.Float64("zr", math.Pi, "Rotation in Z direction")

)
func (h* helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	xrs := r.Form.Get("xr")
	yrs := r.Form.Get("yr")
	zrs := r.Form.Get("zr")

	var xrp float64
	var yrp float64
	var zrp float64

	if x, err := strconv.ParseFloat(xrs, 64); err == nil{
		xrp = x
	}
	if x, err := strconv.ParseFloat(yrs, 64); err == nil{
		yrp = x
	}
	if x, err := strconv.ParseFloat(zrs, 64); err == nil{
		zrp = x
	}
	fmt.Println(xrp, yrp, zrp)
	t := graphics.ApplyTransform(h.T, graphics.Translate(*xt, *yt, *zt).
		Mult(graphics.RotZ(zrp)).
		Mult(graphics.RotY(yrp)).
		Mult(graphics.RotX(xrp)))
	im := image.NewRGBA(image.Rect(0, 0, *size, *size))
	berkeleyBlue := &color.RGBA{0, 50, 98, 255}
	for i := 0; i < *size; i++ {
		for j := 0; j < *size; j++ {
			im.Set(i, j, berkeleyBlue)
		}
	}
	lit := &graphics.Light{
		Norm: (&graphics.Vector3{0, 1, 1}).Normalize(),
		C: &color.RGBA{
			R:255,
			G:255,
			B:255,
			A:255,
		},
	}
	graphics.DrawTriangles(im, t, lit)

	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, im); err != nil {
		log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}

}
func main() {
	flag.Parse()
	californiaGold := &color.RGBA{196, 130, 15, 255}
	triangles, _ := graphics.OpenObj(*inputFile, californiaGold)
	transform := graphics.RotZ(*zr).
		Mult(graphics.RotY(*yr)).
		Mult(graphics.RotX(*xr))
	triangles = graphics.ApplyTransform(triangles, transform)
	handler := &helloHandler{triangles}
	http.Handle("/", handler)
	fmt.Print(http.ListenAndServe(":8080", nil))
}