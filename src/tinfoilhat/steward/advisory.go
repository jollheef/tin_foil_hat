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
		id	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
		team_id	INTEGER NOT NULL,
		score	INTEGER,
		reviewed	INTEGER,
		timestamp	INTEGER DEFAULT CURRENT_TIMESTAMP,
		text	TEXT NOT NULL
	)`)

	return
}

func AddAdvisory(db *sql.DB, team_id int, text string) (id int, err error) {

	stmt, err := db.Prepare("INSERT INTO `advisory` (`team_id`, `text`) " +
		"VALUES (?, ?)")
	if err != nil {
		return
	}

	defer stmt.Close()

	res, err := stmt.Exec(team_id, text)
	if err != nil {
		return
	}

	id64, err := res.LastInsertId()

	if err != nil {
		return
	}

	id = int(id64)

	return
}

func ReviewAdvisory(db *sql.DB, advisory_id int, score int) error {

	stmt, err := db.Prepare(
		"UPDATE `advisory` SET `score`=?, `reviewed`=? WHERE `id`=?")
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
		"SELECT sum(score) FROM `advisory` WHERE `team_id`=? " +
			"AND `reviewed`=?")
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
		"SELECT `id`, `text`, `score`, strftime('%s', `timestamp`) " +
			"FROM `advisory` ")
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var adv Advisory
		var timestamp int64

		err = rows.Scan(&adv.Id, &adv.Text, &adv.Score, &timestamp)
		if err != nil {
			return
		}

		adv.Timestamp = time.Unix(timestamp, 0)

		advisories = append(advisories, adv)
	}

	return
}
