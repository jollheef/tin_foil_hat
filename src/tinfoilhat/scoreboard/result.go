/**
 * @file result.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief result struct with html conversion
 *
 * Contain structures and html conversion functions
 */

package scoreboard

import "fmt"

import "tinfoilhat/steward"

type TeamResult struct {
	Rank            int
	Name            string
	Score           float64
	ScorePercent    float64
	Attack          float64
	AttackPercent   float64
	Defence         float64
	DefencePercent  float64
	Advisory        int
	AdvisoryPercent float64
	Status          []steward.ServiceState
}

func td(s string, best bool) string {
	if best {
		return `<td bgcolor="#00AAAA"><font color="#FFFFFF">` +
			s + "</td>"
	} else {
		return `<td>` + s + `</td>`
	}
}

func (tr TeamResult) ToHTML() string {

	var status string
	for _, s := range tr.Status {

		var label string

		switch s {
		case steward.STATUS_UP:
			label = "success"
		case steward.STATUS_MUMBLE:
		case steward.STATUS_CORRUPT:
			label = "warning"
		case steward.STATUS_UNKNOWN:
			label = "default"
		default:
			label = "danger"
		}

		status += fmt.Sprintf(
			`<td width="10%%"><span class="label label-%s">%s</span></td>`,
			label, s.String())
	}

	var score_best, attack_best, defence_best, advisory_best bool
	if tr.ScorePercent == 100 {
		score_best = true
	}
	if tr.AttackPercent == 100 {
		attack_best = true
	}
	if tr.DefencePercent == 100 {
		defence_best = true
	}
	if tr.AdvisoryPercent == 100 {
		advisory_best = true
	}

	info := fmt.Sprintf("<td>%d</td><td>%s</td>", tr.Rank, tr.Name)
	score := td(fmt.Sprintf("%05.2f&#37", tr.ScorePercent), score_best)
	attack := td(fmt.Sprintf("%.3f", tr.Attack), attack_best)
	defence := td(fmt.Sprintf("%.3f", tr.Defence), defence_best)
	advisory := td(fmt.Sprintf("%d", tr.Advisory), advisory_best)

	return "<tr>" + info + score + attack + defence + advisory + status + "</tr>"
}

type ByScore []TeamResult

func (tr ByScore) Len() int           { return len(tr) }
func (tr ByScore) Swap(i, j int)      { tr[i], tr[j] = tr[j], tr[i] }
func (tr ByScore) Less(i, j int) bool { return tr[i].Score > tr[j].Score }

type Result struct {
	Teams    []TeamResult
	Services []string
}

func (r Result) ToHTML() string {

	var services string
	for _, s := range r.Services {
		services += "<th>" + s + "</th>"
	}

	var teams string
	for _, t := range r.Teams {

		need_add := len(r.Services) - len(t.Status)

		for i := 0; i < need_add; i++ {
			t.Status = append(t.Status, steward.STATUS_UNKNOWN)
		}

		teams += t.ToHTML()
	}

	return fmt.Sprintf("<thead><th>#</th><th>Team</th><th>Score</th>"+
		"<th>Attack</th><th>Defence</th><th>Advisory</th>%s"+
		"</thead><tbody>%s</tbody>", services, teams)
}
