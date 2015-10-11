/**
 * @file receiver.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
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
	"tinfoilhat/steward"
	"tinfoilhat/vexillary"
)

const (
	GreetingMsg         string = "IBST.PSU CTF Flag Receiver\nInput flag: "
	InvalidFlagMsg      string = "Invalid flag\n"
	AlreadyCapturedMsg  string = "Flag already captured\n"
	CapturedMsg         string = "Captured!\n"
	InternalErrorMsg    string = "Internal error\n"
	FlagDoesNotExistMsg string = "Flag does not exist\n"
	FlagExpiredMsg      string = "Flag expired\n"
	InvalidTeamMsg      string = "Team does not exist\n"
	AttemptsLimitMsg    string = "Attack attempts limit exceeded\n"
	FlagYoursMsg        string = "Flag belongs to the attacking team\n"
	ServiceNotUpMsg     string = "The attacking team service is not up\n"
)

func ParseAddr(addr string) (subnet_no int, err error) {

	_, err = fmt.Sscanf(strings.Split(addr, ".")[2], "%d", &subnet_no)

	return
}

func TeamByAddr(db *sql.DB, addr string) (team steward.Team, err error) {

	subnet_no, err := ParseAddr(addr)
	if err != nil {
		return
	}

	teams, err := steward.GetTeams(db)
	if err != nil {
		return
	}

	for i := 0; i < len(teams); i++ {

		team = teams[i]

		team_subnet_no, err := ParseAddr(team.Subnet)
		if err != nil {
			return team, err
		}

		if team_subnet_no == subnet_no {
			return team, err
		}
	}

	err = errors.New("team not found")

	return
}

func Handler(conn net.Conn, db *sql.DB, priv *rsa.PrivateKey) {

	addr := conn.RemoteAddr().String()

	defer conn.Close()

	fmt.Fprint(conn, GreetingMsg)

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
		fmt.Fprint(conn, InvalidFlagMsg)
		return
	}

	exist, err := steward.FlagExist(db, flag)
	if err != nil {
		log.Println("\tExist flag check failed:", err)
		fmt.Fprint(conn, InternalErrorMsg)
		return
	}
	if !exist {
		fmt.Fprint(conn, FlagDoesNotExistMsg)
		return
	}

	flg, err := steward.GetFlagInfo(db, flag)
	if err != nil {
		log.Println("\tGet flag info failed:", err)
		fmt.Fprint(conn, InternalErrorMsg)
		return
	}

	captured, err := steward.AlreadyCaptured(db, flg.Id)
	if err != nil {
		log.Println("\tAlready captured check failed:", err)
		fmt.Fprint(conn, InternalErrorMsg)
		return
	}
	if captured {
		fmt.Fprint(conn, AlreadyCapturedMsg)
		return
	}

	team, err := TeamByAddr(db, addr)
	if err != nil {
		log.Println("\tGet team by ip failed:", err)
		fmt.Fprint(conn, InvalidTeamMsg)
		return
	}

	if flg.TeamId == team.Id {
		log.Printf("\tTeam %s try to send their flag", team.Name)
		fmt.Fprint(conn, FlagYoursMsg)
		return
	}

	halfStatus := steward.Status{flg.Round, team.Id, flg.ServiceId,
		steward.STATUS_UNKNOWN}
	state, err := steward.GetState(db, halfStatus)

	if state != steward.STATUS_UP {
		log.Printf("\t%s service not ok, cannot capture", team.Name)
		fmt.Fprint(conn, ServiceNotUpMsg)
		return
	}

	round, err := steward.CurrentRound(db)

	if round.Id != flg.Round {
		log.Printf("\t%s try to send flag from past round", team.Name)
		fmt.Fprint(conn, FlagExpiredMsg)
		return
	}

	round_end_time := round.StartTime.Add(round.Len)

	if time.Now().After(round_end_time) {
		log.Printf("\t%s try to send flag from finished round", team.Name)
		fmt.Fprint(conn, FlagExpiredMsg)
		return
	}

	err = steward.CaptureFlag(db, flg.Id, team.Id)
	if err != nil {
		log.Println("\tCapture flag failed:", err)
		fmt.Fprint(conn, InternalErrorMsg)
		return
	}

	fmt.Fprint(conn, CapturedMsg)
}

func FlagReceiver(db *sql.DB, priv *rsa.PrivateKey, addr string,
	timeout, socket_timeout time.Duration) {

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
			fmt.Fprint(conn, InternalErrorMsg)
			conn.Close()
			continue
		}

		if time.Now().Before(connects[ip].Add(timeout)) {
			log.Println("\tToo fast connects by", ip)
			fmt.Fprint(conn, AttemptsLimitMsg)
			conn.Close()
			continue
		}

		err = conn.SetDeadline(time.Now().Add(socket_timeout))
		if err != nil {
			log.Println("Set deadline fail:", err)
			fmt.Fprint(conn, InternalErrorMsg)
			conn.Close()
			continue
		}

		go Handler(conn, db, priv)

		connects[ip] = time.Now()
	}
}
