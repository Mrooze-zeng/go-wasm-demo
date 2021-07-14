package app

import (
	"bytes"
	"encoding/csv"
	"syscall/js"
)

func Csv() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var buffer bytes.Buffer
		w := csv.NewWriter(&buffer)

		defer func() {
			w.Flush()
		}()

		data := [][]string{
			{"1", "中国", "23"},
			{"2", "美国", "23"},
			{"3", "bb", "23"},
			{"4", "bb", "23"},
			{"5", "bb", "23"},
		}

		w.WriteAll(data)

		dst := js.Global().Get("Uint8Array").New(len(buffer.Bytes()))

		js.CopyBytesToJS(dst, buffer.Bytes())

		return map[string]interface{}{
			"type": "text/csv",
			"data": dst,
		}
	})
}
