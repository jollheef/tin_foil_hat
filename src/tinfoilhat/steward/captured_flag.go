/**
 * @file captured_flag.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief queries for captured_flag table
 */

package steward

import (
	"database/sql"
	"time"
)

type CapturedFlag struct {
	FlagId    int
	TeamId    int
	Timestamp time.Time
}

func createCapturedFlagTable(db *sql.DB) (err error) {

	_, err = db.Exec(`
	CREATE TABLE "captured_flag" (
		id	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
		flag_id	INTEGER NOT NULL,
		team_id	INTEGER NOT NULL,
		timestamp	INTEGER DEFAULT 'CURRENT_TIMESTAMP'
	)`)

	return
}
