/**
 * @file flag_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date September, 2015
 * @brief test work with flag table
 */

package steward_test

import (
	"log"
	"testing"
)

import "github.com/jollheef/tin_foil_hat/steward"

func TestAddFlag(t *testing.T) {

	db, err := openDB()

	defer db.Close()

	err = steward.AddFlag(db.db, steward.Flag{ID: 1, Flag: "lolka",
		Round: 1, TeamID: 2, ServiceID: 3, Cred: "1:2"})
	if err != nil {
		log.Fatalln("Add flag failed:", err)
	}
}

func TestFlagExist(t *testing.T) {

	db, err := openDB()

	defer db.Close()

	flg := steward.Flag{ID: 0, Flag: "tralala", Round: 5, TeamID: 10,
		ServieID: 4, Cred: "1:2"}

	err = steward.AddFlag(db.db, flg)

	exist, err := steward.FlagExist(db.db, flg.Flag)
	if !exist {
		log.Fatalln("Exist flag does not exist:", err)
	}

	exist, err = steward.FlagExist(db.db, "not_exist_flag")
	if exist {
		log.Fatalln("Not exist flag is exist:", err)
	}
}

func TestGetFlagInfo(t *testing.T) {

	db, err := openDB()

	defer db.Close()

	flg := steward.Flag{ID: 1, Flag: "asdfasdf", Round: 5345, TeamID: 433,
		ServiceID: 353, Cred: "1:2"}

	err = steward.AddFlag(db.db, flg)

	new_flg, err := steward.GetFlagInfo(db.db, flg.Flag)
	if err != nil {
		log.Fatalln("Cannot get flag info:", err)
	}

	if new_flg != flg {
		log.Fatalln("Readed flag is not equal to writed before")
	}
}

func TestGetCred(t *testing.T) {

	db, err := openDB()

	defer db.Close()

	flg := steward.Flag{ID: 1, Flag: "asdfasdf", Round: 5345, TeamID: 433,
		ServiceID: 353, Cred: "1:2"}

	err = steward.AddFlag(db.db, flg)

	flag, cred, err := steward.GetCred(db.db, flg.Round, flg.TeamID,
		flg.ServiceID)
	if err != nil {
		log.Fatalln("Get cred failed:", err)
	}

	if flag != flg.Flag || cred != flg.Cred {
		log.Fatalln("Gotten cred invalid")
	}
}
