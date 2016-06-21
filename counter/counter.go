/**
 * @file counter.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date September, 2015
 * @brief cound round results
 *
 * Contain functions for work with round results
 */

package counter

import (
	"database/sql"

	"github.com/jollheef/tin_foil_hat/steward"
)

// CountStatesResult count round states (up/down/etc.) result
func CountStatesResult(db *sql.DB, round, team int,
	service steward.Service) (score float64, err error) {

	halfStatus := steward.Status{round, team, service.ID,
		steward.StatusUnknown}

	states, err := steward.GetStates(db, halfStatus)
	if err != nil {
		return
	}

	if len(states) == 0 {
		return
	}

	ok := 0.0
	for _, state := range states {
		if state == steward.StatusUP {
			ok++
		}
	}

	score = 1.0 / float64(len(states)) * ok

	return
}

// CountDefenceResult count round defence result
func CountDefenceResult(db *sql.DB, round, team int,
	services []steward.Service) (defence float64, err error) {

	defence = 0

	perService := 1.0 / float64(len(services))

	for _, svc := range services {

		score, err := CountStatesResult(db, round, team, svc)
		if err != nil {
			return defence, err
		}

		defence += score * perService
	}

	return
}

// CountRound count round result
func CountRound(db *sql.DB, round int, teams []steward.Team,
	services []steward.Service) (err error) {

	roundRes := make(map[steward.Team]steward.RoundResult)

	for _, team := range teams {

		res := steward.RoundResult{TeamID: team.ID, Round: round}

		def, err := CountDefenceResult(db, round, team.ID, services)
		if err != nil {
			return err
		}

		res.DefenceScore = def * 2

		roundRes[team] = res
	}

	perService := 1.0 / float64(len(services))

	for _, team := range teams {

		cflags, err := steward.GetCapturedFlags(db, round, team.ID)
		if err != nil {
			return err
		}

		for _, flag := range cflags {

			attackedTeam, err := steward.GetTeam(db, flag.TeamID)
			if err != nil {
				return err
			}

			res := roundRes[attackedTeam]
			res.DefenceScore -= perService
			if res.DefenceScore < 0 {
				res.DefenceScore = 0
			}
			roundRes[attackedTeam] = res

			attackRes := roundRes[team]
			attackRes.AttackScore += perService
			roundRes[team] = attackRes
		}
	}

	for _, res := range roundRes {
		_, err := steward.AddRoundResult(db, res)
		if err != nil {
			return err
		}

	}

	return
}
