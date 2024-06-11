//go:build !libsqlite3 && carray
// +build !libsqlite3,carray

package sqlite3

import (
	"database/sql"
	"os"
	"testing"
)

func TestCArray(t *testing.T) {
	tempFilename := TempFilename(t)
	defer os.Remove(tempFilename)
	db, err := sql.Open(driverName, tempFilename)
	if err != nil {
		t.Fatal("Failed to open the test database:", err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		t.Fatal("Failed to connect to the test database:", err)
	}

	_, err = db.Exec(`CREATE TABLE foo (id integer, name string)`)
	if err != nil {
		t.Fatal("Failed to create table in test database:", err)
	}

	_, err = db.Exec(`INSERT INTO foo(id, name) VALUES(1, "first")`)
	if err != nil {
		t.Fatal("Failed to insert data into test database", err)
	}
	_, err = db.Exec(`INSERT INTO foo(id, name) VALUES(2, "second")`)
	if err != nil {
		t.Fatal("Failed to insert data into test database", err)
	}
	_, err = db.Exec(`INSERT INTO foo(id, name) VALUES(3, "third")`)
	if err != nil {
		t.Fatal("Failed to insert data into test database", err)
	}
	_, err = db.Exec(`INSERT INTO foo(id, name) VALUES(4, "fourth")`)
	if err != nil {
		t.Fatal("Failed to insert data into test database", err)
	}

	// select multiple rows using carray()
	rows, err := db.Query("SELECT id, name FROM foo WHERE id in carray(?)", []int{1, 4})
	if err != nil {
		t.Fatal("Failed to select from the test table:", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		if id != 1 && id != 4 {
			t.Fatal("Unexpected row selected")
		}
		if id == 1 && name != "first" {
			t.Fatal("Unexpected value selected")
		} else if id == 4 && name != "fourth" {
			t.Fatal("Unexpected value selected")
		}
	}

	err = rows.Err()
	if err != nil {
		t.Fatal("Failed to read rows:", err)
	}
}
