package app

import (
	"crypto/md5"
	"encoding/hex"
	"syscall/js"
)

func GetMD5() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var res Result
		buffer := getBuffer(args)
		if buffer == nil {
			return js.Undefined()
		}
		md5Code := md5.Sum(buffer)
		return res.new("md5", hex.EncodeToString(md5Code[:]))
		// return res.new("md5", fmt.Sprintf("%x", md5.Sum(buffer)))
	})
}
