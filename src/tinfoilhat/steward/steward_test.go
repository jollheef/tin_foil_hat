/**
 * @file steward_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief test general work with database functions
 */

package steward_test

import (
	"database/sql"
	"log"
	"os"
	"testing"
)

import "tinfoilhat/steward"

type testDB struct {
	db *sql.DB
}

const db_path string = "/tmp/test.sql"

func openDB() (t testDB, err error) {

	os.Remove(db_path)

	t.db, err = steward.PrivateOpenDatabase(db_path)

	return
}

func (t testDB) Close() {

	t.db.Close()

	os.Remove(db_path)
}

func TestOpenDatabase(t *testing.T) {

	db, err := openDB()
	if err != nil {
		log.Fatalln("Database open failed:", err)
	}

	defer db.Close()
}
