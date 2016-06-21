/**
 * @file steward.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
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

	err = createRoundResultTable(db)
	if err != nil {
		return err
	}

	return nil
}

// defer db.Close() after open
func OpenDatabase(path string) (db *sql.DB, err error) {

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

// CleanDatabase remove all data from database and restart sequences
func CleanDatabase(db *sql.DB) (err error) {

	tables := []string{"team", "advisory", "captured_flag", "flag",
		"service", "status", "round", "round_result"}

	for _, table := range tables {

		_, err = db.Exec("DELETE FROM " + table)
		if err != nil {
			return
		}

		_, err = db.Exec("ALTER SEQUENCE " + table + "_id_seq RESTART WITH 1;")
		if err != nil {
			return
		}
	}

	return
}
