/**
 * @file flag.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief queries for flag table
 */

package steward

import "database/sql"

type Flag struct {
	Id        int
	Flag      string
	Round     int
	TeamId    int
	ServiceId int
}

func createFlagTable(db *sql.DB) (err error) {

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS "flag" (
		id	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
		round	INTEGER NOT NULL,
		flag	TEXT NOT NULL UNIQUE,
		team_id	INTEGER NOT NULL,
		service_id	INTEGER NOT NULL
	)`)

	return
}

func AddFlag(db *sql.DB, flg Flag) error {

	stmt, err := db.Prepare(
		"INSERT INTO `flag` (`round`, `team_id`, `service_id`, `flag`) " +
			"VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(flg.Round, flg.TeamId, flg.ServiceId, flg.Flag)
	if err != nil {
		return err
	}

	return nil
}

func FlagExist(db *sql.DB, flag string) (exist bool, err error) {

	stmt, err := db.Prepare(
		"SELECT EXISTS(SELECT `id` FROM `flag` WHERE `flag`=?)")
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

func GetFlagInfo(db *sql.DB, flag string) (flg Flag, err error) {

	flg.Flag = flag

	stmt, err := db.Prepare(
		"SELECT `id`, `round`, `team_id`, `service_id` " +
			"FROM `flag` WHERE `flag`=?")
	if err != nil {
		return
	}

	defer stmt.Close()

	err = stmt.QueryRow(flag).Scan(&flg.Id, &flg.Round, &flg.TeamId,
		&flg.ServiceId)
	if err != nil {
		return
	}

	return
}
