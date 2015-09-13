/**
 * @file receiver_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief test receiver package
 */

package receiver_test

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"testing"
	"time"
)

import "tinfoilhat/receiver"
import "tinfoilhat/steward"
import "tinfoilhat/vexillary"

func TestParseAddr(t *testing.T) {

	for id := 0; id < 255; id++ {
		addr := fmt.Sprintf("127.0.%d.1:44259", id)

		team_id, err := receiver.ParseAddr(addr)
		if err != nil {
			log.Fatalln("Parse addr failed:", err)
		}

		if team_id != id {
			log.Fatalf("Parsed [%v] instead [%v]", team_id, id)
		}
	}
}

func TestTeamByAddr(t *testing.T) {

	db, err := steward.PrivateOpenDatabase(":memory:")
	if err != nil {
		log.Fatalln("Open database failed:", err)
	}

	defer db.Close()

	for i := 10; i < 15; i++ {

		subnet := fmt.Sprintf("127.0.%d.1/24", i)

		name := fmt.Sprintf("Team_%d", i)

		team_id, err := steward.AddTeam(db, name, subnet)
		if err != nil {
			log.Fatalln("Add team failed:", err)
		}

		addr := fmt.Sprintf("127.0.%d.115:3542", i)

		team, err := receiver.TeamByAddr(db, addr)
		if err != nil {
			log.Fatalln("Get team failed:", err)
		}

		if team.Id != team_id {
			log.Fatalf("Get team with id [%v] instead [%v]",
				team.Id, team_id)
		}
	}
}

func testFlag(addr, flag, response string) {

	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		log.Fatalln("Connect to receiver failed:", err)
	}

	good_msg := strings.Split(receiver.GreetingMsg, "\n")[0]

	msg, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatalln("Invalid greeting:", err)
	}

	if strings.Trim(good_msg, "\n") != good_msg {
		log.Fatalf("Invalid message [%v] instead [%v]", msg, good_msg)
	}

	fmt.Fprintf(conn, flag+"\n")

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

func TestReceiver(t *testing.T) {

	db, err := steward.PrivateOpenDatabase(":memory:")
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

	err = steward.AddFlag(db, steward.Flag{-1, flag, 1, 8, 1})
	if err != nil {
		log.Fatalln("Add flag failed:", err)
	}

	first_round, err := steward.NewRound(db, time.Minute*2)
	if err != nil {
		log.Fatalln("New round failed:", err)
	}

	go receiver.Receiver(db, priv, addr, time.Nanosecond)

	time.Sleep(time.Second) // wait for init listener

	// The attacker must appear to be a team (e.g. jury cannot attack)
	testFlag(addr, flag, receiver.InvalidTeamMsg)

	// Correct flag must be captured
	team_id, err := steward.AddTeam(db, "TestTeam", "127.0.0.1/24")
	if err != nil {
		log.Fatalln("Add team failed:", err)
	}

	service_id := 1

	// Flag must be captured only if service status ok
	steward.PutStatus(db, steward.Status{first_round, team_id, service_id,
		steward.STATUS_OK})

	testFlag(addr, flag, receiver.CapturedMsg)

	// Correct flag must be captured only one
	testFlag(addr, flag, receiver.AlreadyCapturedMsg)

	// Incorrect (non-signed or signed on other key) flag must be invalid
	testFlag(addr, "1e7b642f2282886377d1655af6097dd6101eac5b=",
		receiver.InvalidFlagMsg)

	// Correct flag that does not exist in database must not be captured
	new_flag, err := vexillary.GenerateFlag(priv)
	if err != nil {
		log.Fatalln("Generate flag failed:", err)
	}

	testFlag(addr, new_flag, receiver.FlagDoesNotExistMsg)

	// Submitted flag does not belongs to the attacking team
	flag4, err := vexillary.GenerateFlag(priv)
	if err != nil {
		log.Fatalln("Generate flag failed:", err)
	}

	err = steward.AddFlag(db, steward.Flag{-1, flag4, 1, team_id, 1})
	if err != nil {
		log.Fatalln("Add flag failed:", err)
	}

	testFlag(addr, flag4, receiver.FlagYoursMsg)

	// Correct flag from another round must not be captured
	flag2, err := vexillary.GenerateFlag(priv)
	if err != nil {
		log.Fatalln("Generate flag failed:", err)
	}

	cur_round, err := steward.CurrentRound(db)

	err = steward.AddFlag(db, steward.Flag{-1, flag2, cur_round.Id, 8, 1})
	if err != nil {
		log.Fatalln("Add flag failed:", err)
	}

	_, err = steward.NewRound(db, time.Minute*2)
	if err != nil {
		log.Fatalln("New round failed:", err)
	}

	testFlag(addr, flag2, receiver.FlagExpiredMsg)

	// Correct flag from expired round must not be captured
	round_len := time.Second
	round_id, err := steward.NewRound(db, round_len)
	if err != nil {
		log.Fatalln("New round failed:", err)
	}

	flag3, err := vexillary.GenerateFlag(priv)
	if err != nil {
		log.Fatalln("Generate flag failed:", err)
	}

	err = steward.AddFlag(db, steward.Flag{-1, flag3, round_id, 8, 1})
	if err != nil {
		log.Fatalln("Add flag failed:", err)
	}

	time.Sleep(round_len) // wait end of round

	testFlag(addr, flag3, receiver.FlagExpiredMsg)

	// If service status down flag must not be captured
	round_id, err = steward.NewRound(db, time.Minute)
	if err != nil {
		log.Fatalln("New round failed:", err)
	}

	flag5, err := vexillary.GenerateFlag(priv)
	if err != nil {
		log.Fatalln("Generate flag failed:", err)
	}

	err = steward.AddFlag(db, steward.Flag{-1, flag5, round_id, 8,
		service_id})
	if err != nil {
		log.Fatalln("Add flag failed:", err)
	}

	steward.PutStatus(db, steward.Status{round_id, team_id, service_id,
		steward.STATUS_DOWN})

	testFlag(addr, flag5, receiver.ServiceNotUpMsg)

	steward.PutStatus(db, steward.Status{round_id, team_id, service_id,
		steward.STATUS_OK})

	// If attempts limit exceeded flag must not be captured
	new_addr := "127.0.0.1:64000"

	// Start new receiver for test timeouts
	go receiver.Receiver(db, priv, new_addr, time.Second)

	time.Sleep(time.Second) // wait for init listener

	// Just for take timeout
	testFlag(new_addr, flag3, receiver.FlagExpiredMsg)

	// Can't use testFlag, if attempts limit exceeded server does not send
	// greeting message, and client has not able to send flag
	conn, err := net.DialTimeout("tcp", new_addr, time.Second)
	if err != nil {
		log.Fatalln("Connect to receiver failed:", err)
	}

	msg, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatalln("Invalid response:", err)
	}

	response := receiver.AttemptsLimitMsg

	if msg != response {
		log.Fatalf("Invalid message [%v] instead [%v]",
			strings.Trim(msg, "\n"),
			strings.Trim(response, "\n"))
	}
}
