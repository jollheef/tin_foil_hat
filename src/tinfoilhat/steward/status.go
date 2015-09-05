/**
 * @file status.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief queries for status table
 */

package steward

import (
	"database/sql"
	"time"
)

type Status struct {
	Round     int
	TeamId    int
	ServiceId int
	Status    int
	Timestamp time.Time
}

func createStatusTable(db *sql.DB) (err error) {

	_, err = db.Exec(`
	CREATE TABLE "status" (
		id	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
		round	INTEGER NOT NULL,
		team_id	INTEGER NOT NULL,
		service_id	INTEGER NOT NULL,
		status	INTEGER NOT NULL,
		timestamp	INTEGER DEFAULT 'CURRENT_TIMESTAMP'
	)`)

	return
}
