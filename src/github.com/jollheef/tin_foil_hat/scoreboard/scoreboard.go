/**
 * @file scoreboard.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief web scoreboard
 *
 * Contain web ui and several helpers for convert round results to table
 */

package scoreboard

import (
	"database/sql"
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"time"
)

import "github.com/jollheef/tin_foil_hat/steward"

const (
	CONTEST_STATE_NOT_AVAILABLE = "state n/a"
	CONTEST_NOT_STARTED         = "not started"
	CONTEST_RUNNING             = "running"
	CONTEST_PAUSED              = "paused"
	CONTEST_COMPLETED           = "completed"
)

var (
	current_result string
	current_round  string
	last_updated   string
	contest_status string
	round          int
)

func ScoreboardHandler(ws *websocket.Conn) {

	defer ws.Close()

	fmt.Fprint(ws, current_result)
	sended_result := current_result
	last_update := time.Now()

	for {
		if sended_result != current_result ||
			time.Now().After(last_update.Add(time.Minute)) {

			sended_result = current_result
			last_update = time.Now()

			_, err := fmt.Fprint(ws, current_result)
			if err != nil {
				log.Println("Socket closed:", err)
				return
			}
		}

		time.Sleep(time.Second)
	}
}

func GetInfo() string {

	alert_type := ""

	if contest_status == CONTEST_RUNNING {
		alert_type = "alert-danger"
	}

	info := fmt.Sprintf(
		`<span class="alert %s">Contest %s</span>`+
			`<span class="alert">Round %d</span>`+
			`<span class="alert">Updated at %s</span>`,
		alert_type, contest_status, round, last_updated)

	return info
}

func InfoHandler(ws *websocket.Conn) {

	defer ws.Close()
	for {
		_, err := fmt.Fprint(ws, GetInfo())
		if err != nil {
			log.Println("Socket closed:", err)
			return
		}

		time.Sleep(time.Second)
	}
}

func ResultUpdater(db *sql.DB, update_timeout time.Duration,
	darkest_time time.Time) {

	for {
		res, err := CollectLastResult(db)
		if err != nil {
			log.Println("Collect last result fail:", err)
			time.Sleep(update_timeout)
			continue
		}

		if time.Now().Before(darkest_time) {
			CountScoreAndSort(&res)
			current_result = res.ToHTML(false)
		} else {
			current_result = res.ToHTML(true) // hide score
		}

		now := time.Now()
		last_updated = fmt.Sprintf("%02d:%02d:%02d", now.Hour(),
			now.Minute(), now.Second())

		r, err := steward.CurrentRound(db)
		if err != nil {
			round = 0
		} else {
			round = r.Id
		}

		time.Sleep(update_timeout)
	}
}

func StateUpdater(start, lunch_start_time, lunch_end_time, end_time time.Time,
	timeout time.Duration) {

	for {

		if time.Now().Before(start) {
			contest_status = CONTEST_NOT_STARTED
		} else if time.Now().Before(lunch_start_time) {
			contest_status = CONTEST_RUNNING
		} else if time.Now().Before(lunch_end_time) {
			contest_status = CONTEST_PAUSED
		} else if time.Now().Before(end_time) {
			contest_status = CONTEST_RUNNING
		} else {
			contest_status = CONTEST_COMPLETED
		}

		time.Sleep(timeout)
	}
}

func Scoreboard(db *sql.DB, www_path, addr string, update_timeout time.Duration,
	start time.Time, half, lunch, darkest time.Duration) (err error) {

	contest_status = CONTEST_STATE_NOT_AVAILABLE

	lunch_start := start.Add(half)
	lunch_end := lunch_start.Add(lunch)
	end_time := lunch_end.Add(half)

	darkest_time := end_time.Add(-darkest)

	go ResultUpdater(db, update_timeout, darkest_time)
	go StateUpdater(start, lunch_start, lunch_end, end_time, update_timeout)

	go AdvisoryUpdater(db, update_timeout)

	http.Handle("/scoreboard", websocket.Handler(ScoreboardHandler))
	http.Handle("/advisory", websocket.Handler(AdvisoryHandler))
	http.Handle("/info", websocket.Handler(InfoHandler))
	http.Handle("/", http.FileServer(http.Dir(www_path)))
	http.HandleFunc("/static-scoreboard", StaticScoreboard)

	log.Println("Launching scoreboard at", addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		return
	}

	return
}
