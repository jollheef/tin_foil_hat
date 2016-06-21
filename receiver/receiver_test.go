/**
 * @file receiver_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date September, 2015
 * @brief test receiver package
 */

package receiver

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net"
	"strings"
	"testing"
	"time"
)

import "github.com/jollheef/tin_foil_hat/steward"
import "github.com/jollheef/tin_foil_hat/vexillary"

type testDB struct {
	db *sql.DB
}

const dbPath string = "user=postgres dbname=tinfoilhat_test sslmode=disable"

func openDB() (t testDB, err error) {

	t.db, err = steward.OpenDatabase(dbPath)
	if err != nil {
		return
	}

	err = steward.CleanDatabase(t.db)

	return
}

func (t testDB) Close() {

	t.db.Exec("DROP SCHEMA public CASCADE")
	t.db.Exec("CREATE SCHEMA public")

	t.db.Close()
}

func TestparseAddr(t *testing.T) {

	for id := 0; id < 255; id++ {
		addr := fmt.Sprintf("127.0.%d.1:44259", id)

		teamID, err := parseAddr(addr)
		if err != nil {
			log.Fatalln("Parse addr failed:", err)
		}

		if teamID != id {
			log.Fatalf("Parsed [%v] instead [%v]", teamID, id)
		}
	}
}

func TestteamByAddr(*testing.T) {

	db, err := openDB()
	if err != nil {
		log.Fatalln("Open database failed:", err)
	}

	defer db.Close()

	for i := 10; i < 15; i++ {

		subnet := fmt.Sprintf("127.0.%d.1/24", i)

		name := fmt.Sprintf("Team_%d", i)

		t := steward.Team{ID: -1, Name: name, Subnet: subnet,
			Vulnbox: subnet}

		teamID, err := steward.AddTeam(db.db, t)
		if err != nil {
			log.Fatalln("Add team failed:", err)
		}

		addr := fmt.Sprintf("127.0.%d.115:3542", i)

		team, err := teamByAddr(db.db, addr)
		if err != nil {
			log.Fatalln("Get team failed:", err)
		}

		if team.ID != teamID {
			log.Fatalf("Get team with id [%v] instead [%v]",
				team.ID, teamID)
		}
	}
}

func testFlag(addr, flag, response string) {

	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		log.Fatalln("Connect to receiver failed:", err)
	}

	goodMsg := strings.Split(greetingMsg, "\n")[0]

	msg, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatalln("Invalid greeting:", err)
	}

	if strings.Trim(goodMsg, "\n") != goodMsg {
		log.Fatalf("Invalid message [%v] instead [%v]", msg, goodMsg)
	}

	fmt.Fprint(conn, flag+"\n")

	msg, err = bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatalln("Invalid response:", err)
	}

	if msg != response {
		log.Fatalf("Invalid message [%v] instead [%v]",
			strings.Trim(msg, "\n"),
			strings.Trim(response, "\n"))
	}
}

