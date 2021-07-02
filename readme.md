# go wasm demo 学习


- 用 go md5 包在 web 前端生成文件或者字符的 MD5

- 随机旋转图片
### 已知问题：
-  在Tinygo中(用go不会出现内存大量占用的情况，且耗时比tinygo的更少，但打包后的文件很大)
```go
	buffer := make([]byte, args[0].Get("length").Int())
```
  大量的buffer会消耗大量时间,并堵塞其他程序运行(处理方法，将重复的文件缓存在变量中，减少重复占用内存)