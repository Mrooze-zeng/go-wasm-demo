package app

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
	"syscall/js"
	"time"
)

type ChunkUploader struct {
	apiURL  string
	buffer  []byte
	md5     string
	name    string
	mintype string
}

func NewChunkUploader(apiURL string, buffer []byte, name string, mintype string) *ChunkUploader {
	md5Code := md5.Sum(buffer)
	return &ChunkUploader{
		apiURL:  apiURL,
		buffer:  buffer,
		md5:     hex.EncodeToString(md5Code[:]),
		name:    name,
		mintype: mintype,
	}
}

func (c *ChunkUploader) getChunksInfo() {

}

func (c *ChunkUploader) uploadChunks(chunkLen, chunkSize int) {
	for i := 1; i <= chunkLen; i++ {
		start := (i - 1) * chunkSize
		end := start + chunkSize
		if end > len(c.buffer) {
			end = len(c.buffer)
		}
		go c.uploadChunk(c.buffer[start:end], i)
	}
}

func (c *ChunkUploader) uploadChunk(chunk []byte, index int) {

	md5Code := md5.Sum(chunk)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("file", c.md5)
	writer.WriteField("size", strconv.Itoa(len(c.buffer)))
	writer.WriteField("md5", hex.EncodeToString(md5Code[:]))
	writer.WriteField("index", strconv.Itoa(index))

	fmt.Println(hex.EncodeToString(md5Code[:]))

	chunkField, _ := writer.CreateFormFile("binary", "binary")
	io.Copy(chunkField, bytes.NewReader(chunk))

	writer.Close()
	req, _ := http.NewRequest("POST", c.apiURL, body)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	data, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	fmt.Println(string(data))
}

func SliceUpload() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		s := time.Now().UnixNano() / 1e6
		defer func() {
			fmt.Println(time.Now().UnixNano()/1e6 - s)
		}()

		options := args[1].JSValue()
		apiURL := args[2].JSValue().String()
		buffer := getBuffer(args)
		name := options.Get("name").String()
		mintype := options.Get("mintype").String()
		chunkSize := options.Get("chunkSize").Int()
		var chunkLen = len(buffer) / chunkSize
		if len(buffer)%chunkSize > 0 {
			chunkLen += 1
		}

		uploader := NewChunkUploader(apiURL, buffer, name, mintype)

		go uploader.uploadChunks(chunkLen, chunkSize)

		return nil
	})
}
