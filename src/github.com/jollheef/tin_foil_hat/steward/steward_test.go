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
	"testing"
)

import "github.com/jollheef/tin_foil_hat/steward"

type testDB struct {
	db *sql.DB
}

const db_path string = "user=postgres dbname=tinfoilhat_test sslmode=disable"

func openDB() (t testDB, err error) {

	t.db, err = steward.OpenDatabase(db_path)

	return
}

func (t testDB) Close() {

	t.db.Exec("DROP SCHEMA public CASCADE")
	t.db.Exec("CREATE SCHEMA public")

	t.db.Close()
}

func TestOpenDatabase(t *testing.T) {

	db, err := openDB()
	if err != nil {
		log.Fatalln("Database open failed:", err)
	}

	defer db.Close()
}

func TestCleanDatabase(*testing.T) {

	db, err := openDB()
	if err != nil {
		log.Fatalln("Database open failed:", err)
	}

	err = steward.CleanDatabase(db.db)
	if err != nil {
		log.Fatalln("Clean database failed:", err)
	}

	defer db.Close()
}
