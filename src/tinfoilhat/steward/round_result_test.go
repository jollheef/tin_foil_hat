/**
 * @file round_result_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief tewt work with round result table
 */

package steward_test

import (
	"log"
	"testing"
)

import "tinfoilhat/steward"

func TestAddRoundResult(t *testing.T) {

	db, err := openDB()
	if err != nil {
		log.Fatalln("Open database failed:", err)
	}

	defer db.Close()

	first := steward.RoundResult{-1, 10, 1, 30, 40}
	second := steward.RoundResult{-1, first.TeamId, first.Round + 1, 130, 140}

	_, err = steward.AddRoundResult(db.db, first)
	if err != nil {
		log.Fatalln("Add round result failed:", err)
	}

	_, err = steward.AddRoundResult(db.db, second)
	if err != nil {
		log.Fatalln("Add round result failed:", err)
	}
}

func TestGetRoundResult(t *testing.T) {

	db, err := openDB()
	if err != nil {
		log.Fatalln("Open database failed:", err)
	}

	defer db.Close()

	first := steward.RoundResult{-1, 10, 1, 30, 40}
	second := steward.RoundResult{-1, first.TeamId, first.Round + 1, 130, 140}

	_, err = steward.AddRoundResult(db.db, first)
	if err != nil {
		log.Fatalln("Add round result failed:", err)
	}

	_, err = steward.AddRoundResult(db.db, second)
	if err != nil {
		log.Fatalln("Add round result failed:", err)
	}

	last_res, err := steward.GetLastResult(db.db, second.TeamId)
	if err != nil {
		log.Fatalln("Get last result failed:", err)
	}

	res, err := steward.GetRoundResult(db.db, second.TeamId, second.Round)
	if err != nil {
		log.Fatalln("Get round result failed:", err)
	}

	if res != last_res {
		log.Fatalln("Last result != round result", res, last_res)
	}

	attack_sum := first.AttackScore + second.AttackScore
	defence_sum := first.DefenceScore + second.DefenceScore

	if res.AttackScore != attack_sum {
		log.Fatalln("Invalid attack score value", res.AttackScore,
			attack_sum)
	}

	if res.DefenceScore != defence_sum {
		log.Fatalln("Invalid defence score value", res.DefenceScore,
			defence_sum)
	}
}
