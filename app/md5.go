package app

import (
	"crypto/md5"
	"fmt"
	"syscall/js"
)

func GetMD5() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		uint8Array := js.Global().Get("Uint8Array")
		if len(args) < 1 || !args[0].InstanceOf(uint8Array) {
			return js.Undefined()
		}
		buffer := make([]byte, args[0].Get("length").Int())
		js.CopyBytesToGo(buffer, args[0])
		return fmt.Sprintf("%x", md5.Sum(buffer))
	})
}
