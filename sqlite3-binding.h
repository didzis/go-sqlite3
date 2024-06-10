#ifndef USE_LIBSQLITE3
#include "sqlite3/sqlite3.h"
#else // USE_LIBSQLITE3
// If users really want to link against the system sqlite3 we
// need to make this file a noop.
#endif