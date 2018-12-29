package main

import (
	"fmt"
	"net/http"
	"flag"
	"math"
	"github.com/wizgrao/simple3d/graphics"
	"image"
	"image/color"
	"bytes"
	"log"
	"strconv"
	"image/png"
)

type helloHandler struct {T []*graphics.Triangle}
type otherHandler struct {T []*graphics.Triangle}
var (
	inputFile = flag.String("i", "in.obj", "Input file (png)")
	size = flag.Int("s", 2000, "Size of output image")
	xt = flag.Float64("xt", 0, "Translation in X direction")
	yt = flag.Float64("yt", 0, "Translation in Y direction")
	zt = flag.Float64("zt", 1.8, "Translation in Z direction")
	xr = flag.Float64("xr", 0, "Rotation in X direction")
	yr = flag.Float64("yr", math.Pi, "Rotation in Y direction")
	zr = flag.Float64("zr", math.Pi, "Rotation in Z direction")
	port = flag.String("p", "8080", "port")

)

func parseForm(r *http.Request) (xrp, yrp, zrp, xtp, ytp, ztp float64) {
	r.ParseForm()
	xrs := r.Form.Get("xr")
	yrs := r.Form.Get("yr")
	zrs := r.Form.Get("zr")
	xts := r.Form.Get("xt")
	yts := r.Form.Get("yt")
	zts := r.Form.Get("zt")

	if x, err := strconv.ParseFloat(xrs, 64); err == nil{
		xrp = x
	}
	if x, err := strconv.ParseFloat(yrs, 64); err == nil{
		yrp = x
	}
	if x, err := strconv.ParseFloat(zrs, 64); err == nil{
		zrp = x
	}

	if x, err := strconv.ParseFloat(xts, 64); err == nil{
		xtp = x
	}
	if x, err := strconv.ParseFloat(yts, 64); err == nil{
		ytp = x
	}
	if x, err := strconv.ParseFloat(zts, 64); err == nil {
		ztp = x
	}
	return
}

func (h* helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request")
	xrp, yrp, zrp, xtp, ytp, ztp := parseForm(r)
	t := graphics.ApplyTransform(h.T, graphics.Translate(*xt + xtp, *yt + ytp, *zt + ztp).
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
	lit := &graphics.DirectionLight{
		Direction: (&graphics.Vector3{1, 1, 1}).Normalize(),
		Color: &color.RGBA{
			R:255,
			G:255,
			B:255,
			A:255,
		},
	}
	graphics.DrawTrianglesParallel(im, t, []graphics.Light{lit})

	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, im); err != nil {
		log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
	fmt.Println("done")

}

func (h* otherHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request")
	xrp, yrp, zrp, xtp, ytp, ztp := parseForm(r)
	t := graphics.ApplyTransform(h.T, graphics.Translate(*xt + xtp, *yt + ytp, *zt + ztp).
		Mult(graphics.RotZ(zrp)).
		Mult(graphics.RotY(yrp)).
		Mult(graphics.RotX(xrp)))
	im := image.NewRGBA(image.Rect(0, 0, *size, *size))
	bg := &color.RGBA{0, 0, 0, 255}
	for i := 0; i < *size; i++ {
		for j := 0; j < *size; j++ {
			im.Set(i, j, bg)
		}
	}
	lit1 := &graphics.PointLight{
		Location: &graphics.Vector3{2, 1, 0},
		R: 500,
	}
	lit2 := &graphics.PointLight{
		Location: &graphics.Vector3{-2, 1, 0},
		B: 500,
	}
	lit3 := &graphics.PointLight{
		Location: &graphics.Vector3{0, 1, -1},
		G: 500,
	}
	graphics.DrawTrianglesParallel(im, t, []graphics.Light{lit1, lit2, lit3})

	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, im); err != nil {
		log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
	fmt.Println("done")

}

func main() {
	flag.Parse()
	fg := &color.RGBA{255, 255, 255, 255}
	triangles, _ := graphics.OpenObj(*inputFile, fg)
	transform := graphics.RotZ(*zr).
		Mult(graphics.RotY(*yr)).
		Mult(graphics.RotX(*xr))
	triangles = graphics.ApplyTransform(triangles, transform)
	handler := &helloHandler{triangles}
	http.Handle("/points", &otherHandler{triangles})
	http.Handle("/", handler)
	fmt.Println("Listening on port "+*port)
	fmt.Print(http.ListenAndServe(":" + *port, nil))
}
