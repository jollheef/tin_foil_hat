/**
 * @file steward.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief work with database
 *
 * Contain functions for work with database.
 */

package steward

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func createSchema(db *sql.DB) error {

	err := createFlagTable(db)
	if err != nil {
		return err
	}

	err = createAdvisoryTable(db)
	if err != nil {
		return err
	}

	err = createCapturedFlagTable(db)
	if err != nil {
		return err
	}

	err = createTeamTable(db)
	if err != nil {
		return err
	}

	err = createServiceTable(db)
	if err != nil {
		return err
	}

	err = createStatusTable(db)
	if err != nil {
		return err
	}

	err = createRoundTable(db)
	if err != nil {
		return err
	}

	return nil
}

func openDatabase(path string) (db *sql.DB, err error) {

	db, err = sql.Open("postgres", path)
	if err != nil {
		return
	}

	err = createSchema(db)
	if err != nil {
		return
	}

	return
}

// for test purpose
var PrivateOpenDatabase = openDatabase

// defer db.Close() after open
func OpenDatabase() (db *sql.DB, err error) {
	return openDatabase("user=postgres dbname=tinfoilhat sslmode=disable")
}
