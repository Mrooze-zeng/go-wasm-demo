package app

import (
	"bytes"
	"syscall/js"
)

func getBuffer(args []js.Value) []byte {
	uint8Array := js.Global().Get("Uint8Array")
	if len(args) < 1 || !args[0].InstanceOf(uint8Array) {
		return nil
	}
	//todo 消耗大量时间
	buffer := make([]byte, args[0].Get("length").Int())
	js.CopyBytesToGo(buffer, args[0])

	return buffer
}

func isJPG(buffer []byte) bool {
	return bytes.HasPrefix(buffer, []byte{255, 216}) && bytes.HasSuffix(buffer, []byte{255, 217})
	// return binary.BigEndian.Uint16(buffer[0:2]) == 0xffd8 && binary.BigEndian.Uint16(buffer[len(buffer)-2:]) == 0xffd9
}
