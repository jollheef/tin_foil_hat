/**
 * @file status.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief queries for status table
 */

package steward

import "database/sql"

type ServiceState int

const (
	// Service is online, serves the requests, stores and
	// returns flags and behaves as expected
	STATUS_OK ServiceState = iota
	// Service is online, but behaves not as expected, e.g. if HTTP server
	// listens the port, but doesn't respond on request
	STATUS_MUMBLE
	// Service is online, but past flags cannot be retrieved
	STATUS_CORRUPT
	// Service is offline
	STATUS_DOWN
	// Checker error
	STATUS_ERROR
	// Unknown
	STATUS_UNKNOWN
)

type Status struct {
	Round     int
	TeamId    int
	ServiceId int
	State     ServiceState
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

func GetStates(db *sql.DB, halfStatus Status) (states []ServiceState,
	err error) {

	stmt, err := db.Prepare(
		"SELECT `state` FROM `status` WHERE `round`=? AND `team_id`=? " +
			"AND `service_id`=?")
	if err != nil {
		return
	}

	defer stmt.Close()

	rows, err := stmt.Query(halfStatus.Round, halfStatus.TeamId,
		halfStatus.ServiceId)
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

		states = append(states, ServiceState(state))
	}

	return
}

func GetState(db *sql.DB, halfStatus Status) (state ServiceState, err error) {

	stmt, err := db.Prepare(
		"SELECT `state` FROM `status` WHERE `round`=? AND `team_id`=? " +
			"AND `service_id`=? " +
			"AND ID = (SELECT MAX(ID) FROM `status`)")
	if err != nil {
		return
	}

	defer stmt.Close()

	err = stmt.QueryRow(halfStatus.Round, halfStatus.TeamId,
		halfStatus.ServiceId).Scan(&state)
	if err != nil {
		return
	}

	return
}
