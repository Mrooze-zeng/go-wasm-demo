package app

import (
	"syscall/js"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func Excel() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		f := excelize.NewFile()

		index := f.NewSheet("Sheet2")

		f.SetCellValue("Sheet2", "A2", "Hello world.")
		f.SetCellValue("Sheet1", "B2", 100)

		f.SetActiveSheet(index)

		buf, err := f.WriteToBuffer()

		if err != nil {
			return nil
		}

		dst := js.Global().Get("Uint8Array").New(len(buf.Bytes()))

		js.CopyBytesToJS(dst, buf.Bytes())

		return map[string]interface{}{
			"type": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			"data": dst,
		}
	})
}
