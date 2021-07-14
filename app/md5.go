package app

import (
	"crypto/md5"
	"fmt"
	"syscall/js"
)

func GetMD5() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		buffer := getBuffer(args)
		if buffer == nil {
			return js.Undefined()
		}
		return map[string]interface{}{
			"type": "md5",
			"data": fmt.Sprintf("%x", md5.Sum(buffer)),
		}
	})
}
