/**
 * @file round.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date September, 2015
 * @brief queries for round table
 */

package steward

import (
	"database/sql"
	"time"
)

// Round contains info about round
type Round struct {
	ID        int
	Len       time.Duration
	StartTime time.Time
}

func createRoundTable(db *sql.DB) (err error) {

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS "round" (
		id	SERIAL PRIMARY KEY,
		len_seconds INTEGER  NOT NULL,
		start_time	TIMESTAMP with time zone DEFAULT now()
	)`)

	return
}

// NewRound add new round to database
func NewRound(db *sql.DB, len time.Duration) (round int, err error) {

	stmt, err := db.Prepare("INSERT INTO round (len_seconds) " +
		"VALUES ($1) RETURNING id")
	if err != nil {
		return
	}

	defer stmt.Close()

	err = stmt.QueryRow(len / time.Second).Scan(&round)
	if err != nil {
		return
	}

	return
}

// CurrentRound returns current round
func CurrentRound(db *sql.DB) (round Round, err error) {

	stmt, err := db.Prepare("SELECT id, len_seconds, start_time " +
		"FROM round WHERE ID = (SELECT MAX(ID) FROM round)")
	if err != nil {
		return
	}

	defer stmt.Close()

	var lenSeconds int64

	err = stmt.QueryRow().Scan(&round.ID, &lenSeconds, &round.StartTime)
	if err != nil {
		return
	}

	round.Len = time.Duration(lenSeconds) * time.Second

	return
}
