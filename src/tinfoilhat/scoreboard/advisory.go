/**
 * @file scoreboard.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief web security advisory
 *
 * Contain web ui and several helpers for show advisory results
 */

package scoreboard

import (
	"database/sql"
	"fmt"
	"golang.org/x/net/websocket"
	"time"
)

import "tinfoilhat/steward"

func AdvisoryToHtml(adv steward.Advisory) (html string) {

	html = fmt.Sprintf("<h3>ISA-%d-%04d</h3>",
		adv.Timestamp.Year(), adv.Id)

	html += "<br><h4>Summary:</h4>"

	html += `<pre style="background-color: #000084; color: #ffffff">` +
		adv.Text + "</pre>"

	html += fmt.Sprintf("<h4>Published: %02d.%02d.%d %02d:%02d</h3>",
		adv.Timestamp.Day(),
		adv.Timestamp.Month(),
		adv.Timestamp.Year(),
		adv.Timestamp.Hour(),
		adv.Timestamp.Minute())

	html += fmt.Sprintf("<h4>Score: %d</h3><br>", adv.Score)

	return
}

var advisories string

func AdvisoryUpdater(db *sql.DB, update_timeout time.Duration) {

	for {
		var tmp_advisories string

		advs, err := steward.GetAdvisories(db)
		if err != nil {
			time.Sleep(update_timeout)
			continue
		}

		for i := range advs {
			adv := advs[len(advs)-i-1]
			if adv.Reviewed {
				tmp_advisories += AdvisoryToHtml(adv)
			}
		}

		if len(tmp_advisories) == 0 {
			advisories = "Current no advisories"
		} else {
			advisories = tmp_advisories
		}

		time.Sleep(update_timeout)
	}
}

func AdvisoryHandler(ws *websocket.Conn) {
	defer ws.Close()
	fmt.Fprintf(ws, advisories)
}
