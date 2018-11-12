package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"fmt"
	"github.com/wizgrao/raytracer/graphics"
	"flag"
)

var (
	outputFile = flag.String("o", "out.png", "Output File (png)")
	inputFile = flag.String("i", "in.obj", "Input file (png)")
	size = flag.Int("s", 2000, "Size of output image")
)

func main() {
	flag.Parse()
	image := image.NewRGBA(image.Rect(0, 0, *size, *size))
	californiaGold := &color.RGBA{196, 130, 15, 255}
	berkeleyBlue := &color.RGBA{0, 50, 98, 255}
	for i := 0; i < *size; i++ {
		for j := 0; j < *size; j++ {
			image.Set(i, j, berkeleyBlue)
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
	triangles, _ := graphics.OpenObj(*inputFile, californiaGold)
	fmt.Println(len(triangles))
	graphics.DrawTriangles(image, triangles, lit)
	f, _ := os.Create(*outputFile)
	png.Encode(f, image)
}
