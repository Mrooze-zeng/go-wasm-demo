tiny_exec_path = $(shell tinygo env TINYGOROOT)/targets/wasm_exec.js
exec_path = $(shell go env GOROOT)/misc/wasm/wasm_exec.js


.ONESHELL:
build:
	@GOOS=js GOARCH=wasm go build -o ./demo/worker/app.wasm .
	@cp ${exec_path} ./demo/worker
	

.ONESHELL:
build-tiny:
	@GOOS=js GOARCH=wasm tinygo build -o ./demo/worker/app.wasm -target wasm .
	@cp ${tiny_exec_path} ./demo/worker

.ONESHELL:
serve:
	@go run server/main.go