/**
 * @file pulse.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief generate game events
 *
 * Routine that generates game events
 */

package pulse

import (
	"crypto/rsa"
	"database/sql"
	"log"
	"time"
)

// Wait for time
func Wait(end time.Time, timeout time.Duration) (waited bool) {

	if time.Now().After(end) {
		return false
	}

	for time.Now().Before(end) {
		time.Sleep(timeout)
	}

	return true
}

// Pulse manage game
func Pulse(db *sql.DB, priv *rsa.PrivateKey, startTime time.Time,
	half, lunch, roundLen, checkTimeout time.Duration) (err error) {

	log.Println("Launching pulse...")

	lunchStartTime := startTime.Add(half)
	lunchEndTime := lunchStartTime.Add(lunch)
	endTime := lunchEndTime.Add(half)

	log.Println("Pulse start time", time.Now())

	log.Println("Contest start time", startTime)

	game, err := NewGame(db, priv, roundLen, checkTimeout)

	defer game.Over()

	timeout := 100 * time.Millisecond

	log.Println("Wait start time...")
	if Wait(startTime, timeout) || time.Now().Before(lunchStartTime) {
		log.Println("game run")
		err = game.Run(lunchStartTime)
		if err != nil {
			return
		}
	}

	log.Println("Wait lunch time")
	Wait(lunchStartTime, timeout)

	log.Println("Lunch...")
	if Wait(lunchEndTime, timeout) || time.Now().Before(endTime) {
		err = game.Run(endTime)
		if err != nil {
			return
		}
	}

	log.Println("Wait end time")
	Wait(endTime, timeout)

	return
}
