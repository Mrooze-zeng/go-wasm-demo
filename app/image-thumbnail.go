package app

import (
	"bytes"
	"image/jpeg"
	"syscall/js"

	"github.com/disintegration/imaging"
)

func ImageThumbnail() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var res Result
		buf := getBuffer(args)
		r := bytes.NewReader(buf)
		// img, _, err := image.Decode(r)
		// if err != nil {
		// 	return nil
		// }

		// image := resize.Resize(160, 0, img, resize.Lanczos3)
		srcImage, err := imaging.Decode(r)

		if err != nil {
			return nil
		}

		image := imaging.Resize(srcImage, 250, 0, imaging.Lanczos)

		result := new(bytes.Buffer)

		err = jpeg.Encode(result, image, nil)

		if err != nil {
			return nil
		}

		dst := exportDataToJS(result.Bytes())
		return res.new("image/jpeg", dst)
	})
}
