/**
 * @file scoreboard_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief test work with scoreboard
 */

package scoreboard_test

import (
	"database/sql"
	"golang.org/x/net/websocket"
	"log"
	"sort"
	"sync"
	"testing"
	"time"
)

import (
	"github.com/jollheef/tin_foil_hat/scoreboard"
	"github.com/jollheef/tin_foil_hat/steward"
)

const (
	db_path  string = "user=postgres dbname=tinfoilhat_test sslmode=disable"
	www_path string = "www"
)

func TestCountScoreboard(*testing.T) {

	res := scoreboard.Result{}

	res.Teams = append(res.Teams, scoreboard.TeamResult{
		Attack:   10,
		Defence:  10,
		Advisory: 50})

	res.Teams = append(res.Teams, scoreboard.TeamResult{
		Attack:   0,
		Defence:  0,
		Advisory: 0})

	res.Teams = append(res.Teams, scoreboard.TeamResult{
		Attack:   100,
		Defence:  100,
		Advisory: 100})

	res.Teams = append(res.Teams, scoreboard.TeamResult{
		Attack:   0,
		Defence:  10,
		Advisory: 10})

	scoreboard.CountScoreAndSort(&res)

	sort.Sort(scoreboard.ByScore(res.Teams))

	for rank := 1; rank <= 4; rank++ {
		if res.Teams[rank-1].Rank != rank {
			log.Fatalln("team", rank, "rank is not", rank)
		}
	}

	if res.Teams[0].ScorePercent != 100 {
		log.Fatalln("First team score != 100%")
	}
}

func dialWebsocket(db *sql.DB, wg *sync.WaitGroup, i int) {

	origin := "http://localhost/"
	url := "ws://localhost:8080/scoreboard"

	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	res, err := scoreboard.CollectLastResult(db)
	if err != nil {
		log.Fatal(err)
	}

	html_res := res.ToHTML(false)

	var msg = make([]byte, len(html_res))
	if _, err = ws.Read(msg); err != nil {
		log.Fatal(err)
	}

	if string(msg) != html_res {
		log.Fatalln("Received result invalid",
			html_res, msg)
	}

	wg.Done()
}

func TestParallelWebSocketConnect(*testing.T) {

	db, err := steward.OpenDatabase(db_path)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(50) // default == 100

	err = steward.CleanDatabase(db)
	if err != nil {
		log.Fatal(err)
	}

	addr := ":8080"

	go func() {
		err := scoreboard.Scoreboard(db, www_path, addr, time.Second,
			time.Now(), time.Minute, time.Minute, time.Second)
		if err != nil {
			log.Fatal(err)
		}
	}()

	time.Sleep(time.Second) // wait scoreboard launching

	connects := 1000

	log.Printf("%d parallel connects\n", connects)

	var wg sync.WaitGroup
	for i := 0; i < connects; i++ {
		wg.Add(1)
		go dialWebsocket(db, &wg, i)
	}

	wg.Wait()
}
