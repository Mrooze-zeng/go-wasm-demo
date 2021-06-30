export GOOS=js
export GOARCH=wasm
tiny_exec_path = $(shell tinygo env TINYGOROOT)/targets/wasm_exec.js
exec_path = $(shell go env GOROOT)/misc/wasm/wasm_exec.js


.ONESHELL:
build:;
	@go build -o ./demo//worker/app.wasm ./main.go
	@cp ${exec_path} ./demo/worker
	

.ONESHELL:
build-tiny:
	@tinygo build -o ./demo//worker/app.wasm -target wasm main.go
	@cp ${tiny_exec_path} ./demo/worker