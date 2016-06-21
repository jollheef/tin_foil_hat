/**
 * @file round_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date September, 2015
 * @brief test work with round table
 */

package steward_test

import (
	"log"
	"testing"
	"time"
)

import "github.com/jollheef/tin_foil_hat/steward"

func TestRoundWork(t *testing.T) {

	db, err := openDB()

	defer db.Close()

	round_len := time.Minute * 2

	_, err = steward.CurrentRound(db.db)
	if err == nil {
		log.Fatalln("Current round in empty database already exist")
	}

	var i int
	for i = 1; i < 5; i++ {
		new_round, err := steward.NewRound(db.db, round_len)
		if err != nil {
			log.Fatalln("Start new round fail:", err)
		}
		if new_round != i {
			log.Fatalln("New round number invalid", new_round, i)
		}

		current_round, err := steward.CurrentRound(db.db)
		if err != nil {
			log.Fatalln("Get current round fail:", err)
		}
		if current_round.ID != new_round {
			log.Fatalln("Current round number invalid")
		}
		if current_round.Len != round_len {
			log.Fatalln("Current round len invalid")
		}
		if time.Now().Sub(current_round.StartTime) > 5*time.Second {
			log.Fatalln("Time must be ~ current:",
				current_round.StartTime)
		}
	}
}
