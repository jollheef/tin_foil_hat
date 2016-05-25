/**
 * @file advisory.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date September, 2015
 * @brief routine for receive advisories from commands
 *
 * Provide tcp server for receive advisory.
 */

package receiver

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net"
	"regexp"
	"time"
)

import "github.com/jollheef/tin_foil_hat/steward"

func HasUnacceptableSymbols(s, regex string) bool {

	r, err := regexp.Compile(regex)
	if err != nil {
		log.Println("Compile regex fail:", err)
		return true
	}

	return !r.MatchString(s)
}

func AdvisoryHandler(conn net.Conn, db *sql.DB) {

	addr := conn.RemoteAddr().String()

	defer conn.Close()

	round, err := steward.CurrentRound(db)
	if err != nil {
		log.Println("Get current round fail:", err)
		fmt.Fprint(conn, InternalErrorMsg)
		return
	}

	round_end_time := round.StartTime.Add(round.Len)

	if time.Now().After(round_end_time) {
		fmt.Fprintln(conn, "Current contest not runned")
		return
	}

	fmt.Fprint(conn, "IBST.PSU CTF Advisory Receiver\n"+
		"Insert empty line for close\n"+
		"Input advisory: ")

	scanner := bufio.NewScanner(conn)
	var advisory string
	for scanner.Scan() {
		advisory += scanner.Text() + "\n"
		if len(advisory) > 2 {
			if advisory[len(advisory)-2:len(advisory)-1] == "\n" {
				// remove last newline
				advisory = advisory[:len(advisory)-1]
				break
			}
		}
	}

	http_get_root := "GET / HTTP/1.1"
	if len(advisory) > len(http_get_root) {
		if advisory[0:len(http_get_root)] == http_get_root {
			fmt.Fprintf(conn, "\n\nIt's not a HTTP server! "+
				"Use netcat for communication.")
			return
		}
	}

	r := `[ -~]`
	if HasUnacceptableSymbols(advisory, r) {
		fmt.Fprintf(conn, "Accept only %s\n", r)
		return
	}

	team, err := TeamByAddr(db, addr)
	if err != nil {
		log.Println("\tGet team by ip failed:", err)
		fmt.Fprint(conn, InvalidTeamMsg)
		return
	}

	_, err = steward.AddAdvisory(db, team.Id, advisory)
	if err != nil {
		log.Println("\tAdd advisory failed:", err)
		fmt.Fprint(conn, InternalErrorMsg)
		return
	}

	fmt.Fprint(conn, "Accepted\n")
}

func AdvisoryReceiver(db *sql.DB, addr string, timeout,
	socket_timeout time.Duration) {

	log.Println("Launching advisory receiver at", addr, "...")

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
			fmt.Fprintf(conn, "Attempts limit exceeded (wait %s)\n",
				connects[ip].Add(timeout).Sub(time.Now()))
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

		go AdvisoryHandler(conn, db)

		connects[ip] = time.Now()
	}
}
