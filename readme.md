# go wasm demo 学习


- 用 go md5 包在 web 前端生成文件或者字符的 MD5

- 随机旋转图片，已知问题，
```go
	buffer := make([]byte, args[0].Get("length").Int())
```
  大量的buffer会消耗大量时间,并堵塞其他程序运行