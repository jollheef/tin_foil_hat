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

var current_result string

func ScoreboardHandler(ws *websocket.Conn) {

	fmt.Fprintf(ws, current_result)
	sended_result := current_result
	last_update := time.Now()

	for {
		if sended_result != current_result ||
			time.Now().After(last_update.Add(time.Minute)) {

			sended_result = current_result
			last_update = time.Now()

			_, err := fmt.Fprintf(ws, current_result)
			if err != nil {
				log.Println("Socket closed:", err)
				return
			}
		}

		time.Sleep(time.Second)
	}
}

func ResultUpdater(db *sql.DB, update_timeout time.Duration) {

	for {
		res, err := CollectLastResult(db)
		if err != nil {
			log.Println("Collect last result fail:", err)
			time.Sleep(update_timeout)
			continue
		}

		current_result = res.ToHTML()

		time.Sleep(update_timeout)
	}
}

func Scoreboard(db *sql.DB, www_path, addr string,
	update_timeout time.Duration) (err error) {

	go ResultUpdater(db, update_timeout)

	http.Handle("/scoreboard", websocket.Handler(ScoreboardHandler))
	http.Handle("/", http.FileServer(http.Dir(www_path)))

	log.Println("Launching scoreboard at", addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		return
	}

	return
}
