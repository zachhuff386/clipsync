package clipboard

/*
#cgo LDFLAGS: -lX11 -lXfixes
#include <stdio.h>
#include <stdlib.h>
extern void clipboardWait();
extern char * clipboardGet();
extern void clipboardSet(char *data);
*/
import "C"

import (
	"sync"
	"time"
	"unsafe"
)

var (
	lock       = sync.Mutex{}
	lastChange = time.Now()
)

func Wait() {
	for {
		C.clipboardWait()
		lock.Lock()
		if time.Since(lastChange) < 200*time.Millisecond {
			lock.Unlock()
			continue
		} else {
			lastChange = time.Now()
			lock.Unlock()
			break
		}
	}
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
	lastChange = time.Now()
	lock.Unlock()
}
