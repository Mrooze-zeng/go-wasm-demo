package app

import (
	"bytes"
	"syscall/js"
)

func getBuffer(args []js.Value) []byte {
	uint8Array := js.Global().Get("Uint8Array")
	uint8ClampedArray := js.Global().Get("Uint8ClampedArray")
	if len(args) < 1 || !(args[0].InstanceOf(uint8Array) || args[0].InstanceOf(uint8ClampedArray)) {
		return nil
	}
	dst := make([]byte, args[0].Get("length").Int())
	js.CopyBytesToGo(dst, args[0])
	return dst
}

func isJPG(buffer []byte) bool {
	return bytes.HasPrefix(buffer, []byte{255, 216}) && bytes.HasSuffix(buffer, []byte{255, 217})
	// return binary.BigEndian.Uint16(buffer[0:2]) == 0xffd8 && binary.BigEndian.Uint16(buffer[len(buffer)-2:]) == 0xffd9
}

func exportDataToJS(buffer []byte) js.Value {
	dst := js.Global().Get("Uint8Array").New(len(buffer))
	js.CopyBytesToJS(dst, buffer)
	return dst
}

type Result struct {
	Value map[string]interface{}
}

func (r *Result) new(dataType string, data interface{}) map[string]interface{} {
	r.Value = make(map[string]interface{})
	r.Value["type"] = dataType
	r.Value["data"] = data
	return r.Value
}

func (r Result) add(name string, value interface{}) map[string]interface{} {
	r.Value[name] = value
	return r.Value
}

func (r Result) remove(name string) map[string]interface{} {
	delete(r.Value, name)
	return r.Value
}
