package main

/*
#cgo LDFLAGS: -lX11 -lXfixes
#include <stdio.h>
#include <stdlib.h>
extern void clipboardInit();
extern void clipboardExit();
extern void clipboardWait();
extern char * clipboardGet();
extern void clipboardSet(char *data);
*/
import "C"

import (
	"fmt"
	"time"
	"unsafe"
)

func main() {
	for {
		C.clipboardWait()

		start1 := time.Now()
		buffer := C.clipboardGet()
		data := C.GoString(buffer)
		C.free(unsafe.Pointer(buffer))
		end1 := time.Since(start1)

		start2 := time.Now()
		bufferC := C.CString(data)
		C.clipboardSet(bufferC)
		C.free(unsafe.Pointer(bufferC))
		end2 := time.Since(start2)

		println("***********************************************")
		fmt.Println(end1)
		fmt.Println(end2)
		println("***********************************************")
	}
}
