package app

import (
	"bytes"
	"encoding/csv"
	"syscall/js"
)

func Csv() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var buffer bytes.Buffer
		var res Result

		w := csv.NewWriter(&buffer)

		defer w.Flush()

		data := [][]string{
			{"1", "中国", "23"},
			{"2", "美国", "23"},
			{"3", "bb", "23"},
			{"4", "bb", "23"},
			{"5", "bb", "23"},
		}

		err := w.WriteAll(data)

		if err != nil {
			return nil
		}

		dst := exportDataToJS(buffer.Bytes())

		return res.new("text/csv", dst)
	})
}
