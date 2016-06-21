/**
 * @file counter_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date September, 2015
 * @brief test counter package
 */

package counter_test

import (
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"
)

import (
	"github.com/jollheef/tin_foil_hat/counter"
	"github.com/jollheef/tin_foil_hat/steward"
	"github.com/jollheef/tin_foil_hat/vexillary"
)

type testDB struct {
	db *sql.DB
}

const db_path string = "user=postgres dbname=tinfoilhat_test sslmode=disable"

func openDB() (t testDB, err error) {

	t.db, err = steward.OpenDatabase(db_path)

	t.Close()

	t.db, err = steward.OpenDatabase(db_path)

	return
}

func (t testDB) Close() {

	t.db.Exec("DROP TABLE team")
	t.db.Exec("DROP TABLE advisory")
	t.db.Exec("DROP TABLE captured_flag")
	t.db.Exec("DROP TABLE flag")
	t.db.Exec("DROP TABLE service")
	t.db.Exec("DROP TABLE status")
	t.db.Exec("DROP TABLE round")
	t.db.Exec("DROP TABLE round_result")

	t.db.Close()
}

func TestCountStatesResult(*testing.T) {

	db, err := openDB()
	if err != nil {
		log.Fatalln("Open database failed:", err)
	}

	defer db.Close()

	r := 1 // round
	t := 1 // team id
	s := 1 // service id

	svc := steward.Service{s, "foo", 8080, "", false}

	steward.PutStatus(db.db, steward.Status{r, t, s, steward.StatusUP})
	steward.PutStatus(db.db, steward.Status{r, t, s, steward.StatusUP})
	steward.PutStatus(db.db, steward.Status{r, t, s, steward.StatusUP})
	steward.PutStatus(db.db, steward.Status{r, t, s, steward.StatusMumble})

	res, err := counter.CountStatesResult(db.db, r, t, svc)
	if err != nil {
		log.Fatalln("Count states failed:", err)
	}

	must_be := 0.75

	if res != must_be {
		log.Fatalln("Result invalid:", res, "instead", must_be)
	}

}

func TestCountDefenceResult(*testing.T) {

	db, err := openDB()
	if err != nil {
		log.Fatalln("Open database failed:", err)
	}

	defer db.Close()

	r := 1
	t := 1

	services := make([]steward.Service, 0)
	services = append(services, steward.Service{1, "foo", 8080, "", false})
	services = append(services, steward.Service{2, "bar", 8081, "", false})
	services = append(services, steward.Service{3, "baz", 8082, "", false})
	services = append(services, steward.Service{4, "qwe", 8083, "", false})

	steward.PutStatus(db.db, steward.Status{r, t, 1, steward.StatusUP})
	steward.PutStatus(db.db, steward.Status{r, t, 1, steward.StatusUP})
	steward.PutStatus(db.db, steward.Status{r, t, 1, steward.StatusUP})
	steward.PutStatus(db.db, steward.Status{r, t, 1, steward.StatusDown})

	steward.PutStatus(db.db, steward.Status{r, t, 2, steward.StatusDown})
	steward.PutStatus(db.db, steward.Status{r, t, 2, steward.StatusDown})
	steward.PutStatus(db.db, steward.Status{r, t, 2, steward.StatusDown})
	steward.PutStatus(db.db, steward.Status{r, t, 2, steward.StatusDown})

	steward.PutStatus(db.db, steward.Status{r, t, 3, steward.StatusUP})
	steward.PutStatus(db.db, steward.Status{r, t, 3, steward.StatusUP})
	steward.PutStatus(db.db, steward.Status{r, t, 3, steward.StatusUP})

	steward.PutStatus(db.db, steward.Status{r, t, 4, steward.StatusUP})
	steward.PutStatus(db.db, steward.Status{r, t, 4, steward.StatusDown})
	steward.PutStatus(db.db, steward.Status{r, t, 4, steward.StatusDown})
	steward.PutStatus(db.db, steward.Status{r, t, 4, steward.StatusDown})

	res, err := counter.CountDefenceResult(db.db, r, t, services)
	if err != nil {
		log.Fatalln("Count defence result failed:", err)
	}

	must_be := 0.75*0.25 + 0 + 1*0.25 + 0.25*0.25

	if res != must_be {
		log.Fatalln("Result invalid:", res, "instead", must_be)

	}
}

