/**
 * @file game.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
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

	"github.com/jollheef/tin_foil_hat/checker"
	"github.com/jollheef/tin_foil_hat/counter"
	"github.com/jollheef/tin_foil_hat/steward"
)

// RandomizeTimeout Do not forget something like rand.Seed(time.Now().UnixNano())
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

// Game contains game info
type Game struct {
	db       *sql.DB
	priv     *rsa.PrivateKey
	roundLen time.Duration
	timeout  time.Duration
	teams    []steward.Team
	services []steward.Service
}

// NewGame create new Game object
func NewGame(db *sql.DB, priv *rsa.PrivateKey, roundLen time.Duration,
	timeout time.Duration) (g Game, err error) {

	g.priv = priv
	g.roundLen = roundLen
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

// Over stop game
func (g Game) Over() {
	log.Println("Game over")
}

// Run start game
func (g Game) Run(end time.Time) (err error) {

	log.Println("Game start, end:", end)

	var counters sync.WaitGroup

	for time.Now().Add(g.roundLen).Before(end) {
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

// Round start new round
func (g Game) Round(counters *sync.WaitGroup) (err error) {

	roundNo, err := steward.NewRound(g.db, g.roundLen)
	if err != nil {
		return
	}

	log.Println("New round", roundNo)

	err = checker.PutFlags(g.db, g.priv, roundNo, g.teams, g.services)
	if err != nil {
		return
	}

	round, err := steward.CurrentRound(g.db)
	if err != nil {
		return
	}

	roundEnd := round.StartTime.Add(round.Len)

	for time.Now().Before(roundEnd) {

		log.Println("Round", round.ID, "check start")

		err = checker.CheckFlags(g.db, round.ID, g.teams, g.services)
		if err != nil {
			return
		}

		timeout := RandomizeTimeout(g.timeout, g.timeout/3)

		if time.Now().Add(timeout).After(roundEnd) {
			break
		}

		log.Println("Round", round.ID, "check end, timeout", timeout)

		time.Sleep(timeout)
	}

	log.Println("Check", round.ID, "over, wait", time.Now().Sub(roundEnd))

	for time.Now().Before(roundEnd) {
		time.Sleep(time.Second / 10)
	}

	counters.Add(1)
	go func() {
		defer counters.Done()

		log.Println("Count round", round.ID, "start", time.Now())

		err = counter.CountRound(g.db, round.ID, g.teams, g.services)
		if err != nil {
			log.Println("Count round", round.ID, "failed:", err)
		}

		log.Println("Count round", round.ID, "end", time.Now())
	}()

	return
}
