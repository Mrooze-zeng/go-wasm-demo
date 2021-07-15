package app

import (
	"syscall/js"
)

func ParseVideo() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return nil
	})
}
