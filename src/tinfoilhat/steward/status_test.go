/**
 * @file status_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief test work with status table
 */

package steward_test

import (
	"log"
	"testing"
)

import "tinfoilhat/steward"

func TestPutStatus(t *testing.T) {

	db, err := openDB()

	defer db.Close()

	status := steward.Status{10, 10, 10, 10}

	err = steward.PutStatus(db.db, status)
	if err != nil {
		log.Fatalln("Add status failed:", err)
	}
}

func TestGetAllStatus(t *testing.T) {

	db, err := openDB()

	defer db.Close()

	round := 1
	team := 2
	service := 3

	steward.PutStatus(db.db, steward.Status{round, team, service, 4})
	steward.PutStatus(db.db, steward.Status{round, team, service, 10})
	steward.PutStatus(db.db, steward.Status{round, team, service, 20})

	states, err := steward.GetStates(db.db, round, team, service)
	if err != nil {
		log.Fatalln("Get states failed:", err)
	}

	if len(states) != 3 {
		log.Fatalln("Get states moar/less than put:", err)
	}

	if states[0] != 4 || states[1] != 10 || states[2] != 20 {
		log.Fatalln("Get states invalid")
	}
}
