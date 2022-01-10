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
import "unsafe"

func Wait() {
	C.clipboardWait()
}

func Get() string {
	buffer := C.clipboardGet()
	data := C.GoString(buffer)
	C.free(unsafe.Pointer(buffer))
	return data
}

func Set(data string) {
	buffer := C.CString(data)
	C.clipboardSet(buffer)
	C.free(unsafe.Pointer(buffer))
}
