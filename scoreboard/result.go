/**
 * @file result.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date September, 2015
 * @brief result struct with html conversion
 *
 * Contain structures and html conversion functions
 */

package scoreboard

import "fmt"

import "github.com/jollheef/tin_foil_hat/steward"

// TeamResult contain info for scoreboard
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
	}

	return `<td>` + s + `</td>`
}

// ToHTML convert TeamResult to HTML
func (tr TeamResult) ToHTML(hideScore bool) string {

	var status string
	for _, s := range tr.Status {

		var label string

		switch s {
		case steward.StatusUP:
			label = "success"
		case steward.StatusMumble:
		case steward.StatusCorrupt:
			label = "warning"
		case steward.StatusUnknown:
			label = "default"
		default:
			label = "important"
		}

		status += fmt.Sprintf(
			`<td width="10%%"><span class="label label-%s">%s</span></td>`,
			label, s.String())
	}

	var scoreBest, attackBest, defenceBest, advisoryBest bool
	if tr.ScorePercent == 100 {
		scoreBest = true
	}
	if tr.AttackPercent == 100 {
		attackBest = true
	}
	if tr.DefencePercent == 100 {
		defenceBest = true
	}
	if tr.AdvisoryPercent == 100 {
		advisoryBest = true
	}

	var info, score, attack, defence, advisory string

	if hideScore {
		hidden := `<td>&#xFFFD</td>`
		info = hidden + fmt.Sprintf("<td>%s</td>", tr.Name)
		score = hidden
		attack = hidden
		defence = hidden
		defence = hidden
		advisory = hidden
	} else {
		info = fmt.Sprintf("<td>%d</td><td>%s</td>", tr.Rank, tr.Name)
		score = td(fmt.Sprintf("%05.2f&#37", tr.ScorePercent), scoreBest)
		attack = td(fmt.Sprintf("%.3f", tr.Attack), attackBest)
		defence = td(fmt.Sprintf("%.3f", tr.Defence), defenceBest)
		advisory = td(fmt.Sprintf("%d", tr.Advisory), advisoryBest)
	}

	return "<tr>" + info + score + attack + defence + advisory + status + "</tr>"
}

// ByScore sort team result by score
type ByScore []TeamResult

func (tr ByScore) Len() int           { return len(tr) }
func (tr ByScore) Swap(i, j int)      { tr[i], tr[j] = tr[j], tr[i] }
func (tr ByScore) Less(i, j int) bool { return tr[i].Score > tr[j].Score }

// Result contain result and services for scoreboard
type Result struct {
	Teams    []TeamResult
	Services []string
}

// ToHTML convert Result to HTML
func (r Result) ToHTML(hideScore bool) string {

	var services string
	for _, s := range r.Services {
		services += "<th>" + s + "</th>"
	}

	var teams string
	for _, t := range r.Teams {

		needAdd := len(r.Services) - len(t.Status)

		for i := 0; i < needAdd; i++ {
			t.Status = append(t.Status, steward.StatusUnknown)
		}

		teams += t.ToHTML(hideScore)
	}

	return fmt.Sprintf("<thead><th>#</th><th>Team</th><th>Score</th>"+
		"<th>Attack</th><th>Defence</th><th>Advisory</th>%s"+
		"</thead><tbody>%s</tbody>", services, teams)
}
