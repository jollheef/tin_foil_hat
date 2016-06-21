/**
 * @file scoreboard.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date September, 2015
 * @brief web scoreboard
 *
 * Contain web ui and several helpers for convert round results to table
 */

package scoreboard

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

import "github.com/jollheef/tin_foil_hat/steward"

const (
	contestStateNotAvailable = "state n/a"
	contestNotStarted        = "not started"
	contestRunning           = "running"
	contestPaused            = "paused"
	contestCompleted         = "completed"
)

var (
	currentResult string
	currentRound  string
	lastUpdated   string
	contestStatus string
	round         int
)

func scoreboardHandler(ws *websocket.Conn) {

	defer ws.Close()

	fmt.Fprint(ws, currentResult)
	sendedResult := currentResult
	lastUpdate := time.Now()

	for {
		if sendedResult != currentResult ||
			time.Now().After(lastUpdate.Add(time.Minute)) {

			sendedResult = currentResult
			lastUpdate = time.Now()

			_, err := fmt.Fprint(ws, currentResult)
			if err != nil {
				log.Println("Socket closed:", err)
				return
			}
		}

		time.Sleep(time.Second)
	}
}

func getInfo() string {

	alertType := ""

	if contestStatus == contestRunning {
		alertType = "alert-danger"
	}

	info := fmt.Sprintf(
		`<span class="alert %s">Contest %s</span>`+
			`<span class="alert">Round %d</span>`+
			`<span class="alert">Updated at %s</span>`,
		alertType, contestStatus, round, lastUpdated)

	return info
}

func infoHandler(ws *websocket.Conn) {

	defer ws.Close()
	for {
		_, err := fmt.Fprint(ws, getInfo())
		if err != nil {
			log.Println("Socket closed:", err)
			return
		}

		time.Sleep(time.Second)
	}
}

func resultUpdater(db *sql.DB, updateTimeout time.Duration,
	darkestTime time.Time) {

	for {
		res, err := CollectLastResult(db)
		if err != nil {
			log.Println("Collect last result fail:", err)
			time.Sleep(updateTimeout)
			continue
		}

		if time.Now().Before(darkestTime) {
			CountScoreAndSort(&res)
			currentResult = res.ToHTML(false)
		} else {
			currentResult = res.ToHTML(true) // hide score
		}

		now := time.Now()
		lastUpdated = fmt.Sprintf("%02d:%02d:%02d", now.Hour(),
			now.Minute(), now.Second())

		r, err := steward.CurrentRound(db)
		if err != nil {
			round = 0
		} else {
			round = r.ID
		}

		time.Sleep(updateTimeout)
	}
}

func stateUpdater(start, lunchStartTime, lunchEndTime, endTime time.Time,
	timeout time.Duration) {

	for {

		if time.Now().Before(start) {
			contestStatus = contestNotStarted
		} else if time.Now().Before(lunchStartTime) {
			contestStatus = contestRunning
		} else if time.Now().Before(lunchEndTime) {
			contestStatus = contestPaused
		} else if time.Now().Before(endTime) {
			contestStatus = contestRunning
		} else {
			contestStatus = contestCompleted
		}

		time.Sleep(timeout)
	}
}

// Scoreboard run scoreboard page
func Scoreboard(db *sql.DB, wwwPath, addr string, updateTimeout time.Duration,
	start time.Time, half, lunch, darkest time.Duration) (err error) {

	contestStatus = contestStateNotAvailable

	lunchStart := start.Add(half)
	lunchEnd := lunchStart.Add(lunch)
	endTime := lunchEnd.Add(half)

	darkestTime := endTime.Add(-darkest)

	go resultUpdater(db, updateTimeout, darkestTime)
	go stateUpdater(start, lunchStart, lunchEnd, endTime, updateTimeout)

	go advisoryUpdater(db, updateTimeout)

	http.Handle("/scoreboard", websocket.Handler(scoreboardHandler))
	http.Handle("/advisory", websocket.Handler(advisoryHandler))
	http.Handle("/info", websocket.Handler(infoHandler))
	http.Handle("/", http.FileServer(http.Dir(wwwPath)))
	http.HandleFunc("/static-scoreboard", staticScoreboard)

	log.Println("Launching scoreboard at", addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		return
	}

	return
}
