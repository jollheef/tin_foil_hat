/**
 * @file advisory.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief queries for advisory table
 */

package steward

import "time"

import "database/sql"

type Advisory struct {
	Id        int
	Text      string
	Score     int
	Timestamp time.Time
}

func createAdvisoryTable(db *sql.DB) (err error) {

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS "advisory" (
		id	SERIAL PRIMARY KEY,
		team_id	INTEGER NOT NULL,
		score	INTEGER,
		reviewed	BOOLEAN,
		timestamp	TIMESTAMP with time zone DEFAULT now(),
		text	TEXT NOT NULL
	)`)

	return
}

func AddAdvisory(db *sql.DB, team_id int, text string) (id int, err error) {

	stmt, err := db.Prepare("INSERT INTO advisory (team_id, text) " +
		"VALUES ($1, $2) RETURNING id")
	if err != nil {
		return
	}

	defer stmt.Close()

	err = stmt.QueryRow(team_id, text).Scan(&id)
	if err != nil {
		return
	}

	return
}

func ReviewAdvisory(db *sql.DB, advisory_id int, score int) error {

	stmt, err := db.Prepare(
		"UPDATE advisory SET score=$1, reviewed=$2 WHERE id=$3")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(score, true, advisory_id)

	if err != nil {
		return err
	}

	return nil
}

func GetAdvisoryScore(db *sql.DB, team_id int) (score int, err error) {

	stmt, err := db.Prepare(
		"SELECT sum(score) FROM advisory WHERE team_id=$1 " +
			"AND reviewed=$2")
	if err != nil {
		return
	}

	defer stmt.Close()

	err = stmt.QueryRow(team_id, true).Scan(&score)
	if err != nil {
		return
	}

	return
}

func GetAdvisories(db *sql.DB) (advisories []Advisory, err error) {

	rows, err := db.Query(
		"SELECT id, text, score, timestamp FROM advisory ")
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var adv Advisory

		err = rows.Scan(&adv.Id, &adv.Text, &adv.Score, &adv.Timestamp)
		if err != nil {
			return
		}

		advisories = append(advisories, adv)
	}

	return
}
