/**
 * @file status.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief queries for status table
 */

package steward

import "database/sql"

type Status struct {
	Round     int
	TeamId    int
	ServiceId int
	State     int
}

func createStatusTable(db *sql.DB) (err error) {

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS "status" (
		id	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
		round	INTEGER NOT NULL,
		team_id	INTEGER NOT NULL,
		service_id	INTEGER NOT NULL,
		state	INTEGER NOT NULL,
		timestamp	INTEGER DEFAULT 'CURRENT_TIMESTAMP'
	)`)

	return
}

func PutStatus(db *sql.DB, status Status) (err error) {

	stmt, err := db.Prepare("INSERT INTO `status` (`round`, `team_id`, " +
		"`service_id`, `state`) VALUES (?, ?, ?, ?)")
	if err != nil {
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(status.Round, status.TeamId, status.ServiceId,
		status.State)
	if err != nil {
		return
	}

	return
}

func GetStates(db *sql.DB, round int, teamId int, serviceId int) (states []int,
	err error) {

	stmt, err := db.Prepare(
		"SELECT `state` FROM `status` WHERE `round`=? AND `team_id`=? " +
			"AND `service_id`=?")
	if err != nil {
		return
	}

	defer stmt.Close()

	rows, err := stmt.Query(round, teamId, serviceId)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var state int

		err = rows.Scan(&state)
		if err != nil {
			return
		}

		states = append(states, state)
	}

	return
}
