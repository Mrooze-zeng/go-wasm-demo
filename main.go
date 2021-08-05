package main

import (
	"go-wasm-demo/app"
	"syscall/js"
)

func main() {
	// 文件md5
	js.Global().Set("getMd5", app.GetMD5())
	// 随机旋转图片
	js.Global().Set("imageRotate", make(map[string]interface{}))
	imageRotate := app.ImageRotate()
	js.Global().Get("imageRotate").Set("run", imageRotate["run"])
	js.Global().Get("imageRotate").Set("release", imageRotate["release"])
	// 生成excel测试文件
	js.Global().Set("getExcel", app.Excel())
	// 生成csv测试文件
	js.Global().Set("setCSV", app.Csv())
	// 生成图片缩略图
	js.Global().Set("getImageThumbnail", app.ImageThumbnail())
	// 处理视频文件
	js.Global().Set("parseVideo", app.ParseVideo())

	//gzip 解压
	compressGzip := app.CompressGzip()
	js.Global().Set("compress", make(map[string]interface{}))
	js.Global().Get("compress").Set("gzip", compressGzip["gzip"])
	js.Global().Get("compress").Set("ungzip", compressGzip["ungzip"])

	//切片上传
	js.Global().Set("sliceUpload", app.SliceUpload())

	//切片下载
	js.Global().Set("sliceDownload", app.SliceDownload())

	select {}
}
