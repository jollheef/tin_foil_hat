/**
 * @file round_result.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief queries for round result table
 */

package steward

import "database/sql"

type RoundResult struct {
	Id           int
	TeamId       int
	Round        int
	AttackScore  int
	DefenceScore int
}

func createRoundResultTable(db *sql.DB) (err error) {

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS "round_result" (
		id	SERIAL PRIMARY KEY,
		team_id	INTEGER NOT NULL,
		round	INTEGER,
		attack_score	INTEGER,
		defence_score	INTEGER
	)`)

	return
}

func AddRoundResult(db *sql.DB, res RoundResult) (id int, err error) {

	if res.Round > 1 { // if not first round
		previous, err := GetRoundResult(db, res.TeamId, res.Round-1)
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

	err = stmt.QueryRow(res.TeamId, res.Round, res.AttackScore,
		res.DefenceScore).Scan(&id)
	if err != nil {
		return
	}

	return
}

func GetRoundResult(db *sql.DB, team_id, round int) (res RoundResult, err error) {

	stmt, err := db.Prepare("SELECT id, attack_score, defence_score " +
		"FROM round_result WHERE team_id=$1 AND round=$2")
	if err != nil {
		return
	}

	defer stmt.Close()

	err = stmt.QueryRow(team_id, round).Scan(&res.Id, &res.AttackScore,
		&res.DefenceScore)
	if err != nil {
		return
	}

	res.TeamId = team_id
	res.Round = round

	return
}
