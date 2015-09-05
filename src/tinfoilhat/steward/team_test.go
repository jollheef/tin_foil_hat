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

import "tinfoilhat/steward"

func TestAddTeam(t *testing.T) {

	db, err := openDB()

	defer db.Close()

	_, err = steward.AddTeam(db.db, "MySUperTeam", "192.168.111/24")
	if err != nil {
		log.Fatalln("Add team failed:", err)
	}
}

func TestGetTeams(t *testing.T) {

	db, err := openDB()

	defer db.Close()

	team1 := steward.Team{-1, "MySuperTeam", "192.168.111/24"}
	team2 := steward.Team{-1, "MyFooTeam", "192.168.112/24"}

	team1.Id, _ = steward.AddTeam(db.db, team1.Name, team1.Subnet)
	team2.Id, _ = steward.AddTeam(db.db, team2.Name, team2.Subnet)

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
