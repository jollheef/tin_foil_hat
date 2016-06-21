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

func hasUnacceptableSymbols(s, regex string) bool {

	r, err := regexp.Compile(regex)
	if err != nil {
		log.Println("Compile regex fail:", err)
		return true
	}

	return !r.MatchString(s)
}

func advisoryHandler(conn net.Conn, db *sql.DB) {

	addr := conn.RemoteAddr().String()

	defer conn.Close()

	round, err := steward.CurrentRound(db)
	if err != nil {
		log.Println("Get current round fail:", err)
		fmt.Fprint(conn, internalErrorMsg)
		return
	}

	roundEndTime := round.StartTime.Add(round.Len)

	if time.Now().After(roundEndTime) {
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

	httpGetRoot := "GET / HTTP/1.1"
	if len(advisory) > len(httpGetRoot) {
		if advisory[0:len(httpGetRoot)] == httpGetRoot {
			fmt.Fprintf(conn, "\n\nIt's not a HTTP server! "+
				"Use netcat for communication.")
			return
		}
	}

	r := `[ -~]`
	if hasUnacceptableSymbols(advisory, r) {
		fmt.Fprintf(conn, "Accept only %s\n", r)
		return
	}

	team, err := teamByAddr(db, addr)
	if err != nil {
		log.Println("\tGet team by ip failed:", err)
		fmt.Fprint(conn, invalidTeamMsg)
		return
	}

	_, err = steward.AddAdvisory(db, team.ID, advisory)
	if err != nil {
		log.Println("\tAdd advisory failed:", err)
		fmt.Fprint(conn, internalErrorMsg)
		return
	}

	fmt.Fprint(conn, "Accepted\n")
}

// AdvisoryReceiver starts advisory receiver
func AdvisoryReceiver(db *sql.DB, addr string, timeout,
	socketTimeout time.Duration) {

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
			fmt.Fprint(conn, internalErrorMsg)
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

		err = conn.SetDeadline(time.Now().Add(socketTimeout))
		if err != nil {
			log.Println("Set deadline fail:", err)
			fmt.Fprint(conn, internalErrorMsg)
			conn.Close()
			continue
		}

		go advisoryHandler(conn, db)

		connects[ip] = time.Now()
	}
}
