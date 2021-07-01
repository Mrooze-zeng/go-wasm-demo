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
		// fmt.Println(isJPG(buffer))
		return fmt.Sprintf("%x", md5.Sum(buffer))
	})
}
