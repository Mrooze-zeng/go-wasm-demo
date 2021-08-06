package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type FileCollection struct {
	md5    string
	chunks map[int][]byte
}

func NewFileCollection(md5 string) *FileCollection {
	return &FileCollection{
		md5:    md5,
		chunks: make(map[int][]byte),
	}
}

func (f *FileCollection) add(chunk []byte, index int) {
	if f.chunks == nil {
		f.chunks = make(map[int][]byte)
	}
	f.chunks[index] = chunk
}

func (f *FileCollection) concat() []byte {
	var res []byte
	for i := 0; i < len(f.chunks); i++ {
		res = append(res, f.chunks[i+1]...)
	}
	return res
}

func (f *FileCollection) isComplete(size int) bool {
	total := 0
	for _, c := range f.chunks {
		total += len(c)
	}
	return total == size
}

type TempCollection struct {
	files map[string]FileCollection
}

func NewTempCollection() *TempCollection {
	return &TempCollection{
		files: make(map[string]FileCollection),
	}
}

func (t *TempCollection) add(md5 string, chunk []byte, index int) {
	_, ok := t.files[md5]
	if !ok {
		t.files[md5] = *NewFileCollection(md5)
	}
	fc := t.files[md5]
	fc.add(chunk, index)
}
func (t *TempCollection) remove(md5 string) {
	_, ok := t.files[md5]
	if ok {
		delete(t.files, md5)
	}
}

var TC = NewTempCollection()

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
	r.ParseMultipartForm(10 << 20)
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

	TC.add(fileMd5, b, chunkIndex)

	fc := TC.files[fileMd5]

	if fc.isComplete(size) {
		defer TC.remove(fileMd5)
		buf := fc.concat()
		md5String := getMd5String(buf)
		fmt.Fprintf(w, "Received file md5:%s-----Send file md5:%s", md5String, fileMd5)
		return
	}

	md5String := getMd5String(b)
	fmt.Fprintf(w, "Received chunk md5:%s-----Send file md5:%s", md5String, fileMd5)
}

func main() {

	http.HandleFunc("/", cors(indexHandler))
	http.HandleFunc("/upload", cors(uploadHandler))

	http.ListenAndServe(":8080", nil)
}
