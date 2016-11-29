/**
 * @file api.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date November, 2016
 * @brief json api
 */

package scoreboard

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

// Attack describe attack for api
type Attack struct {
	Attacker  int
	Victim    int
	Service   int
	Timestamp int64
}

func attackFlowHandler(ws *websocket.Conn, attackFlow chan Attack) {
	defer ws.Close()
	for {
		attack := <-attackFlow

		buf, err := json.Marshal(attack)
		if err != nil {
			log.Println("Serialization error:", err)
			return
		}

		_, err = ws.Write(buf)
		if err != nil {
			log.Println("Attack flow write error:", err)
			return
		}
	}
}

func resultHandler(w http.ResponseWriter, r *http.Request) {
	buf, err := json.Marshal(lastResult)
	if err != nil {
		log.Println("Serialization error:", err)
		return
	}

	_, err = w.Write(buf)
	if err != nil {
		log.Println("Result write error:", err)
		return
	}
}
