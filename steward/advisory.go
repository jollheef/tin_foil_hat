/**
 * @file advisory.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date September, 2015
 * @brief queries for advisory table
 */

package steward

import "time"

import "database/sql"

// Advisory contains info about advisory
type Advisory struct {
	ID        int
	Text      string
	Reviewed  bool
	Score     int
	Timestamp time.Time
}

func createAdvisoryTable(db *sql.DB) (err error) {

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS "advisory" (
		id	SERIAL PRIMARY KEY,
		team_id	INTEGER NOT NULL,
		score	INTEGER DEFAULT 0,
		reviewed	BOOLEAN DEFAULT false,
		hided	BOOLEAN DEFAULT false,
		timestamp	TIMESTAMP with time zone DEFAULT now(),
		text	TEXT NOT NULL
	)`)

	return
}

// AddAdvisory add advisory for team to database
func AddAdvisory(db *sql.DB, teamID int, text string) (id int, err error) {

	stmt, err := db.Prepare("INSERT INTO advisory (team_id, text) " +
		"VALUES ($1, $2) RETURNING id")
	if err != nil {
		return
	}

	defer stmt.Close()

	err = stmt.QueryRow(teamID, text).Scan(&id)
	if err != nil {
		return
	}

	return
}

// ReviewAdvisory used for set score for advisory
func ReviewAdvisory(db *sql.DB, advisoryID int, score int) error {

	stmt, err := db.Prepare(
		"UPDATE advisory SET score=$1, reviewed=$2 WHERE id=$3")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(score, true, advisoryID)

	if err != nil {
		return err
	}

	return nil
}

// HideAdvisory used for hide advisory from cli
func HideAdvisory(db *sql.DB, advisoryID int, hide bool) error {

	stmt, err := db.Prepare(
		"UPDATE advisory SET hided=$1 WHERE id=$2")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(hide, advisoryID)

	if err != nil {
		return err
	}

	return nil
}

// GetAdvisoryScore get advisory score for team
func GetAdvisoryScore(db *sql.DB, teamID int) (score int, err error) {

	stmt, err := db.Prepare(
		"SELECT sum(score) FROM advisory WHERE team_id=$1 " +
			"AND reviewed=$2")
	if err != nil {
		return
	}

	defer stmt.Close()

	err = stmt.QueryRow(teamID, true).Scan(&score)
	if err != nil {
		return
	}

	return
}

// GetAdvisories get all advisories
func GetAdvisories(db *sql.DB) (advisories []Advisory, err error) {

	rows, err := db.Query(
		"SELECT id, text, score, timestamp, reviewed FROM advisory " +
			"WHERE hided=false")
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var adv Advisory

		err = rows.Scan(&adv.ID, &adv.Text, &adv.Score, &adv.Timestamp,
			&adv.Reviewed)
		if err != nil {
			return
		}

		advisories = append(advisories, adv)
	}

	return
}
