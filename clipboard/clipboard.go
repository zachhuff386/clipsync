package clipboard

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
	"sync"
	"unsafe"
)

var (
	lock = sync.Mutex{}
)

func Wait() {
	C.clipboardWait()
}

func Get() string {
	lock.Lock()
	buffer := C.clipboardGet()
	data := C.GoString(buffer)
	C.free(unsafe.Pointer(buffer))
	lock.Unlock()
	return data
}

func Set(data string) {
	lock.Lock()
	buffer := C.CString(data)
	C.clipboardSet(buffer)
	C.free(unsafe.Pointer(buffer))
	lock.Unlock()
}
