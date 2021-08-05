package app

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"syscall/js"
)

func getChunk(url string, start, end int) (int, []byte, *Result) {

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Range", setRange(start, end))
	res, _ := http.DefaultClient.Do(req)
	chunk, _ := ioutil.ReadAll(res.Body)

	contentRange := res.Header.Get("Content-Range")
	contentRangeSlice := strings.Split(contentRange, "/")
	total, _ := strconv.Atoi(contentRangeSlice[1])

	fmt.Println(contentRange)

	var result Result
	disposition := res.Header.Get("Content-Disposition")
	result.new(res.Header.Get("Content-Type"), "")
	result.add("name", disposition[strings.Index(disposition, "\"")+1:strings.LastIndex(disposition, "\"")])

	return total, chunk, &result
}

func setRange(begin, end int) string {
	beginstr := strconv.Itoa(begin)
	endstr := strconv.Itoa(end)
	return "bytes=" + beginstr + "-" + endstr
}

func SliceDownload() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		url := args[0].JSValue().String()
		size := args[1].JSValue().Int()
		firstChunkSize := 15
		return MyPromise(func() (res map[string]interface{}, err error) {

			total, firstChunk, output := getChunk(url, 0, firstChunkSize)

			var buffer []byte
			var wg sync.WaitGroup

			result := make(map[int][]byte)
			result[0] = firstChunk

			isBreak := false

			length := (total-firstChunkSize)/size + 1
			if (total-firstChunkSize)%size > 0 {
				length++
			}

			for i := 1; i < length; i++ {
				start := firstChunkSize + (i-1)*size
				end := start + size
				if end >= total {
					end = total
					isBreak = true
				}

				wg.Add(1)

				go func(start, end, index int) {
					_, chunk, _ := getChunk(url, start, end)
					result[index] = chunk
					defer wg.Done()
				}(start+1, end, i)
				if isBreak {
					break
				}
			}

			wg.Wait()

			for i := 0; i < len(result); i++ {

				buffer = append(buffer, result[i]...)
			}

			md5Code := md5.Sum(buffer)

			output.add("data", exportDataToJS(buffer))
			output.add("md5", hex.EncodeToString(md5Code[:]))
			output.add("size", len(buffer))

			return output.Value, err
		})
	})
}
