/**
 * @file counter.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief cound round results
 *
 * Contain functions for work with round results
 */

package counter

import "database/sql"

import "github.com/jollheef/tin_foil_hat/steward"

func CountStatesResult(db *sql.DB, round, team int,
	service steward.Service) (score float64, err error) {

	halfStatus := steward.Status{round, team, service.Id,
		steward.STATUS_UNKNOWN}

	states, err := steward.GetStates(db, halfStatus)
	if err != nil {
		return
	}

	if len(states) == 0 {
		return
	}

	ok := 0.0
	for _, state := range states {
		if state == steward.STATUS_UP {
			ok++
		}
	}

	score = 1.0 / float64(len(states)) * ok

	return
}

func CountDefenceResult(db *sql.DB, round, team int,
	services []steward.Service) (defence float64, err error) {

	defence = 0

	per_service := 1.0 / float64(len(services))

	for _, svc := range services {

		score, err := CountStatesResult(db, round, team, svc)
		if err != nil {
			return defence, err
		}

		defence += score * per_service
	}

	return
}

func CountRound(db *sql.DB, round int, teams []steward.Team,
	services []steward.Service) (err error) {

	round_res := make(map[steward.Team]steward.RoundResult)

	for _, team := range teams {

		res := steward.RoundResult{TeamId: team.Id, Round: round}

		def, err := CountDefenceResult(db, round, team.Id, services)
		if err != nil {
			return err
		}

		res.DefenceScore = def * 2

		round_res[team] = res
	}

	per_service := 1.0 / float64(len(services))

	for _, team := range teams {

		cflags, err := steward.GetCapturedFlags(db, round, team.Id)
		if err != nil {
			return err
		}

		for _, flag := range cflags {

			attacked_team, err := steward.GetTeam(db, flag.TeamId)
			if err != nil {
				return err
			}

			res := round_res[attacked_team]
			res.DefenceScore -= per_service
			if res.DefenceScore < 0 {
				res.DefenceScore = 0
			}
			round_res[attacked_team] = res

			attack_res := round_res[team]
			attack_res.AttackScore += per_service
			round_res[team] = attack_res
		}
	}

	for _, res := range round_res {
		_, err := steward.AddRoundResult(db, res)
		if err != nil {
			return err
		}

	}

	return
}
