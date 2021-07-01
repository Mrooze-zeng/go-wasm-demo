package main

import (
	"go-wasm-demo/app"
	"syscall/js"
)

func main() {
	js.Global().Set("getMd5", app.GetMD5())
	js.Global().Set("imageRotate", app.ImageRotate())
	select {}
}
