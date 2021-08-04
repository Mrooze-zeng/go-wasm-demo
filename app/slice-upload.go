package app

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"syscall/js"
	"time"
)

type ChunkUploader struct {
	apiURL  string
	buffer  []byte
	md5     string
	name    string
	mintype string
	chunks  chan map[int][]byte
	dones   map[int]string
}

func NewChunkUploader(apiURL string, buffer []byte, name string, mintype string) *ChunkUploader {
	return &ChunkUploader{
		apiURL:  apiURL,
		buffer:  buffer,
		md5:     fmt.Sprintf("%x", md5.Sum(buffer)),
		name:    name,
		mintype: mintype,
		chunks:  make(chan map[int][]byte),
		dones:   make(map[int]string),
	}
}

func (c *ChunkUploader) setChunks(chunkLen, chunkSize int) {
	defer func() {
		close(c.chunks)
	}()
	for i := 1; i <= chunkLen; i++ {
		start := (i - 1) * chunkSize
		end := i * chunkSize
		if end > len(c.buffer) {
			end = len(c.buffer)
		}
		c.chunks <- map[int][]byte{
			i: c.buffer[start:end],
		}
	}
}

func (c *ChunkUploader) uploadAll() {
	for chunkMap := range c.chunks {
		for i, chunk := range chunkMap {
			// _, ok := c.dones[i]
			// if !ok {
			go c.uploadChunk(chunk, i)
			// } else {
			// 	fmt.Println("It has been upload....")
			// }
		}
	}
}

func (c *ChunkUploader) setDone(chunkMd5 string, index int) {
	c.dones[index] = chunkMd5
}

func (c *ChunkUploader) storeUploadedChunk() {
	localStorage := js.Global().Get("sessionStorage")
	// data := sessionStorage.Call("getItem", c.md5)
	fmt.Println(localStorage)
}

func (c *ChunkUploader) uploadChunk(chunk []byte, index int) {

	chunkMd5 := fmt.Sprintf("%x", md5.Sum(chunk))

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fileField, _ := writer.CreateFormField("file")
	fileField.Write([]byte(c.md5))

	lengthField, _ := writer.CreateFormField("size")
	lengthField.Write([]byte(fmt.Sprint(len(c.buffer))))

	md5Field, _ := writer.CreateFormField("md5")
	md5Field.Write([]byte(chunkMd5))

	indexField, _ := writer.CreateFormField("index")
	indexField.Write([]byte(fmt.Sprint(index)))

	chunkField, _ := writer.CreateFormFile("binary", "binary")
	io.Copy(chunkField, bytes.NewReader(chunk))

	writer.Close()
	req, _ := http.NewRequest("POST", c.apiURL, body)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, _ := http.DefaultClient.Do(req)
	data, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	// c.setDone(chunkMd5, index)
	fmt.Println(string(data), c.dones)
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

		// uploader.storeUploadedChunk()

		go uploader.setChunks(chunkLen, chunkSize)

		uploader.uploadAll()

		return nil
	})
}
