/**
 * @file captured_flag_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief test work with captured_flag table
 */

package steward_test

import (
	"log"
	"testing"
)

import "tinfoilhat/steward"

func TestCaptureFlag(t *testing.T) {

	db, err := openDB()

	defer db.Close()

	flg := steward.CapturedFlag{steward.Flag{10, "", 0, 0, 0}, 20}

	err = steward.CaptureFlag(db.db, flg)
	if err != nil {
		log.Fatalln("Add status failed:", err)
	}
}

func TestGetCapturedFlags(t *testing.T) {

	db, err := openDB()

	defer db.Close()

	round := 1
	team_id := 1

	flg1 := steward.Flag{1, "f", round, team_id, 1}
	flg2 := steward.Flag{2, "b", round, team_id, 1}

	err = steward.AddFlag(db.db, flg1)
	if err != nil {
		log.Fatalln("Add flag failed:", err)
	}

	err = steward.AddFlag(db.db, flg2)
	if err != nil {
		log.Fatalln("Add flag failed:", err)
	}

	cflg1 := steward.CapturedFlag{flg1, 20}
	cflg2 := steward.CapturedFlag{flg2, 30}

	err = steward.CaptureFlag(db.db, cflg1)
	err = steward.CaptureFlag(db.db, cflg2)

	flags, err := steward.GetCapturedFlags(db.db, round, team_id)
	if err != nil {
		log.Fatalln("Get captured flags failed:", err)
	}

	if len(flags) != 2 {
		log.Fatalln("Get captured flags more/less than added")
	}

	if flags[0] != cflg1 || flags[1] != cflg2 {
		log.Fatalln("Getted flags invalid", flags[0], cflg1, flags[1], cflg2)
	}
}

func TestAlreadyCaptured(t *testing.T) {

	db, err := openDB()

	defer db.Close()

	flg1 := steward.Flag{1, "f", 1, 1, 1}
	flg2 := steward.Flag{2, "b", 1, 1, 1}

	cflg1 := steward.CapturedFlag{flg1, 20}
	cflg2 := steward.CapturedFlag{flg2, 30}

	err = steward.CaptureFlag(db.db, cflg1)

	captured, err := steward.AlreadyCaptured(db.db, cflg1.Flag.Id)
	if err != nil {
		log.Fatalln("Already captured check failed:", err)
	}

	if !captured {
		log.Fatalln("Captured flag is not captured")
	}

	captured, err = steward.AlreadyCaptured(db.db, cflg2.Flag.Id)
	if err != nil {
		log.Fatalln("Already captured check failed:", err)
	}

	if captured {
		log.Fatalln("Not captured flag is captured")
	}
}
