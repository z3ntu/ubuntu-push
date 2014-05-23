package nih

/*
#cgo pkg-config: dbus-1 libnih libnih-dbus
#include <stdlib.h>
#include <nih/alloc.h>
#include <libnih-dbus.h>

char *cuote (const char *id) {
    return nih_dbus_path (NULL, "", id, NULL);
}
*/
import "C"

import (
	"unsafe"
)

func cuote(s []byte) string {
	cs := C.CString(string(s))
	defer C.free(unsafe.Pointer(cs))

	cq := C.cuote(cs)
	defer C.nih_free(unsafe.Pointer(cq))

	return C.GoString(cq)[1:]
}
