/**
 * @file flag_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief test work with flag table
 */

package steward_test

import (
	"log"
	"testing"
)

import "tinfoilhat/steward"

func TestAddFlag(t *testing.T) {

	db, err := openDB()

	defer db.Close()

	err = steward.AddFlag(db.db, steward.Flag{1, "lolka", 1, 2, 3, "1:2"})
	if err != nil {
		log.Fatalln("Add flag failed:", err)
	}
}

func TestFlagExist(t *testing.T) {

	db, err := openDB()

	defer db.Close()

	flg := steward.Flag{0, "tralala", 5, 10, 4, "1:2"}

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

	flg := steward.Flag{1, "asdfasdf", 5345, 433, 353, "1:2"}

	err = steward.AddFlag(db.db, flg)

	new_flg, err := steward.GetFlagInfo(db.db, flg.Flag)
	if err != nil {
		log.Fatalln("Cannot get flag info:", err)
	}

	if new_flg != flg {
		log.Fatalln("Readed flag is not equal to writed before")
	}
}
