package app

/*
#include <stdio.h>

void printint(int v) {
    printf("printint: %d\n", v);
}
*/
import "C"

import (
	"syscall/js"
)

func Cgodemo() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		v := 42
		C.printint(C.int(v))
		return nil
	})
}
