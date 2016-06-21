/**
 * @file team.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date September, 2015
 * @brief queries for team table
 */

package steward

import "database/sql"

// Team contains info about team
type Team struct {
	ID        int
	Name      string
	Subnet    string
	Vulnbox   string
	UseNetbox bool
	Netbox    string
}

func createTeamTable(db *sql.DB) (err error) {

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS team (
		id		SERIAL PRIMARY KEY,
		name		TEXT NOT NULL UNIQUE,
		subnet		TEXT NOT NULL UNIQUE,
		vulnbox		TEXT NOT NULL UNIQUE,
		use_netbox	BOOLEAN NOT NULL,
                netbox		TEXT NOT NULL
	)`)
	return
}

// AddTeam add team to database
func AddTeam(db *sql.DB, team Team) (id int, err error) {

	stmt, err := db.Prepare("INSERT INTO team (name, subnet, vulnbox, " +
		"use_netbox, netbox) " +
		"VALUES ($1, $2, $3, $4, $5) RETURNING id")
	if err != nil {
		return
	}

	defer stmt.Close()

	err = stmt.QueryRow(team.Name, team.Subnet, team.Vulnbox,
		team.UseNetbox, team.Netbox).Scan(&id)
	if err != nil {
		return
	}

	return
}

// GetTeams get all teams from database
func GetTeams(db *sql.DB) (teams []Team, err error) {

	rows, err := db.Query(
		"SELECT id, name, subnet, vulnbox, use_netbox, netbox FROM team")
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var team Team

		err = rows.Scan(&team.ID, &team.Name, &team.Subnet,
			&team.Vulnbox, &team.UseNetbox, &team.Netbox)
		if err != nil {
			return
		}

		teams = append(teams, team)
	}

	return
}

// GetTeam get team by id from database
func GetTeam(db *sql.DB, teamID int) (team Team, err error) {

	stmt, err := db.Prepare(
		"SELECT name, subnet, vulnbox, use_netbox, netbox FROM team " +
			"WHERE id=$1")
	if err != nil {
		return
	}

	defer stmt.Close()

	team.ID = teamID

	err = stmt.QueryRow(teamID).Scan(&team.Name, &team.Subnet,
		&team.Vulnbox, &team.UseNetbox, &team.Netbox)
	if err != nil {
		return
	}

	return
}
