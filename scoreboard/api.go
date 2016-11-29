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
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

// Attack describe attack for api
type Attack struct {
	Attacker  int
	Victim    int
	Service   int
	Timestamp int64
}

type broadcast struct {
	listeners  map[chan<- Attack]bool
	mutex      sync.Mutex
	attackFlow chan Attack
}

func newBroadcast(attackFlow chan Attack) *broadcast {
	b := broadcast{}
	b.mutex = sync.Mutex{}
	b.listeners = make(map[chan<- Attack]bool)
	b.attackFlow = attackFlow
	return &b
}

func (b *broadcast) Run() {
	for {
		if len(b.listeners) == 0 {
			time.Sleep(time.Second)
			continue
		}
		attack := <-b.attackFlow
		for l := range b.listeners {
			go func() { l <- attack }()
		}
	}
}

func (b *broadcast) NewListener() (c chan Attack) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	c = make(chan Attack, cap(b.attackFlow))
	b.listeners[c] = true
	return
}

func (b *broadcast) Detach(c chan Attack) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	delete(b.listeners, c)
	close(c)
}

func attackFlowHandler(ws *websocket.Conn, b *broadcast) {
	defer ws.Close()
	localFlow := b.NewListener()
	defer b.Detach(localFlow)
	for {
		attack := <-localFlow

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

func roundHandler(w http.ResponseWriter, r *http.Request) {
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
