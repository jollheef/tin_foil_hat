/**
 * @file flag.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date September, 2015
 * @brief queries for flag table
 */

package steward

import "database/sql"

// Flag contains info about flag
type Flag struct {
	ID        int
	Flag      string
	Round     int
	TeamID    int
	ServiceID int
	Cred      string
}

func createFlagTable(db *sql.DB) (err error) {

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS "flag" (
		id	SERIAL PRIMARY KEY,
		round	INTEGER NOT NULL,
		flag	TEXT NOT NULL UNIQUE,
		team_id	INTEGER NOT NULL,
		service_id	INTEGER NOT NULL,
		cred	TEXT NOT NULL
	)`)

	return
}

// AddFlag add flag to database
func AddFlag(db *sql.DB, flg Flag) error {

	stmt, err := db.Prepare("INSERT INTO flag " +
		"(round, team_id, service_id, flag, cred) " +
		"VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(flg.Round, flg.TeamID, flg.ServiceID,
		flg.Flag, flg.Cred)
	if err != nil {
		return err
	}

	return nil
}

// FlagExist check for flag exist in database
func FlagExist(db *sql.DB, flag string) (exist bool, err error) {

	stmt, err := db.Prepare(
		"SELECT EXISTS(SELECT id FROM flag WHERE flag=$1)")
	if err != nil {
		return
	}

	defer stmt.Close()

	err = stmt.QueryRow(flag).Scan(&exist)
	if err != nil {
		return
	}

	return
}

// GetFlagInfo returns info about flag
func GetFlagInfo(db *sql.DB, flag string) (flg Flag, err error) {

	flg.Flag = flag

	stmt, err := db.Prepare(
		"SELECT id, round, team_id, service_id, cred " +
			"FROM flag WHERE flag=$1")
	if err != nil {
		return
	}

	defer stmt.Close()

	err = stmt.QueryRow(flag).Scan(&flg.ID, &flg.Round, &flg.TeamID,
		&flg.ServiceID, &flg.Cred)
	if err != nil {
		return
	}

	return
}

// GetCred returns credentials for check flag in service
func GetCred(db *sql.DB, round, team, service int) (flag, cred string, err error) {

	stmt, err := db.Prepare("SELECT flag, cred FROM flag WHERE round=$1" +
		" AND team_id=$2 AND service_id=$3")
	if err != nil {
		return
	}

	defer stmt.Close()

	err = stmt.QueryRow(round, team, service).Scan(&flag, &cred)
	if err != nil {
		return
	}

	return
}