func TestReceiver(*testing.T) {

	db, err := openDB()
	if err != nil {
		log.Fatalln("Open database failed:", err)
	}

	defer db.Close()

	priv, err := vexillary.GenerateKey()
	if err != nil {
		log.Fatalln("Generate key failed:", err)
	}

	addr := "127.0.0.1:65000"

	flag, err := vexillary.GenerateFlag(priv)
	if err != nil {
		log.Fatalln("Generate flag failed:", err)
	}

	err = steward.AddFlag(db.db, steward.Flag{-1, flag, 1, 8, 1, ""})
	if err != nil {
		log.Fatalln("Add flag failed:", err)
	}

	firstRound, err := steward.NewRound(db.db, time.Minute*2)
	if err != nil {
		log.Fatalln("New round failed:", err)
	}

	go FlagReceiver(db.db, priv, addr, time.Nanosecond, time.Minute)

	time.Sleep(time.Second) // wait for init listener

	// The attacker must appear to be a team (e.g. jury cannot attack)
	testFlag(addr, flag, invalidTeamMsg)

	t := steward.Team{ID: -1, Name: "TestTeam", Subnet: "127.0.0.1/24",
		Vulnbox: "1"}

	// Correct flag must be captured
	teamID, err := steward.AddTeam(db.db, t)
	if err != nil {
		log.Fatalln("Add team failed:", err)
	}

	serviceID := 1

	// Flag must be captured only if service status ok
	steward.PutStatus(db.db, steward.Status{firstRound, teamID, serviceID,
		steward.StatusUP})

	testFlag(addr, flag, capturedMsg)

	// Correct flag must be captured only one
	testFlag(addr, flag, alreadyCapturedMsg)

	// Incorrect (non-signed or signed on other key) flag must be invalid
	testFlag(addr, "1e7b642f2282886377d1655af6097dd6101eac5b=",
		invalidFlagMsg)

	// Correct flag that does not exist in database must not be captured
	newFlag, err := vexillary.GenerateFlag(priv)
	if err != nil {
		log.Fatalln("Generate flag failed:", err)
	}

	testFlag(addr, newFlag, flagDoesNotExistMsg)

	// Submitted flag does not belongs to the attacking team
	flag4, err := vexillary.GenerateFlag(priv)
	if err != nil {
		log.Fatalln("Generate flag failed:", err)
	}

	err = steward.AddFlag(db.db, steward.Flag{-1, flag4, 1, teamID, 1, ""})
	if err != nil {
		log.Fatalln("Add flag failed:", err)
	}

	testFlag(addr, flag4, flagYoursMsg)

	// Correct flag from another round must not be captured
	flag2, err := vexillary.GenerateFlag(priv)
	if err != nil {
		log.Fatalln("Generate flag failed:", err)
	}

	curRound, err := steward.CurrentRound(db.db)

	err = steward.AddFlag(db.db, steward.Flag{-1, flag2, curRound.ID, 8, 1, ""})
	if err != nil {
		log.Fatalln("Add flag failed:", err)
	}

	_, err = steward.NewRound(db.db, time.Minute*2)
	if err != nil {
		log.Fatalln("New round failed:", err)
	}

	testFlag(addr, flag2, flagExpiredMsg)

	// Correct flag from expired round must not be captured
	roundLen := time.Second
	roundID, err := steward.NewRound(db.db, roundLen)
	if err != nil {
		log.Fatalln("New round failed:", err)
	}

	flag3, err := vexillary.GenerateFlag(priv)
	if err != nil {
		log.Fatalln("Generate flag failed:", err)
	}

	err = steward.AddFlag(db.db, steward.Flag{-1, flag3, roundID, 8, 1, ""})
	if err != nil {
		log.Fatalln("Add flag failed:", err)
	}

	time.Sleep(roundLen) // wait end of round

	testFlag(addr, flag3, flagExpiredMsg)

	// If service status down flag must not be captured
	roundID, err = steward.NewRound(db.db, time.Minute)
	if err != nil {
		log.Fatalln("New round failed:", err)
	}

	flag5, err := vexillary.GenerateFlag(priv)
	if err != nil {
		log.Fatalln("Generate flag failed:", err)
	}

	err = steward.AddFlag(db.db, steward.Flag{-1, flag5, roundID, 8,
		serviceID, ""})
	if err != nil {
		log.Fatalln("Add flag failed:", err)
	}

	steward.PutStatus(db.db, steward.Status{roundID, teamID, serviceID,
		steward.StatusDown})

	testFlag(addr, flag5, serviceNotUpMsg)

	steward.PutStatus(db.db, steward.Status{roundID, teamID, serviceID,
		steward.StatusUP})

	// If attempts limit exceeded flag must not be captured
	newAddr := "127.0.0.1:64000"

	// Start new receiver for test timeouts
	go FlagReceiver(db.db, priv, newAddr, time.Second, time.Minute)

	time.Sleep(time.Second) // wait for init listener

	// Just for take timeout
	testFlag(newAddr, flag3, flagExpiredMsg)

	// Can't use testFlag, if attempts limit exceeded server does not send
	// greeting message, and client has not able to send flag
	conn, err := net.DialTimeout("tcp", newAddr, time.Second)
	if err != nil {
		log.Fatalln("Connect to receiver failed:", err)
	}

	msg, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatalln("Invalid response:", err)
	}

	response := attemptsLimitMsg

	if msg != response {
		log.Fatalf("Invalid message [%v] instead [%v]",
			strings.Trim(msg, "\n"),
			strings.Trim(response, "\n"))
	}
}
