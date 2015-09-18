/**
 * @file game.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief work with game
 */

package pulse

import (
	"crypto/rsa"
	"database/sql"
	"log"
	"math/rand"
	"sync"
	"time"
)

import (
	"tinfoilhat/checker"
	"tinfoilhat/counter"
	"tinfoilhat/steward"
)

// Do not forget something like rand.Seed(time.Now().UnixNano())
func RandomizeTimeout(duration, deviation time.Duration) (ret time.Duration) {

	if duration < time.Second || deviation < time.Second {
		return time.Second
	}

	defer func() {
		if r := recover(); r != nil {
			ret = time.Second
		}
	}()

	s := int(duration / time.Second)
	d := int(deviation / time.Second)

	min := s - s/d
	max := s + s/d

	val := min + rand.Intn(max-min)

	return time.Second * time.Duration(val)
}

type Game struct {
	db       *sql.DB
	priv     *rsa.PrivateKey
	roundLen time.Duration
	timeout  time.Duration
	teams    []steward.Team
	services []steward.Service
}

func NewGame(db *sql.DB, priv *rsa.PrivateKey, round_len time.Duration,
	timeout time.Duration) (g Game, err error) {

	g.priv = priv
	g.roundLen = round_len
	g.timeout = timeout
	g.db = db

	rand.Seed(time.Now().UnixNano())

	g.teams, err = steward.GetTeams(g.db)
	if err != nil {
		return
	}

	g.services, err = steward.GetServices(g.db)
	if err != nil {
		return
	}

	return
}

func (g Game) Over() {

	g.db.Close()

	log.Println("Game over")
}

func (g Game) Run(end time.Time) (err error) {

	log.Println("Game start, end:", end)

	var counters sync.WaitGroup

	for time.Now().Before(end) {
		err = g.Round(&counters)
		if err != nil {
			return
		}
	}

	log.Println("Wait counters")

	counters.Wait()

	log.Println("Game end")

	return
}

func (g Game) Round(counters *sync.WaitGroup) (err error) {

	round_no, err := steward.NewRound(g.db, g.roundLen)
	if err != nil {
		return
	}

	log.Println("New round", round_no)

	err = checker.PutFlags(g.db, g.priv, round_no, g.teams, g.services)
	if err != nil {
		return
	}

	round, err := steward.CurrentRound(g.db)
	if err != nil {
		return
	}

	round_end := round.StartTime.Add(round.Len)

	for time.Now().Before(round_end) {

		log.Println("Round", round.Id, "check start")

		err = checker.CheckFlags(g.db, round.Id, g.teams, g.services)
		if err != nil {
			return
		}

		timeout := RandomizeTimeout(g.timeout, g.timeout/3)

		if time.Now().Add(timeout).After(round_end) {
			break
		}

		log.Println("Round", round.Id, "check end, timeout", timeout)

		time.Sleep(timeout)
	}

	log.Println("Check", round.Id, "over, wait", time.Now().Sub(round_end))

	for time.Now().Before(round_end) {
		time.Sleep(time.Second / 10)
	}

	counters.Add(1)
	go func() {
		defer counters.Done()

		log.Println("Count round", round.Id, "start", time.Now())

		err = counter.CountRound(g.db, round.Id, g.teams, g.services)
		if err != nil {
			log.Println("Count round", round.Id, "failed:", err)
		}

		log.Println("Count round", round.Id, "end", time.Now())
	}()

	return
}
