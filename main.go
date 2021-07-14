package main

import (
	"go-wasm-demo/app"
	"syscall/js"
)

func main() {
	js.Global().Set("getMd5", app.GetMD5())
	js.Global().Set("imageRotate", make(map[string]interface{}))
	imageRotate := app.ImageRotate()
	js.Global().Get("imageRotate").Set("run", imageRotate["run"])
	js.Global().Get("imageRotate").Set("release", imageRotate["release"])

	js.Global().Set("getExcel", app.Excel())
	select {}
}
