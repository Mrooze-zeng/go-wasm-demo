package app

import "syscall/js"

func MyPromise(fn func() (res map[string]interface{}, err error)) js.Value {
	var handler js.Func
	handler = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer handler.Release()
		resolve := func(value interface{}) {
			args[0].Invoke(value)
		}

		reject := func(value interface{}) {
			args[1].Invoke(value)
		}
		go func() {
			res, err := fn()
			if err != nil {
				reject(js.Global().Get("Error").New(err.Error()))
			} else {
				resolve(res)
			}
		}()
		return js.Undefined()
	})
	return js.Global().Get("Promise").New(handler)
}
