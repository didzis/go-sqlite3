//go:build !libsqlite3 && carray
// +build !libsqlite3,carray

package sqlite3

/*
#cgo CFLAGS: -DSQLITE_ENABLE_CARRAY

#ifndef USE_LIBSQLITE3
#include "sqlite3-binding.h"
#else
#include <sqlite3.h>
#endif
#include <stdlib.h>
#include <stdint.h>
#include <string.h>
#include <sys/uio.h>
#include "sqlite3/carray.h"
*/
import "C"
import (
	"database/sql/driver"
	"unsafe"
)

func (c *SQLiteConn) CheckNamedValue(nv *driver.NamedValue) (err error) {
	switch nv.Value.(type) {
	case []int:
	case []int32:
	case []int64:
	case []uint32:
	case []float32:
	case []float64:
	case []string:
	case [][]byte:
	default:
		err = driver.ErrSkip
	}
	return
}

func (s *SQLiteStmt) bind_carray(value driver.Value, n C.int) (rv C.int) {
	rv = C.SQLITE_OK
	switch v := value.(type) {
	case []int:
		if unsafe.Sizeof(int(1)) == 4 {
			rv = C.sqlite3_carray_bind(s.s, n, unsafe.Pointer(&v[0]), C.int(len(v)), C.CARRAY_INT32, nil)
		} else {
			rv = C.sqlite3_carray_bind(s.s, n, unsafe.Pointer(&v[0]), C.int(len(v)), C.CARRAY_INT64, nil)
		}
	case []int32:
		rv = C.sqlite3_carray_bind(s.s, n, unsafe.Pointer(&v[0]), C.int(len(v)), C.CARRAY_INT32, nil)
	case []int64:
		rv = C.sqlite3_carray_bind(s.s, n, unsafe.Pointer(&v[0]), C.int(len(v)), C.CARRAY_INT64, nil)
	case []uint32:
		ptr := C.malloc(C.size_t(C.sizeof_int64_t * len(v)))
		next := unsafe.Pointer(uintptr(ptr))
		for _, i := range v {
			*(*C.int64_t)(next) = (C.int64_t)(i)
			next = unsafe.Pointer(uintptr(next) + uintptr(C.sizeof_int64_t))
		}
		rv = C.sqlite3_carray_bind(s.s, n, ptr, C.int(len(v)), C.CARRAY_INT64, (*[0]byte)(unsafe.Pointer(C.free)))
	case []float32:
		ptr := C.malloc(C.size_t(C.sizeof_double * len(v)))
		next := unsafe.Pointer(uintptr(ptr))
		for _, f := range v {
			*(*C.double)(next) = (C.double)(f)
			next = unsafe.Pointer(uintptr(next) + uintptr(C.sizeof_double))
		}
		rv = C.sqlite3_carray_bind(s.s, n, ptr, C.int(len(v)), C.CARRAY_DOUBLE, (*[0]byte)(unsafe.Pointer(C.free)))
	case []float64:
		rv = C.sqlite3_carray_bind(s.s, n, unsafe.Pointer(&v[0]), C.int(len(v)), C.CARRAY_DOUBLE, nil)
	case []string:
		totalSize := len(v) * C.sizeof_uintptr_t
		for _, s := range v {
			totalSize += len(s) + 1
		}

		ptr := C.malloc(C.size_t(totalSize))
		// pidx := (*[100000000]*C.char)(ptr)[:len(v)]
		next := unsafe.Pointer(uintptr(ptr) + uintptr(len(v)*C.sizeof_uintptr_t))

		for i, s := range v {
			C.memcpy(next, unsafe.Pointer(C.CString(s)), C.size_t(len(s)))

			// pidx[i] = (*C.char)(next)
			*(**C.char)(unsafe.Pointer(uintptr(ptr) + C.sizeof_uintptr_t*uintptr(i))) = (*C.char)(next)

			*(*C.char)(unsafe.Pointer(uintptr(next) + uintptr(len(s)))) = 0 // zero terminating string

			next = unsafe.Pointer(uintptr(next) + uintptr(len(s)) + 1)
		}

		rv = C.sqlite3_carray_bind(s.s, n, ptr, C.int(len(v)), C.CARRAY_TEXT, (*[0]byte)(unsafe.Pointer(C.free)))
	case [][]byte:
		ptr := C.malloc(C.size_t(len(v) * C.sizeof_uintptr_t))
		iovecs := (*[1 << 30]C.struct_iovec)(ptr)[:n:n]
		for i, b := range v {
			iovecs[i].iov_base = unsafe.Pointer(&b[0])
			iovecs[i].iov_len = C.size_t(len(b))
		}
		rv = C.sqlite3_carray_bind(s.s, n, ptr, C.int(len(v)), C.CARRAY_BLOB, (*[0]byte)(unsafe.Pointer(C.free)))
	}
	return
}
