/**
 * @file checker.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief functions for check services
 *
 * Provide functions for check service status, put flags and check flags.
 */

package checker

import (
	"crypto/rsa"
	"database/sql"
	"fmt"
	"log"
	"net"
	"sync"
)

import (
	"github.com/jollheef/tin_foil_hat/steward"
	"github.com/jollheef/tin_foil_hat/vexillary"
)

func tcpPortOpen(team steward.Team, svc steward.Service) bool {

	addr := fmt.Sprintf("%s:%d", team.Vulnbox, svc.Port)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return false
	} else {
		conn.Close()
		return true
	}
}

func putFlag(db *sql.DB, priv *rsa.PrivateKey, round int, team steward.Team,
	svc steward.Service) (err error) {

	flag, err := vexillary.GenerateFlag(priv)
	if err != nil {
		log.Println("Generate flag failed:", err)
		return
	}

	portOpen := true
	if !svc.Udp {
		portOpen = tcpPortOpen(team, svc)
	}

	var cred, logs string
	var state steward.ServiceState
	if portOpen {
		cred, logs, state, err = put(svc.CheckerPath, team.Vulnbox,
			svc.Port, flag)
		if err != nil {
			log.Println("Put flag to service failed:", err)
			return
		}

		if state != steward.STATUS_UP {
			log.Printf("Put flag, round %d, team %s, service %s: %s",
				round, team.Name, svc.Name, logs)
		}
	} else {
		state = steward.STATUS_DOWN
	}

	err = steward.PutStatus(db,
		steward.Status{round, team.Id, svc.Id, state})
	if err != nil {
		log.Println("Add status to database failed:", err)
		return
	}

	err = steward.AddFlag(db,
		steward.Flag{-1, flag, round, team.Id, svc.Id, cred})
	if err != nil {
		log.Println("Add flag to database failed:", err)
		return
	}

	return
}

func getFlag(db *sql.DB, round int, team steward.Team,
	svc steward.Service) (state steward.ServiceState, err error) {

	flag, cred, err := steward.GetCred(db, round, team.Id, svc.Id)
	if err != nil {
		log.Println("Get cred failed:", err)
		state = steward.STATUS_CORRUPT
		return
	}

	service_flag, logs, state, err := get(svc.CheckerPath, team.Vulnbox,
		svc.Port, cred)
	if err != nil {
		log.Println("Check service failed:", err)
		return
	}

	if flag != service_flag {
		state = steward.STATUS_CORRUPT
	}

	if state != steward.STATUS_UP {
		log.Printf("Get flag, round %d, team %s, service %s: %s",
			round, team.Name, svc.Name, logs)
	}

	return
}

func checkService(db *sql.DB, round int, team steward.Team,
	svc steward.Service) (state steward.ServiceState, err error) {

	state, logs, err := check(svc.CheckerPath, team.Vulnbox, svc.Port)
	if err != nil {
		log.Println("Check service failed:", err)
		return
	}

	if state != steward.STATUS_UP {
		log.Printf("Check, round %d, team %s, service %s: %s",
			round, team.Name, svc.Name, logs)
	}

	return
}

func PutFlags(db *sql.DB, priv *rsa.PrivateKey, round int,
	teams []steward.Team, services []steward.Service) (err error) {

	var wg sync.WaitGroup

	for _, team := range teams {
		for _, svc := range services {
			wg.Add(1)
			go func(team steward.Team, svc steward.Service) {
				defer wg.Done()
				putFlag(db, priv, round, team, svc)
			}(team, svc)
		}
	}

	wg.Wait()

	return
}

func CheckFlags(db *sql.DB, round int, teams []steward.Team,
	services []steward.Service) (err error) {

	var wg sync.WaitGroup

	for _, team := range teams {
		for _, svc := range services {
			wg.Add(1)
			go func(team steward.Team, svc steward.Service) {
				defer wg.Done()

				// First check service logic
				state, _ := checkService(db, round, team, svc)
				if state == steward.STATUS_UP {
					// If logic is correct, do flag check
					state, _ = getFlag(db, round, team, svc)
				}

				err = steward.PutStatus(db, steward.Status{round,
					team.Id, svc.Id, state})
				if err != nil {
					log.Println("Add status failed:", err)
					return
				}

			}(team, svc)
		}
	}

	wg.Wait()

	return
}
