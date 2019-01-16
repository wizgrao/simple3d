package main

import (
	"net/http"

	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/wizgrao/blow/gorillaconnection"
	"github.com/wizgrao/blow/maps"
	"github.com/wizgrao/simple3d/graphics"
	"image"
	"io"
	"os"
	"image/png"
)

type wasmHandler int

var wasm wasmHandler

var upgrader = websocket.Upgrader{}

var mapper *graphics.RayTraceMapper

var pool *maps.WorkerPool

type newConnectionHandler int

func (newConnectionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading http", err)
		return
	}
	connection := &gorillaconnection.Connection{c}
	fmt.Println("New Connection")
	pool.AddWorker(connection)
	select {}
}

func (h wasmHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/wasm")
	f, _ := os.Open("slave/main.wasm")
	p := make([]byte, 4)
	for {
		n, err := f.Read(p)
		if err == io.EOF {
			break
		}
		w.Write(p[:n])
	}

}

func main() {
	flag.Parse()
	im := image.NewRGBA(image.Rect(0, 0, 128, 128))
	writer := &graphics.WriterMapper{im, 0, 128* 128}
	source := &graphics.PixelSource{4, im}
	var nch newConnectionHandler
	pool = maps.NewWorkerPool()
	pool.Register(mapper)
	go func() {
		maps.GeneratorSource(source, pool).MapDispatch(mapper).MapLocal(writer).Sink()
		f, _ := os.Create("out.png")
		png.Encode(f, im)
	}()
	http.Handle("/main.wasm", wasm)
	http.Handle("/sock", nch)
	http.Handle("/", http.FileServer(http.Dir("slave/")))
	fmt.Print(http.ListenAndServe(":8090", nil))
}
