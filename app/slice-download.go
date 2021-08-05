package app

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"syscall/js"
)

func getChunk(url string, begin, end int, stop bool, chunks []byte, output chan Result) error {
	var isFinal bool

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Range", setRange(begin, end))
	res, _ := http.DefaultClient.Do(req)

	buf, _ := ioutil.ReadAll(res.Body)

	contentRange := res.Header.Get("Content-Range")
	contentRangeSlice := strings.Split(contentRange, "/")
	total, _ := strconv.Atoi(contentRangeSlice[1])
	next := 2*end - begin
	if next >= total {
		next = total
		isFinal = true
	}

	chunks = append(chunks, buf...)

	if stop {
		fmt.Println("end....")
		md5Code := md5.Sum(chunks)
		fmt.Println(hex.EncodeToString(md5Code[:]))
		var result Result
		disposition := res.Header.Get("Content-Disposition")
		result.new(res.Header.Get("Content-Type"), exportDataToJS(chunks))
		result.add("name", disposition[strings.Index(disposition, "\"")+1:strings.LastIndex(disposition, "\"")])
		output <- result
		return nil
	}
	go getChunk(url, end+1, next, isFinal, chunks, output)
	return nil
}

func setRange(begin, end int) string {
	beginstr := strconv.Itoa(begin)
	endstr := strconv.Itoa(end)
	return "bytes=" + beginstr + "-" + endstr
}

func SliceDownload() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return MyPromise(func() (res map[string]interface{}, err error) {

			url := args[0].JSValue().String()
			size := args[1].JSValue().Int()

			var chunks []byte
			output := make(chan Result, 1)

			go getChunk(url, 0, size, false, chunks, output)

			o := <-output

			return o.Value, err
		})
	})
}
