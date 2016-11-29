/**
 * @file api_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date November, 2016
 * @brief json api test
 */

package scoreboard

import (
	"encoding/json"
	"log"
	"net/http"
	"testing"
	"time"

	"golang.org/x/net/websocket"
)

func TestAttackFlowHandler(*testing.T) {

	attackFlow := make(chan Attack, 3)

	addr := "127.0.0.1:49000"
	apiURL := "/attack_flow_handler_test"

	http.Handle(apiURL, websocket.Handler(
		func(ws *websocket.Conn) {
			attackFlowHandler(ws, attackFlow)
		}))

	go func() {
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Second)

	go func() {
		for i := 0; i < 10; i++ {
			attackFlow <- Attack{i, i * 2, i * 3, int64(i * 4)}
		}
	}()

	var msg = make([]byte, 4096)

	ws, err := websocket.Dial("ws://"+addr+apiURL, "", "http://"+addr)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		var n int
		if n, err = ws.Read(msg); err != nil {
			panic(err)
		}

		var attack Attack
		err = json.Unmarshal(msg[0:n], &attack)
		if err != nil {
			panic(err)
		}

		attackEtalon := Attack{i, i * 2, i * 3, int64(i * 4)}

		if attack != attackEtalon {
			log.Println(attack, attackEtalon)
			panic("Attack not matched")
		}
	}
}
