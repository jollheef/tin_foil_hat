/**
 * @file captured_flag_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date September, 2015
 * @brief test work with captured_flag table
 */

package steward_test

import (
	"log"
	"testing"
)

import "github.com/jollheef/tin_foil_hat/steward"

func TestCaptureFlag(t *testing.T) {

	db, err := openDB()

	defer db.Close()

	err = steward.CaptureFlag(db.db, 10, 20)
	if err != nil {
		log.Fatalln("Capture flag failed:", err)
	}
}

func TestGetCapturedFlags(t *testing.T) {

	db, err := openDB()

	defer db.Close()

	round := 1
	team_id := 1

	flg1 := steward.Flag{ID: 1, Flag: "f", Round: round, TeamID: team_id,
		ServiceID: 1, Cred: "1:2"}
	flg2 := steward.Flag{ID: 2, Flag: "b", Round: round, TeamID: team_id,
		ServiceID: 1, Cred: "1:2"}

	err = steward.AddFlag(db.db, flg1)
	if err != nil {
		log.Fatalln("Add flag failed:", err)
	}

	err = steward.AddFlag(db.db, flg2)
	if err != nil {
		log.Fatalln("Add flag failed:", err)
	}

	err = steward.CaptureFlag(db.db, flg1.ID, 20)
	err = steward.CaptureFlag(db.db, flg2.ID, 30)

	flags1, err := steward.GetCapturedFlags(db.db, round, 20)
	if err != nil {
		log.Fatalln("Get captured flags failed:", err)
	}

	if len(flags1) != 1 {
		log.Fatalln("Get captured flags more/less than added")
	}

	flags2, err := steward.GetCapturedFlags(db.db, round, 30)
	if err != nil {
		log.Fatalln("Get captured flags failed:", err)
	}

	if len(flags2) != 1 {
		log.Fatalln("Get captured flags more/less than added")
	}

	if flags1[0] != flg1 || flags2[0] != flg2 {
		log.Fatalln("Getted flags invalid", flags1[0], flg1, flags2[0], flg2)
	}
}

func TestAlreadyCaptured(t *testing.T) {

	db, err := openDB()

	defer db.Close()

	flg1 := steward.Flag{ID: 1, Flag: "f", Round: 1, TeamID: 1,
		ServiceID: 1, Cred: "1:2"}
	flg2 := steward.Flag{ID: 2, Flag: "b", Round: 1, TeamID: 1,
		ServiceID: 1, Cred: "1:2"}

	err = steward.CaptureFlag(db.db, flg1.ID, 20)

	captured, err := steward.AlreadyCaptured(db.db, flg1.ID)
	if err != nil {
		log.Fatalln("Already captured check failed:", err)
	}

	if !captured {
		log.Fatalln("Captured flag is not captured")
	}

	captured, err = steward.AlreadyCaptured(db.db, flg2.ID)
	if err != nil {
		log.Fatalln("Already captured check failed:", err)
	}

	if captured {
		log.Fatalln("Not captured flag is captured")
	}
}
