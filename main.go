package main

import (
	"fmt"
	"syscall/js"
)

var renderer js.Value
var scene js.Value
var camera js.Value
var mesh js.Value
var animator js.Func

func main() {
	c := make(chan struct{}, 0)
	fmt.Println("Go/WASM main()")

	doc := js.Global().Get("document")
	app := doc.Call("getElementById", "app")
	cw := app.Get("clientWidth")
	ch := app.Get("clientHeight")

	// camera
	three := js.Global().Get("THREE")
	cameraCon := three.Get("PerspectiveCamera")
	camera = cameraCon.New(70, cw.Float()/ch.Float(), 0.01, 10)
	cpos := camera.Get("position")
	cpos.Set("z", 1)

	// geometry
	boxCon := three.Get("BoxGeometry")
	box := boxCon.New(0.2, 0.2, 0.2)
	matCon := three.Get("MeshNormalMaterial")
	mat := matCon.New()

	// scene
	sceneCon := three.Get("Scene")
	scene = sceneCon.New()

	// mesh
	meshCon := three.Get("Mesh")
	mesh = meshCon.New(box, mat)
	scene.Call("add", mesh)

	// renderer
	rendCon := three.Get("WebGLRenderer")
	renderer = rendCon.New("{antialias: true}")
	renderer.Call("setSize", cw.Int(), ch.Int())
	rdom := renderer.Get("domElement")
	app.Call("appendChild", rdom)

	// animator
	animator = js.FuncOf(animate)
	js.Global().Set("animate", animator)
	defer animator.Release()
	js.Global().Call("animate")

	fmt.Println("wait")
	<-c
	fmt.Println("exit")
}

func animate(this js.Value, args []js.Value) interface{} {
	win := js.Global()
	win.Call("requestAnimationFrame", animator)
	rot := mesh.Get("rotation")
	rx := rot.Get("x")
	ry := rot.Get("y")
	rot.Set("x", rx.Float()+0.01)
	rot.Set("y", ry.Float()+0.02)
	renderer.Call("render", scene, camera)
	return nil
}
