/**
 * @file collect.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief collect results
 *
 * Collect all result in database, and calculate scoreboard table struct
 */

package scoreboard

import (
	"database/sql"
	"sort"
)

import "tinfoilhat/steward"

func CollectTeamResult(db *sql.DB, team steward.Team,
	services []steward.Service) (tr TeamResult, err error) {

	tr.Name = team.Name

	rr, err := steward.GetLastResult(db, team.Id)
	if err != nil {
		// At game start, no result exist
		rr = steward.RoundResult{AttackScore: 0, DefenceScore: 0}
	}

	tr.Attack = rr.AttackScore
	tr.Defence = rr.DefenceScore

	advisory, err := steward.GetAdvisoryScore(db, team.Id)
	if err != nil {
		tr.Advisory = 0
	} else {
		tr.Advisory = advisory
	}

	round, err := steward.CurrentRound(db)
	if err != nil {
		// At game start, no round exist
		return tr, nil
	}

	for _, svc := range services {
		s := steward.Status{round.Id, team.Id, svc.Id, -1}
		state, err := steward.GetState(db, s)
		if err != nil {
			state = steward.STATUS_DOWN
		}

		tr.Status = append(tr.Status, state)
	}

	return
}

func Max(r *Result, fn func(TeamResult) float64) (max float64) {
	max = 0
	for _, tr := range r.Teams {
		if fn(tr) > max {
			max = fn(tr)
		}
	}
	return
}

func CountScore(r *Result) {

	max_attack := Max(r,
		func(tr TeamResult) float64 { return tr.Attack })
	max_defence := Max(r,
		func(tr TeamResult) float64 { return tr.Defence })
	max_advisory := Max(r,
		func(tr TeamResult) float64 { return float64(tr.Advisory) })

	if max_attack == 0 {
		max_attack = 1
	}

	if max_defence == 0 {
		max_defence = 1
	}

	if max_advisory == 0 {
		max_advisory = 1
	}

	for i, _ := range r.Teams {

		tr := &r.Teams[i]

		tr.AttackPercent = tr.Attack / max_attack * 100
		tr.DefencePercent = tr.Defence / max_defence * 100
		tr.AdvisoryPercent = float64(tr.Advisory) / max_advisory * 100

		tr.Score = (tr.AttackPercent + tr.DefencePercent +
			tr.AdvisoryPercent) / 3
	}

	max_score := Max(r,
		func(tr TeamResult) float64 { return tr.Score })

	if max_score == 0 {
		max_score = 1
	}

	for i, _ := range r.Teams {
		tr := &r.Teams[i]
		tr.ScorePercent = tr.Score / max_score * 100
	}

	sort.Sort(ByScore(r.Teams))

	for i, _ := range r.Teams {
		tr := &r.Teams[i]
		tr.Rank = i + 1
	}
}

func CollectLastResult(db *sql.DB) (r Result, err error) {

	teams, err := steward.GetTeams(db)
	if err != nil {
		return
	}

	services, err := steward.GetServices(db)
	if err != nil {
		return
	}

	for _, svc := range services {
		r.Services = append(r.Services, svc.Name)
	}

	for _, team := range teams {

		tr, err := CollectTeamResult(db, team, services)
		if err != nil {
			return r, err
		}

		r.Teams = append(r.Teams, tr)
	}

	CountScore(&r)

	return
}
