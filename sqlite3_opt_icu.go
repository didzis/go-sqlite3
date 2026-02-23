// Copyright (C) 2019 Yasuhiro Matsumoto <mattn.jp@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

//go:build sqlite_icu || icu
// +build sqlite_icu icu

package sqlite3

/*
#cgo CFLAGS: -DSQLITE_ENABLE_ICU
#cgo linux LDFLAGS: -licudata -licui18n -licuio -licutest -licutu -licuuc

#cgo darwin,amd64 CFLAGS:  -I/usr/local/opt/icu4c/include
#cgo darwin,arm64 LDFLAGS: /opt/homebrew/opt/icu4c/lib/libicudata.a /opt/homebrew/opt/icu4c/lib/libicui18n.a /opt/homebrew/opt/icu4c/lib/libicuio.a /opt/homebrew/opt/icu4c/lib/libicutest.a /opt/homebrew/opt/icu4c/lib/libicutu.a /opt/homebrew/opt/icu4c/lib/libicuuc.a

#cgo darwin,arm64 CFLAGS:  -I/opt/homebrew/opt/icu4c/include
#cgo darwin,amd64 LDFLAGS: /usr/local/opt/icu4c/lib/libicudata.a /usr/local/opt/icu4c/lib/libicui18n.a /usr/local/opt/icu4c/lib/libicuio.a /usr/local/opt/icu4c/lib/libicutest.a /usr/local/opt/icu4c/lib/libicutu.a /usr/local/opt/icu4c/lib/libicuuc.a

#cgo openbsd LDFLAGS: -lsqlite3
*/
import "C"
