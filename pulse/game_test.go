/**
 * @file game_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date September, 2015
 * @brief test game struct
 */

package pulse_test

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"testing"
	"time"
)

import (
	"github.com/jollheef/tin_foil_hat/pulse"
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

type dummyService struct {
	path string
	port int
}

func newDummyService(path string, port int) (svc dummyService) {
	svc.path = path
	svc.port = port
	return
}

func (svc dummyService) Start() {
	port := fmt.Sprintf("%d", svc.port)
	exec.Command(svc.path, port).Start()
}

func (svc dummyService) Stop() {
	exec.Command("pkill", "-f", svc.path).Run()
}

func TestRandomizeTimeout(*testing.T) {

	rand.Seed(time.Now().UnixNano())

	if pulse.RandomizeTimeout(0, 0) != time.Second {
		log.Fatalln("Invalid value for invalid parameters")
	}

	timeout := 10 * time.Second
	deviation := timeout / 3

	for i := 0; i < 100; i++ {
		val := pulse.RandomizeTimeout(timeout, deviation)
		if val > timeout+deviation || val < timeout-deviation {
			log.Fatalln("Value is not within the expected range")
		}
	}
}

func TestGame(*testing.T) {

	db, err := openDB()
	if err != nil {
		log.Fatalln("Open database fail:", err)
	}

	min_port_num := 1024
	max_port_num := 65535

	rand.Seed(time.Now().UnixNano())

	port := min_port_num + rand.Intn(max_port_num-min_port_num)

	log.Println("Use port", port)

	service_path := "../checker/python-api/dummy_service.py"

	svc := newDummyService(service_path, port)
	svc.Stop() // kill
	svc.Start()

	time.Sleep(time.Second) // wait for service init

	defer svc.Stop()

	priv, err := vexillary.GenerateKey()
	if err != nil {
		log.Fatalln("Generate key fail:", err)
	}

	for index, team := range []string{"FooTeam", "BarTeam", "BazTeam"} {

		// just trick for bypass UNIQUE team subnet
		subnet := fmt.Sprintf("127.%d.0.1/24", index)

		vulnbox := fmt.Sprintf("127.0.%d.3", index)

		t := steward.Team{ID: -1, Name: team, Subnet: subnet,
			Vulnbox: vulnbox}

		_, err = steward.AddTeam(db.db, t)
		if err != nil {
			log.Fatalln("Add team failed:", err)
		}
	}

	checker_path := "../checker/python-api/dummy_checker.py"

	for _, service := range []string{"Foo", "Bar", "Baz"} {

		err = steward.AddService(db.db,
			steward.Service{ID: -1, Name: service, Port: port,
				CheckerPath: checker_path, UDP: false})
		if err != nil {
			log.Fatalln("Add service failed:", err)
		}
	}

	round_len := 30 * time.Second
	timeout_between_check := 10 * time.Second

	game, err := pulse.NewGame(db.db, priv, round_len, timeout_between_check)

	defer game.Over()

	end_time := time.Now().Add(time.Minute + 10*time.Second)

	err = game.Run(end_time)

	if err != nil {
		log.Fatalln("Game error:", err)
	}

	for round := 1; round <= 2; round++ {

		for team_id := 1; team_id <= 3; team_id++ {

			res, err := steward.GetRoundResult(db.db, team_id, round)
			if err != nil {
				log.Fatalf("Get round %d result fail: %s",
					round, err)
			}

			if res.DefenceScore != float64(round*2) {
				log.Fatalln("Invalid defence score")
			}

			if res.AttackScore != 0 {
				log.Fatalln("Invalid attack score")
			}
		}
	}
}
