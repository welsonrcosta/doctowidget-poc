package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"doctogadget/doctowidget"
	"fmt"
	"unsafe"
)

var dw doctowidget.Doctowidget

func main() {
	f := func(s string) {
		fmt.Println(s)
	}
	dw = doctowidget.NewDoctowidget(&f, true)
	defer dw.Destroy()
	dw.Run()
}

type Callback func(*C.char)

//export startQtGadgets
func startQtGadgets(doctoWidgetCallback Callback,
	initPromptCallback Callback,
	workingCallback Callback,
	retryPromptCallback Callback, lang C.int) {
	f := func(s string) {
		cstr := C.CString(s)
		defer C.free(unsafe.Pointer(cstr))
		doctoWidgetCallback(cstr)
	}
	dw = doctowidget.NewDoctowidget(&f, false)

	dw.Run()
}
