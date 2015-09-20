/**
 * @file scoreboard_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief test work with scoreboard
 */

package scoreboard_test

import (
	"log"
	"sort"
	"testing"
)

import "tinfoilhat/scoreboard"

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

	scoreboard.CountScore(&res)

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
