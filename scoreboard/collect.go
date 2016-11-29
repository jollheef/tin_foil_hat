/**
 * @file collect.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
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

import "github.com/jollheef/tin_foil_hat/steward"

var advisoryEnabled = true

// DisableAdvisory turn off advisory in overall score
func DisableAdvisory() {
	advisoryEnabled = false
}

func collectTeamResult(db *sql.DB, team steward.Team,
	services []steward.Service) (tr TeamResult, err error) {

	tr.ID = team.ID
	tr.Name = team.Name

	rr, err := steward.GetLastResult(db, team.ID)
	if err != nil {
		// At game start, no result exist
		rr = steward.RoundResult{AttackScore: 0, DefenceScore: 0}
	}

	tr.Attack = rr.AttackScore
	tr.Defence = rr.DefenceScore

	advisory, err := steward.GetAdvisoryScore(db, team.ID)
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
		s := steward.Status{round.ID, team.ID, svc.ID, -1}
		state, err := steward.GetState(db, s)
		if err != nil {
			// Try to get status from previous round
			s.Round--
			state, err = steward.GetState(db, s)
			if err != nil {
				state = steward.StatusDown
			}
		}

		tr.Status = append(tr.Status, state)
	}

	return
}

func max(r *Result, fn func(TeamResult) float64) (max float64) {
	max = 0
	for _, tr := range r.Teams {
		if fn(tr) > max {
			max = fn(tr)
		}
	}
	return
}

// CountScoreAndSort update scoreboard helper
func CountScoreAndSort(r *Result) {

	maxAttack := max(r,
		func(tr TeamResult) float64 { return tr.Attack })
	maxDefence := max(r,
		func(tr TeamResult) float64 { return tr.Defence })
	maxAdvisory := max(r,
		func(tr TeamResult) float64 { return float64(tr.Advisory) })

	if maxAttack == 0 {
		maxAttack = 1
	}

	if maxDefence == 0 {
		maxDefence = 1
	}

	if maxAdvisory == 0 {
		maxAdvisory = 1
	}

	for i := range r.Teams {

		tr := &r.Teams[i]

		tr.AttackPercent = tr.Attack / maxAttack * 100
		tr.DefencePercent = tr.Defence / maxDefence * 100
		tr.AdvisoryPercent = float64(tr.Advisory) / maxAdvisory * 100

		if advisoryEnabled {
			tr.Score = (tr.AttackPercent + tr.DefencePercent +
				tr.AdvisoryPercent) / 3
		} else {
			tr.Score = (tr.AttackPercent + tr.DefencePercent) / 2
		}
	}

	maxScore := max(r,
		func(tr TeamResult) float64 { return tr.Score })

	if maxScore == 0 {
		maxScore = 1
	}

	for i := range r.Teams {
		tr := &r.Teams[i]
		tr.ScorePercent = tr.Score / maxScore * 100
	}

	sort.Sort(ByScore(r.Teams))

	for i := range r.Teams {
		tr := &r.Teams[i]
		tr.Rank = i + 1
	}
}

// CollectLastResult returns actual scoreboard
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

		tr, err := collectTeamResult(db, team, services)
		if err != nil {
			return r, err
		}

		r.Teams = append(r.Teams, tr)
	}

	return
}
