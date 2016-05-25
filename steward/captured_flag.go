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

func CaptureFlag(db *sql.DB, flag_id, team_id int) (err error) {

	stmt, err := db.Prepare(
		"INSERT INTO captured_flag (flag_id, team_id) " +
			"VALUES ($1, $2)")
	if err != nil {
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(flag_id, team_id)
	if err != nil {
		return
	}

	return
}

func GetCapturedFlags(db *sql.DB, round, team_id int) (flgs []Flag, err error) {

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

		err = rows.Scan(&flag.Id, &flag.Flag, &flag.TeamId,
			&flag.ServiceId, &flag.Cred)
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

		err = nstmt.QueryRow(flag.Id, team_id).Scan(&captured)
		if err != nil {
			return flgs, err
		}

		if captured {
			flgs = append(flgs, flag)
		}
	}

	return
}

func AlreadyCaptured(db *sql.DB, flagId int) (captured bool, err error) {

	stmt, err := db.Prepare("SELECT EXISTS(SELECT id FROM captured_flag " +
		"WHERE flag_id=$1)")
	if err != nil {
		return
	}

	defer stmt.Close()

	err = stmt.QueryRow(flagId).Scan(&captured)
	if err != nil {
		return
	}

	return
}
