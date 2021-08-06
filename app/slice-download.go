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

func getChunk(url string, start, end int) (int, []byte, *http.Header) {

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Range", setRange(start, end))
	res, _ := http.DefaultClient.Do(req)
	chunk, _ := ioutil.ReadAll(res.Body)

	defer res.Body.Close()

	contentRange := res.Header.Get("Content-Range")
	contentRangeSlice := strings.Split(contentRange, "/")
	total, _ := strconv.Atoi(contentRangeSlice[1])

	fmt.Println(contentRange)

	return total, chunk, &res.Header
}

func setRange(begin, end int) string {
	beginstr := strconv.Itoa(begin)
	endstr := strconv.Itoa(end)
	return "bytes=" + beginstr + "-" + endstr
}

func splitRequestProcess(s, l, firstChunkSize, size, total int, wg sync.WaitGroup, url string, result map[int][]byte) {
	isBreak := false

	for i := s; i < l; i++ {
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
}

func SliceDownload() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		url := args[0].JSValue().String()
		size := args[1].JSValue().Int()
		firstChunkSize := 15
		return MyPromise(func() (res map[string]interface{}, err error) {
			var output Result
			var buffer []byte
			var wg sync.WaitGroup

			total, firstChunk, baseHeaders := getChunk(url, 0, firstChunkSize)

			result := make(map[int][]byte)
			result[0] = firstChunk

			length := (total-firstChunkSize)/size + 1
			if (total-firstChunkSize)%size > 0 {
				length++
			}

			if length/6 > 1 {
				l := length / 6
				if length%6 > 0 {
					l++
				}
				for i := 0; i < l; i++ {
					splitRequestProcess(i*6+1, (i+1)*6+1, firstChunkSize, size, total, wg, url, result)
				}
			} else {
				splitRequestProcess(1, length, firstChunkSize, size, total, wg, url, result)
			}

			for i := 0; i < len(result); i++ {

				buffer = append(buffer, result[i]...)
			}

			md5Code := md5.Sum(buffer)

			disposition := baseHeaders.Get("Content-Disposition")
			output.new(baseHeaders.Get("Content-Type"), "")
			output.add("name", disposition[strings.Index(disposition, "\"")+1:strings.LastIndex(disposition, "\"")])

			output.add("data", exportDataToJS(buffer))
			output.add("md5", hex.EncodeToString(md5Code[:]))
			output.add("size", len(buffer))

			return output.Value, err
		})
	})
}
