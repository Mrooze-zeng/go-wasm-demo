package app

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"syscall/js"
)

func CompressGzip() map[string]js.Func {
	return map[string]js.Func{
		"gzip":   js.FuncOf(gzipFn),
		"ungzip": js.FuncOf(ungzipFn),
	}
}

func gzipFn(this js.Value, args []js.Value) interface{} {
	var res Result
	var buffer bytes.Buffer
	buf := getBuffer(args)
	name := args[1].String()

	gw, err := gzip.NewWriterLevel(&buffer, 9)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	gw.Name = name
	gw.Write(buf)

	if err := gw.Close(); err != nil {
		return nil
	}

	dst := exportDataToJS(buffer.Bytes())

	return res.new("application/x-gzip", dst)
}

func ungzipFn(this js.Value, args []js.Value) interface{} {
	var res Result
	buffer := getBuffer(args)
	r := bytes.NewReader(buffer)
	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil
	}
	defer gr.Close()

	buf, _ := ioutil.ReadAll(gr)

	dst := exportDataToJS(buf)

	res.new("application/x-gzip", dst)
	return res.add("name", gr.Name)
}
