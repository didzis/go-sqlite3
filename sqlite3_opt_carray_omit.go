//go:build libsqlite3 || !carray
// +build libsqlite3 !carray

package sqlite3

/*
#ifndef USE_LIBSQLITE3
#include "sqlite3-binding.h"
#else
#include <sqlite3.h>
#endif
*/
import "C"
import "database/sql/driver"

func (s *SQLiteStmt) bind_carray(value driver.Value, n C.int) (rv C.int) {
	return C.SQLITE_OK
}
