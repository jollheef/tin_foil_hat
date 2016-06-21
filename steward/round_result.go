/**
 * @file round_result.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date September, 2015
 * @brief queries for round result table
 */

package steward

import (
	"database/sql"
	"sync"
)

// RoundResult contains info about result of round
type RoundResult struct {
	ID           int
	TeamID       int
	Round        int
	AttackScore  float64
	DefenceScore float64
}

func createRoundResultTable(db *sql.DB) (err error) {

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS "round_result" (
		id	SERIAL PRIMARY KEY,
		team_id	INTEGER NOT NULL,
		round	INTEGER,
		attack_score	FLOAT(24),
		defence_score	FLOAT(24),
		UNIQUE (team_id, round)
	);`)

	return
}

var addRoundResultMutex sync.Mutex // Use as FIFO queue

// AddRoundResult add round result to database
func AddRoundResult(db *sql.DB, res RoundResult) (id int, err error) {

	addRoundResultMutex.Lock()

	defer addRoundResultMutex.Unlock()

	if res.Round > 1 { // if not first round
		previous, err := GetLastResult(db, res.TeamID)
		if err != nil {
			return id, err
		}
		res.AttackScore += previous.AttackScore
		res.DefenceScore += previous.DefenceScore
	}

	stmt, err := db.Prepare("INSERT INTO round_result " +
		"(team_id, round, attack_score, defence_score) " +
		"VALUES ($1, $2, $3, $4) RETURNING id")
	if err != nil {
		return
	}

	defer stmt.Close()

	err = stmt.QueryRow(res.TeamID, res.Round, res.AttackScore,
		res.DefenceScore).Scan(&id)
	if err != nil {
		return
	}

	return
}

// GetRoundResult get result for team and round
func GetRoundResult(db *sql.DB, teamID, round int) (res RoundResult, err error) {

	stmt, err := db.Prepare("SELECT id, attack_score, defence_score " +
		"FROM round_result WHERE team_id=$1 AND round=$2")
	if err != nil {
		return
	}

	defer stmt.Close()

	err = stmt.QueryRow(teamID, round).Scan(&res.ID, &res.AttackScore,
		&res.DefenceScore)
	if err != nil {
		return
	}

	res.TeamID = teamID
	res.Round = round

	return
}

// GetLastResult get last round result for team
func GetLastResult(db *sql.DB, teamID int) (res RoundResult, err error) {

	stmt, err := db.Prepare("SELECT id, round, attack_score, defence_score " +
		"FROM round_result WHERE team_id=$1 " +
		"AND round = (SELECT MAX(round) FROM round_result " +
		"WHERE team_id=$1)")
	if err != nil {
		return
	}

	defer stmt.Close()

	err = stmt.QueryRow(teamID).Scan(&res.ID, &res.Round, &res.AttackScore,
		&res.DefenceScore)
	if err != nil {
		return
	}

	res.TeamID = teamID

	return

}
