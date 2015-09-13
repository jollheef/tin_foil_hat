/**
 * @file team.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief queries for team table
 */

package steward

import "database/sql"

type Team struct {
	Id     int
	Name   string
	Subnet string
}

func createTeamTable(db *sql.DB) (err error) {

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS team (
		id	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
		name	TEXT NOT NULL UNIQUE,
		subnet	TEXT NOT NULL UNIQUE
	)`)

	return
}

func AddTeam(db *sql.DB, name string, subnet string) (id int, err error) {

	stmt, err := db.Prepare("INSERT INTO `team` (`name`, `subnet`) " +
		"VALUES (?, ?)")
	if err != nil {
		return
	}

	defer stmt.Close()

	res, err := stmt.Exec(name, subnet)
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

func GetTeams(db *sql.DB) (teams []Team, err error) {

	rows, err := db.Query("SELECT `id`, `name`, `subnet` FROM `team`")
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var team Team

		err = rows.Scan(&team.Id, &team.Name, &team.Subnet)
		if err != nil {
			return
		}

		teams = append(teams, team)
	}

	return
}

func GetTeam(db *sql.DB, team_id int) (team Team, err error) {

	stmt, err := db.Prepare("SELECT `name`, `subnet` FROM `team`" +
		"WHERE `id`=?")
	if err != nil {
		return
	}

	defer stmt.Close()

	team.Id = team_id

	err = stmt.QueryRow(team_id).Scan(&team.Name, &team.Subnet)
	if err != nil {
		return
	}

	return
}
