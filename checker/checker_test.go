/**
 * @file checker_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief test checker package
 */

package checker_test

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os/exec"
	"testing"
	"time"
)

import (
	"github.com/jollheef/tin_foil_hat/checker"
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

func (svc dummyService) SendCmd(command string) (err error) {

	addr := fmt.Sprintf("127.0.0.1:%d", svc.port)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}

	defer conn.Close()

	fmt.Fprint(conn, command+"\n")

	return
}

func (svc dummyService) BrokeLogic() {
	log.Println("Broke service logic")
	svc.SendCmd("REGFAIL")
}

func (svc dummyService) RestoreLogic() {
	log.Println("Restore service logic")
	svc.SendCmd("REGOK")
}

func (svc dummyService) ClearFlags() {
	log.Println("Clear flags")
	svc.SendCmd("CLEAR")
}

func checkServicesStatus(db *sql.DB, round int, teams []steward.Team,
	services []steward.Service, status steward.ServiceState) {

	for _, team := range teams {
		for _, svc := range services {
			halfStatus := steward.Status{round, team.Id, svc.Id, -1}

			state, err := steward.GetState(db, halfStatus)
			if err != nil {
				log.Fatalln("Get state failed:", err)
			}

			if state != status {
				log.Fatalln("One of service status is", state,
					"instead", status)
			}
		}
	}
}

func TestFlagsWork(t *testing.T) {

	db, err := openDB()
	if err != nil {
		log.Fatalln("Open database failed:", err)
	}

	defer db.Close()

	min_port_num := 1024
	max_port_num := 65535

	rand.Seed(time.Now().UnixNano())

	port := min_port_num + rand.Intn(max_port_num-min_port_num)

	log.Println("Use port", port)

	service := newDummyService("python-api/dummy_service.py", port)
	service.Stop() // if already run

	priv, err := vexillary.GenerateKey()
	if err != nil {
		log.Fatalln("Generate key failed:", err)
	}

	for index, team := range []string{"FooTeam", "BarTeam", "BazTeam"} {

		// just trick for bypass UNIQUE team subnet
		subnet := fmt.Sprintf("127.%d.0.1/24", index)

		vulnbox := fmt.Sprintf("127.0.%d.3", index)

		t := steward.Team{-1, team, subnet, vulnbox}

		_, err = steward.AddTeam(db.db, t)
		if err != nil {
			log.Fatalln("Add team failed:", err)
		}
	}

	checker_path := "python-api/dummy_checker.py"

	for _, service := range []string{"Foo", "Bar", "Baz"} {

		err = steward.AddService(db.db,
			steward.Service{-1, service, port, checker_path, false})
		if err != nil {
			log.Fatalln("Add service failed:", err)
		}
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

	err = checker.PutFlags(db.db, priv, round, teams, services)
	if err != nil {
		log.Fatalln("Put flags failed:", err)
	}

	// No services -> all down
	checkServicesStatus(db.db, round, teams, services, steward.STATUS_DOWN)

	// Start service
	service.Start()

	time.Sleep(time.Second)

	service.BrokeLogic()

	err = checker.PutFlags(db.db, priv, round, teams, services)
	if err != nil {
		log.Fatalln("Put flags failed:", err)
	}

	checkServicesStatus(db.db, round, teams, services, steward.STATUS_MUMBLE)

	service.RestoreLogic()

	time.Sleep(time.Second)

	round, err = steward.NewRound(db.db, time.Minute)
	if err != nil {
		log.Fatalln("Create new round failed:", err)
	}

	log.Println("Put flags to correct service...")

	err = checker.PutFlags(db.db, priv, round, teams, services)
	if err != nil {
		log.Fatalln("Put flags failed:", err)
	}

	checkServicesStatus(db.db, round, teams, services, steward.STATUS_UP)

	log.Println("Check flags of correct service...")

	err = checker.CheckFlags(db.db, round, teams, services)
	if err != nil {
		log.Fatalln("Check flags failed:", err)
	}

	checkServicesStatus(db.db, round, teams, services, steward.STATUS_UP)

	service.ClearFlags()

	log.Println("Check flags of service with removed flags...")

	err = checker.CheckFlags(db.db, round, teams, services)
	if err != nil {
		log.Fatalln("Check flags failed:", err)
	}

	checkServicesStatus(db.db, round, teams, services, steward.STATUS_CORRUPT)

	log.Println("Stop service...")
	service.Stop()

	log.Println("Check flags of stopped service...")

	err = checker.CheckFlags(db.db, round, teams, services)
	if err != nil {
		log.Fatalln("Check flags failed:", err)
	}

	checkServicesStatus(db.db, round, teams, services, steward.STATUS_DOWN)
}
