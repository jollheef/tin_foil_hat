/**
 * @file status.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date September, 2015
 * @brief queries for status table
 */

package steward

import "database/sql"

// ServiceState provide type for service status
type ServiceState int

const (
	// StatusUP Service is online, serves the requests, stores and
	// returns flags and behaves as expected
	StatusUP ServiceState = iota
	// StatusMumble Service is online, but behaves not as expected,
	// e.g. if HTTP server listens the port, but doesn't respond on request
	StatusMumble
	// StatusCorrupt Service is online, but past flags cannot be retrieved
	StatusCorrupt
	// StatusDown Service is offline
	StatusDown
	// StatusError Checker error
	StatusError
	// StatusUnknown Unknown
	StatusUnknown
)

func (state ServiceState) String() string {
	switch state {
	case StatusUP:
		return "up"
	case StatusMumble:
		return "mumble"
	case StatusCorrupt:
		return "corrupt"
	case StatusDown:
		return "down"
	case StatusError:
		return "error"
	case StatusUnknown:
		return "unknown"
	}

	return "undefined"
}

// Status contains info about services status
type Status struct {
	Round     int
	TeamID    int
	ServiceID int
	State     ServiceState
}

func createStatusTable(db *sql.DB) (err error) {

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS "status" (
		id	SERIAL PRIMARY KEY,
		round	INTEGER NOT NULL,
		team_id	INTEGER NOT NULL,
		service_id	INTEGER NOT NULL,
		state	INTEGER NOT NULL,
		timestamp	TIMESTAMP with time zone DEFAULT now()
	)`)

	return
}

// PutStatus add status to database
func PutStatus(db *sql.DB, status Status) (err error) {

	stmt, err := db.Prepare("INSERT INTO status (round, team_id, " +
		"service_id, state) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(status.Round, status.TeamID, status.ServiceID,
		status.State)
	if err != nil {
		return
	}

	return
}

// GetStates get states for services status
func GetStates(db *sql.DB, halfStatus Status) (states []ServiceState,
	err error) {

	stmt, err := db.Prepare(
		"SELECT state FROM status WHERE round=$1 AND team_id=$2 " +
			"AND service_id=$3")
	if err != nil {
		return
	}

	defer stmt.Close()

	rows, err := stmt.Query(halfStatus.Round, halfStatus.TeamID,
		halfStatus.ServiceID)
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

// GetState get state for service status
func GetState(db *sql.DB, halfStatus Status) (state ServiceState, err error) {

	stmt, err := db.Prepare(
		"SELECT state FROM status WHERE round=$1 AND team_id=$2 " +
			"AND service_id=$3 " +
			"AND ID = (SELECT MAX(ID) FROM status " +
			"WHERE round=$1 AND team_id=$2 AND service_id=$3)")
	if err != nil {
		return
	}

	defer stmt.Close()

	err = stmt.QueryRow(halfStatus.Round, halfStatus.TeamID,
		halfStatus.ServiceID).Scan(&state)
	if err != nil {
		return
	}

	return
}
