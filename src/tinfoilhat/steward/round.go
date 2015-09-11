/**
 * @file round.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief queries for round table
 */

package steward

import (
	"database/sql"
	"time"
)

func createRoundTable(db *sql.DB) (err error) {

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS "round" (
		id	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
		len_seconds INTEGER  NOT NULL,
		start_timestamp	INTEGER DEFAULT CURRENT_TIMESTAMP
	)`)

	return
}

func NewRound(db *sql.DB, len time.Duration) (round int64, err error) {

	stmt, err := db.Prepare("INSERT INTO `round` (len_seconds) " +
		"VALUES (?)")
	if err != nil {
		return
	}

	defer stmt.Close()

	res, err := stmt.Exec(len / time.Second)
	if err != nil {
		return
	}

	round, err = res.LastInsertId()

	if err != nil {
		return
	}

	return
}

func CurrentRound(db *sql.DB) (round int64, len time.Duration, err error) {

	stmt, err := db.Prepare("SELECT `id`, `len_seconds` FROM `round` " +
		"WHERE ID = (SELECT MAX(ID) FROM `round`)")
	if err != nil {
		return
	}

	defer stmt.Close()

	var len_seconds int64
	err = stmt.QueryRow().Scan(&round, &len_seconds)
	if err != nil {
		return
	}

	len = time.Duration(len_seconds) * time.Second

	return
}
