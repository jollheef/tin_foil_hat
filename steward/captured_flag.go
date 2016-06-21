/**
 * @file captured_flag.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date September, 2015
 * @brief queries for captured_flag table
 */

package steward

import "database/sql"

func createCapturedFlagTable(db *sql.DB) (err error) {

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS "captured_flag" (
		id	SERIAL PRIMARY KEY,
		flag_id	INTEGER NOT NULL,
		team_id	INTEGER NOT NULL,
		timestamp	TIMESTAMP with time zone DEFAULT now()
	)`)

	return
}

// CaptureFlag add correct flag to db
func CaptureFlag(db *sql.DB, flagID, teamID int) (err error) {

	stmt, err := db.Prepare(
		"INSERT INTO captured_flag (flag_id, team_id) " +
			"VALUES ($1, $2)")
	if err != nil {
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(flagID, teamID)
	if err != nil {
		return
	}

	return
}

// GetCapturedFlags get all captured flags for team and round
func GetCapturedFlags(db *sql.DB, round, teamID int) (flgs []Flag, err error) {

	stmt, err := db.Prepare("SELECT id, flag, team_id, " +
		"service_id, cred FROM flag WHERE round=$1")
	if err != nil {
		return
	}

	defer stmt.Close()

	rows, err := stmt.Query(round)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var flag Flag
		flag.Round = round

		err = rows.Scan(&flag.ID, &flag.Flag, &flag.TeamID,
			&flag.ServiceID, &flag.Cred)
		if err != nil {
			return
		}

		nstmt, err := db.Prepare(
			"SELECT EXISTS(SELECT id FROM captured_flag " +
				"WHERE flag_id=$1 AND team_id=$2)")
		if err != nil {
			return flgs, err
		}

		defer nstmt.Close()

		var captured bool

		err = nstmt.QueryRow(flag.ID, teamID).Scan(&captured)
		if err != nil {
			return flgs, err
		}

		if captured {
			flgs = append(flgs, flag)
		}
	}

	return
}

// AlreadyCaptured returns false if flag already captured
func AlreadyCaptured(db *sql.DB, flagID int) (captured bool, err error) {

	stmt, err := db.Prepare("SELECT EXISTS(SELECT id FROM captured_flag " +
		"WHERE flag_id=$1)")
	if err != nil {
		return
	}

	defer stmt.Close()

	err = stmt.QueryRow(flagID).Scan(&captured)
	if err != nil {
		return
	}

	return
}
