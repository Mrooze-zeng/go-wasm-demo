# go wasm demo 学习

```go 
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
```


- 用 go md5 包在 web 前端生成文件或者字符的 MD5

- 随机旋转图片
- 生成 excel/csv 测试文件并下载
- 生成图片缩略图
- 压缩上传文件和解压上传文件
- 前端文件切片上传

### 已知问题：

- 在 Tinygo 中(用 go 不会出现内存大量占用的情况，且耗时比 tinygo 的更少，但打包后的文件很大)

```go
	buffer := make([]byte, args[0].Get("length").Int())
```

大量的 buffer 会消耗大量时间,并堵塞其他程序运行(处理方法，将重复的文件缓存在变量中，减少重复占用内存)

- 在 TinyGo 中，不支持 encoding/xml,生成 excel 不可用 TinyGo
