package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

type FileChunks struct {
	md5    string
	chunks map[int][]byte
}

func NewFileChunks(md5 string) *FileChunks {
	return &FileChunks{
		md5:    md5,
		chunks: make(map[int][]byte),
	}
}

func (f *FileChunks) add(chunk []byte, index int) {
	if f.chunks == nil {
		f.chunks = make(map[int][]byte)
	}
	if !f.has(index) {
		f.chunks[index] = chunk
	}
}

func (f *FileChunks) has(index int) bool {
	_, ok := f.chunks[index]
	return ok
}

func (f *FileChunks) concat() []byte {
	var res []byte
	for i := 0; i < len(f.chunks); i++ {
		res = append(res, f.chunks[i+1]...)
	}
	return res
}

func (f *FileChunks) isComplete(size int) bool {
	total := 0
	for _, c := range f.chunks {
		total += len(c)
	}
	return total == size
}

type FileStore struct {
	mu    sync.Mutex
	files map[string]FileChunks
}

func NewFileStore() *FileStore {
	return &FileStore{
		mu:    sync.Mutex{},
		files: make(map[string]FileChunks),
	}
}

func (t *FileStore) add(md5 string, chunk []byte, index int) {
	_, ok := t.files[md5]
	if !ok {
		t.files[md5] = *NewFileChunks(md5)
	}
	fc := t.files[md5]
	fc.add(chunk, index)
}
func (t *FileStore) remove(md5 string) {
	_, ok := t.files[md5]
	if ok {
		delete(t.files, md5)
	}
}

var fileStore = NewFileStore()

func cors(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")                                                            // 允许访问所有域，可以换成具体url，注意仅具体url才能带cookie信息
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token") //header的类型
		w.Header().Add("Access-Control-Allow-Credentials", "true")                                                    //设置为true，允许ajax异步请求带cookie信息
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")                             //允许请求方法
		// w.Header().Set("content-type", "application/jsjon;charset=UTF-8")                                              //返回数据格式是json
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		f(w, r)
	}
}

func getMd5String(b []byte) string {
	md5Code := md5.Sum(b)
	return hex.EncodeToString(md5Code[:])
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		fmt.Fprintf(w, "Hello World!!")
		return
	}
	file, _, err := r.FormFile("binary")
	if err != nil {
		fmt.Println("Error retrieving the file", err)
		return
	}
	defer file.Close()

	b, _ := ioutil.ReadAll(file)

	fileMd5 := r.FormValue("file")
	sizeStr := r.FormValue("size")
	size, _ := strconv.Atoi(sizeStr)
	chunkIndexStr := r.FormValue("index")
	chunkIndex, _ := strconv.Atoi(chunkIndexStr)

	// fileStore.mu.Lock()
	// defer fileStore.mu.Unlock()

	fileStore.add(fileMd5, b, chunkIndex)

	fc := fileStore.files[fileMd5]

	if fc.isComplete(size) {
		defer fileStore.remove(fileMd5)
		buf := fc.concat()
		md5String := getMd5String(buf)
		fmt.Fprintf(w, "Received file md5:%s-----Send file md5:%s----%s", md5String, fileMd5, chunkIndexStr)
		return
	}

	md5String := getMd5String(b)
	fmt.Fprintf(w, "Received chunk md5:%s-----Send file md5:%s----%s", md5String, fileMd5, chunkIndexStr)
}

func main() {

	http.HandleFunc("/", cors(indexHandler))
	http.HandleFunc("/upload", cors(uploadHandler))

	http.ListenAndServe(":8080", nil)
}