func TestCountRound(*testing.T) {

	db, err := openDB()
	if err != nil {
		log.Fatalln("Open database failed:", err)
	}

	defer db.Close()

	for index, team := range []string{
		"FooTeam", "BarTeam", "BazTeam", "Ololosha"} {

		// just trick for bypass UNIQUE team subnet
		subnet := fmt.Sprintf("127.%d.0.1/24", index)

		vulnbox := fmt.Sprintf("127.0.%d.3/24", index)

		t := steward.Team{ID: -1, Name: team, Subnet: subnet,
			Vulnbox: vulnbox}

		_, err = steward.AddTeam(db.db, t)
		if err != nil {
			log.Fatalln("Add team failed:", err)
		}
	}

	for _, service := range []string{"Foo", "Bar", "Baz", "Boo"} {

		err = steward.AddService(db.db,
			steward.Service{-1, service, 8080, "", false})
		if err != nil {
			log.Fatalln("Add service failed:", err)
		}
	}

	priv, err := vexillary.GenerateKey()
	if err != nil {
		log.Fatalln("Generate key failed:", err)
	}

	round, err := steward.NewRound(db.db, time.Minute)
	if err != nil {
		log.Fatalln("Create new round failed:", err)
	}

	teams, err := steward.GetTeams(db.db)
	if err != nil {
		log.Fatalln("Get teams failed:", err)
	}

	services, err := steward.GetServices(db.db)
	if err != nil {
		log.Fatalln("Get services failed:", err)
	}

	flags := make([]string, 0)

	for _, team := range teams {
		for _, svc := range services {

			flag, err := vexillary.GenerateFlag(priv)
			if err != nil {
				log.Fatalln("Generate flag failed:", err)
			}

			flags = append(flags, flag)

			flg := steward.Flag{-1, flag, round, team.ID, svc.ID, ""}

			err = steward.AddFlag(db.db, flg)
			if err != nil {
				log.Fatalln("Add flag to database failed:", err)
			}

			err = steward.PutStatus(db.db, steward.Status{round,
				team.ID, svc.ID, steward.StatusUP})
			if err != nil {
				log.Fatalln("Put status to database failed:", err)
			}
		}
	}

	flag1, err := steward.GetFlagInfo(db.db, flags[2])
	if err != nil {
		log.Fatalln("Get flag info failed:", err)
	}

	err = steward.CaptureFlag(db.db, flag1.ID, teams[2].ID)
	if err != nil {
		log.Fatalln("Capture flag failed:", err)
	}

	flag2, err := steward.GetFlagInfo(db.db, flags[7])
	if err != nil {
		log.Fatalln("Get flag info failed:", err)
	}

	err = steward.CaptureFlag(db.db, flag2.ID, teams[3].ID)
	if err != nil {
		log.Fatalln("Capture flag failed:", err)
	}

	err = counter.CountRound(db.db, round, teams, services)
	if err != nil {
		log.Fatalln("Count round failed:", err)
	}

	res, err := steward.GetRoundResult(db.db, teams[0].ID, round)
	if res.AttackScore != 0.0 || res.DefenceScore != 1.75 {
		log.Fatalln("Invalid result:", res)
	}

	res, err = steward.GetRoundResult(db.db, teams[1].ID, round)
	if res.AttackScore != 0.0 || res.DefenceScore != 1.75 {
		log.Fatalln("Invalid result:", res)
	}

	res, err = steward.GetRoundResult(db.db, teams[2].ID, round)
	if res.AttackScore != 0.25 || res.DefenceScore != 2.0 {
		log.Fatalln("Invalid result:", res)
	}

	res, err = steward.GetRoundResult(db.db, teams[3].ID, round)
	if res.AttackScore != 0.25 || res.DefenceScore != 2.0 {
		log.Fatalln("Invalid result:", res)
	}

}
