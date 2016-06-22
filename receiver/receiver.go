/**
 * @file receiver.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date September, 2015
 * @brief routine for receive flags from commands
 *
 * Provide tcp server for receive flags. After receive flag daemon perform
 * validate flag, check flag round and write result to db.
 */

package receiver

import (
	"bufio"
	"crypto/rsa"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

import (
	"github.com/jollheef/tin_foil_hat/steward"
	"github.com/jollheef/tin_foil_hat/vexillary"
)

const (
	greetingMsg         string = "IBST.PSU CTF Flag Receiver\nInput flag: "
	invalidFlagMsg      string = "Invalid flag\n"
	alreadyCapturedMsg  string = "Flag already captured\n"
	capturedMsg         string = "Captured!\n"
	internalErrorMsg    string = "Internal error\n"
	flagDoesNotExistMsg string = "Flag does not exist\n"
	flagExpiredMsg      string = "Flag expired\n"
	invalidTeamMsg      string = "Team does not exist\n"
	attemptsLimitMsg    string = "Attack attempts limit exceeded\n"
	flagYoursMsg        string = "Flag belongs to the attacking team\n"
	serviceNotUpMsg     string = "The attacking team service is not up\n"
)

func parseAddr(addr string) (subnetNo int, err error) {

	defer func() {
		if r := recover(); r != nil {
			err = errors.New("Cannot parse '" + addr + "'")
		}
	}()
	_, err = fmt.Sscanf(strings.Split(addr, ".")[2], "%d", &subnetNo)

	return
}

func teamByAddr(db *sql.DB, addr string) (team steward.Team, err error) {

	subnetNo, err := parseAddr(addr)
	if err != nil {
		return
	}

	teams, err := steward.GetTeams(db)
	if err != nil {
		return
	}

	for i := 0; i < len(teams); i++ {

		team = teams[i]

		teamSubnetNo, err := parseAddr(team.Subnet)
		if err != nil {
			return team, err
		}

		if teamSubnetNo == subnetNo {
			return team, err
		}
	}

	err = errors.New("team not found")

	return
}

func handler(conn net.Conn, db *sql.DB, priv *rsa.PrivateKey) {

	addr := conn.RemoteAddr().String()

	defer conn.Close()

	fmt.Fprint(conn, greetingMsg)

	flag, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Println("Read error:", err)
	}

	flag = strings.Trim(flag, "\n")

	log.Printf("\tGet flag %s from %s", flag, addr)

	valid, err := vexillary.ValidFlag(flag, priv.PublicKey)
	if err != nil {
		log.Println("\tValidate flag failed:", err)
	}
	if !valid {
		fmt.Fprint(conn, invalidFlagMsg)
		return
	}

	exist, err := steward.FlagExist(db, flag)
	if err != nil {
		log.Println("\tExist flag check failed:", err)
		fmt.Fprint(conn, internalErrorMsg)
		return
	}
	if !exist {
		fmt.Fprint(conn, flagDoesNotExistMsg)
		return
	}

	flg, err := steward.GetFlagInfo(db, flag)
	if err != nil {
		log.Println("\tGet flag info failed:", err)
		fmt.Fprint(conn, internalErrorMsg)
		return
	}

	captured, err := steward.AlreadyCaptured(db, flg.ID)
	if err != nil {
		log.Println("\tAlready captured check failed:", err)
		fmt.Fprint(conn, internalErrorMsg)
		return
	}
	if captured {
		fmt.Fprint(conn, alreadyCapturedMsg)
		return
	}

	team, err := teamByAddr(db, addr)
	if err != nil {
		log.Println("\tGet team by ip failed:", err)
		fmt.Fprint(conn, invalidTeamMsg)
		return
	}

	if flg.TeamID == team.ID {
		log.Printf("\tTeam %s try to send their flag", team.Name)
		fmt.Fprint(conn, flagYoursMsg)
		return
	}

	halfStatus := steward.Status{flg.Round, team.ID, flg.ServiceID,
		steward.StatusUnknown}
	state, err := steward.GetState(db, halfStatus)

	if state != steward.StatusUP {
		log.Printf("\t%s service not ok, cannot capture", team.Name)
		fmt.Fprint(conn, serviceNotUpMsg)
		return
	}

	round, err := steward.CurrentRound(db)

	if round.ID != flg.Round {
		log.Printf("\t%s try to send flag from past round", team.Name)
		fmt.Fprint(conn, flagExpiredMsg)
		return
	}

	roundEndTime := round.StartTime.Add(round.Len)

	if time.Now().After(roundEndTime) {
		log.Printf("\t%s try to send flag from finished round", team.Name)
		fmt.Fprint(conn, flagExpiredMsg)
		return
	}

	err = steward.CaptureFlag(db, flg.ID, team.ID)
	if err != nil {
		log.Println("\tCapture flag failed:", err)
		fmt.Fprint(conn, internalErrorMsg)
		return
	}

	fmt.Fprint(conn, capturedMsg)
}

// FlagReceiver starts flag receiver
func FlagReceiver(db *sql.DB, priv *rsa.PrivateKey, addr string,
	timeout, socketTimeout time.Duration) {

	log.Println("Launching receiver at", addr, "...")

	connects := make(map[string]time.Time) // { ip : last_connect_time }

	listener, _ := net.Listen("tcp", addr)

	for {
		conn, _ := listener.Accept()

		addr := conn.RemoteAddr().String()

		log.Printf("Connection accepted from %s", addr)

		ip, _, err := net.SplitHostPort(addr)
		if err != nil {
			log.Println("\tCannot split remote addr:", err)
			fmt.Fprint(conn, internalErrorMsg)
			conn.Close()
			continue
		}

		if time.Now().Before(connects[ip].Add(timeout)) {
			log.Println("\tToo fast connects by", ip)
			fmt.Fprint(conn, attemptsLimitMsg)
			conn.Close()
			continue
		}

		err = conn.SetDeadline(time.Now().Add(socketTimeout))
		if err != nil {
			log.Println("Set deadline fail:", err)
			fmt.Fprint(conn, internalErrorMsg)
			conn.Close()
			continue
		}

		go handler(conn, db, priv)

		connects[ip] = time.Now()
	}
}
