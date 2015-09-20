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

func (tr TeamResult) ToHTML() string {

	var status string
	for _, s := range tr.Status {
		status += "<td>" + s.String() + "</td>"
	}

	return fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%05.2f &#37</td><td>%f</td>"+
		"<td>%f</td><td>%d</td>%s</tr>", tr.Rank, tr.Name,
		tr.ScorePercent, tr.Attack, tr.Defence, tr.Advisory, status)
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
		teams += t.ToHTML()
	}

	return fmt.Sprintf("<thead><th>Rank</th><th>Name</th><th>Score</th>"+
		"<th>Attack</th><th>Defence</th><th>Advisory</th>%s"+
		"</thead><tbody>%s</tbody>", services, teams)
}
