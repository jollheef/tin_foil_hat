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

type Round struct {
	Id        int64
	Len       time.Duration
	StartTime time.Time
}

func createRoundTable(db *sql.DB) (err error) {

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS "round" (
		id	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
		len_seconds INTEGER  NOT NULL,
		start_time	INTEGER DEFAULT CURRENT_TIMESTAMP
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

func CurrentRound(db *sql.DB) (round Round, err error) {

	stmt, err := db.Prepare(
		"SELECT `id`, `len_seconds`, strftime('%s', `start_time`) " +
			"FROM `round` WHERE ID = (SELECT MAX(ID) FROM `round`)")
	if err != nil {
		return
	}

	defer stmt.Close()

	var len_seconds int64
	var timestamp int64

	err = stmt.QueryRow().Scan(&round.Id, &len_seconds, &timestamp)
	if err != nil {
		return
	}

	round.Len = time.Duration(len_seconds) * time.Second

	round.StartTime = time.Unix(timestamp, 0)

	return
}
