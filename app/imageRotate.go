package app

import (
	"encoding/binary"
	"fmt"
	"syscall/js"
)

func ImageRotate() map[string]js.Func {
	var buffer []byte
	return map[string]js.Func{
		"run": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if len(args) != 2 {
				return js.Undefined()
			}
			direction := 1
			direction = args[1].Int()

			if buffer == nil {
				buffer = getBuffer(args)
			}

			if buffer == nil || !isJPG(buffer) {
				return js.Undefined()
			}
			b := rotate(buffer, byte(direction))
			if b == nil {
				return js.Undefined()
			}
			result := js.Global().Get("Uint8Array").New(len(b))
			js.CopyBytesToJS(result, b)
			fmt.Println("ok-----")
			return map[string]interface{}{
				"type":   "image/jpeg",
				"buffer": result,
				// "finish": time.Now().UnixNano() / 1e6,
			}
		}),
		"release": js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			if buffer != nil {
				fmt.Println("Release.....")
				buffer = nil
			}
			return nil
		}),
	}
}

func rotate(buffer []byte, direction byte) []byte {
	offset := 0
	length := len(buffer)
	for offset < length {
		if binary.BigEndian.Uint16(buffer[offset:offset+2]) == 0xffe1 {
			break
		} else {
			offset += 2
		}
	}

	if offset >= length {
		fmt.Println("没有找到APP1标识")
		return nil
	}

	app1_offset := offset
	exif_offset := app1_offset + 4

	if binary.BigEndian.Uint32(buffer[exif_offset:exif_offset+4]) != 0x45786966 {
		fmt.Println("无EXIF信息")
		return nil
	}

	tiff_offset := exif_offset + 6

	isLittle := binary.BigEndian.Uint16(buffer[tiff_offset:tiff_offset+2]) == 0x4949

	ifd0_offset := tiff_offset + int(binary.BigEndian.Uint32(buffer[tiff_offset+4:tiff_offset+8]))

	var entries_count int
	if isLittle {
		entries_count = int(binary.LittleEndian.Uint16(buffer[ifd0_offset : ifd0_offset+2]))
	} else {
		entries_count = int(binary.BigEndian.Uint16(buffer[ifd0_offset : ifd0_offset+2]))
	}

	entries_offset := ifd0_offset + 2

	for i := 0; i < entries_count; i++ {
		if isLittle {
			return nil
		} else {
			if binary.BigEndian.Uint16(buffer[entries_offset+i*12:entries_offset+i*12+2]) == 0x0112 {
				buffer[entries_offset+i*12+8+1] = direction
				// v := binary.BigEndian.Uint16(buffer[entries_offset+i*12+8 : entries_offset+i*12+10])
				// fmt.Println(entries_count, entries_offset, "bigEndian", v)
				return buffer
			}
		}
	}

	return nil

}
