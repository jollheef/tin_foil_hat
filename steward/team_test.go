/**
 * @file team_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief test work with team table
 */

package steward_test

import (
	"log"
	"testing"
)

import "github.com/jollheef/tin_foil_hat/steward"

func TestAddTeam(t *testing.T) {

	db, err := openDB()

	defer db.Close()

	team := steward.Team{-1, "MySuperTeam", "192.168.111/24", "pl.hold1"}

	_, err = steward.AddTeam(db.db, team)
	if err != nil {
		log.Fatalln("Add team failed:", err)
	}
}

func TestGetTeams(t *testing.T) {

	db, err := openDB()

	defer db.Close()

	team1 := steward.Team{-1, "MySuperTeam", "192.168.111/24", "pl.hold1"}
	team2 := steward.Team{-1, "MyFooTeam", "192.168.112/24", "pl.hold2"}

	team1.Id, _ = steward.AddTeam(db.db, team1)
	team2.Id, _ = steward.AddTeam(db.db, team2)

	teams, err := steward.GetTeams(db.db)
	if err != nil {
		log.Fatalln("Get teams failed:", err)
	}

	if len(teams) != 2 {
		log.Fatalln("Get teams more than added")
	}

	if teams[0] != team1 || teams[1] != team2 {
		log.Fatalln("Added teams broken")
	}
}

func TestGetTeam(t *testing.T) {

	db, err := openDB()

	defer db.Close()

	team1 := steward.Team{-1, "MySuperTeam", "192.168.111/24", "pl.hold"}

	team1.Id, _ = steward.AddTeam(db.db, team1)

	_team1, err := steward.GetTeam(db.db, team1.Id)
	if err != nil {
		log.Fatalln("Get team failed:", err)
	}

	if _team1 != team1 {
		log.Fatalln("Added team broken")
	}

	_, err = steward.GetTeam(db.db, 10) // invalid team id
	if err == nil {
		log.Fatalln("Get invalid team broken")
	}
}
