/**
 * @file scoreboard.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date September, 2015
 * @brief web security advisory
 *
 * Contain web ui and several helpers for show advisory results
 */

package scoreboard

import (
	"database/sql"
	"fmt"
	"html/template"
	"time"

	"golang.org/x/net/websocket"
)

import "github.com/jollheef/tin_foil_hat/steward"

func advisoryToHTML(adv steward.Advisory) (html string) {

	html = fmt.Sprintf("<h3>ISA-%d-%04d</h3>",
		adv.Timestamp.Year(), adv.ID)

	html += "<br><h4>Summary:</h4>"

	html += `<pre style="background-color: #000084; color: #ffffff">` +
		template.HTMLEscapeString(adv.Text) + "</pre>"

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

func advisoryUpdater(db *sql.DB, updateTimeout time.Duration) {

	for {
		var tmpAdvisories string

		advs, err := steward.GetAdvisories(db)
		if err != nil {
			time.Sleep(updateTimeout)
			continue
		}

		for i := range advs {
			adv := advs[len(advs)-i-1]
			if adv.Reviewed {
				tmpAdvisories += advisoryToHTML(adv)
			}
		}

		if len(tmpAdvisories) == 0 {
			advisories = "Current no advisories"
		} else {
			advisories = tmpAdvisories
		}

		time.Sleep(updateTimeout)
	}
}

func advisoryHandler(ws *websocket.Conn) {
	defer ws.Close()
	fmt.Fprint(ws, advisories)
}
