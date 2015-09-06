/**
 * @file captured_flag.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief queries for captured_flag table
 */

package steward

import "database/sql"

type CapturedFlag struct {
	Flag   Flag
	TeamId int
}

func createCapturedFlagTable(db *sql.DB) (err error) {

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS "captured_flag" (
		id	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
		flag_id	INTEGER NOT NULL,
		team_id	INTEGER NOT NULL,
		timestamp	INTEGER DEFAULT 'CURRENT_TIMESTAMP'
	)`)

	return
}

func CaptureFlag(db *sql.DB, flg CapturedFlag) (err error) {

	stmt, err := db.Prepare(
		"INSERT INTO `captured_flag` (`flag_id`, `team_id`) " +
			"VALUES (?, ?)")
	if err != nil {
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(flg.Flag.Id, flg.TeamId)
	if err != nil {
		return
	}

	return
}

func GetCapturedFlags(db *sql.DB, round int, team_id int) (cflgs []CapturedFlag,
	err error) {

	tx, err := db.Begin()
	if err != nil {
		return
	}

	defer tx.Commit()

	stmt, err := tx.Prepare("SELECT `id`, `flag`, `team_id`, " +
		"`service_id` FROM `flag` WHERE `round`=?")
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
		var cflag CapturedFlag
		cflag.Flag.Round = round

		err = rows.Scan(&cflag.Flag.Id, &cflag.Flag.Flag,
			&cflag.Flag.TeamId, &cflag.Flag.ServiceId)
		if err != nil {
			return
		}

		nstmt, err := tx.Prepare(
			"SELECT `team_id` FROM `captured_flag` " +
				"WHERE `flag_id`=?")
		if err != nil {
			return cflgs, err
		}

		defer nstmt.Close()

		nrows, err := nstmt.Query(cflag.Flag.Id)
		if err != nil {
			return cflgs, err
		}

		defer nrows.Close()

		for nrows.Next() {
			nrows.Scan(&cflag.TeamId)
			cflgs = append(cflgs, cflag)
		}
	}

	return
}

func AlreadyCaptured(db *sql.DB, flagId int) (captured bool, err error) {

	stmt, err := db.Prepare("SELECT EXISTS(SELECT `id` FROM `captured_flag` " +
		"WHERE `flag_id`=?)")
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
