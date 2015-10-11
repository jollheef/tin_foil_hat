/**
 * @file tfhctl.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief contest checking system CLI
 *
 * Entry point for contest checking system CLI
 */

package main

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
)

import (
	"github.com/jollheef/tin_foil_hat/config"
	"github.com/jollheef/tin_foil_hat/scoreboard"
	"github.com/jollheef/tin_foil_hat/steward"
)

var (
	config_path = kingpin.Flag("config",
		"Path to configuration file.").Required().String()

	score = kingpin.Command("scoreboard", "View scoreboard.")

	adv = kingpin.Command("advisory", "Work with advisories.")

	advList        = adv.Command("list", "List advisories.")
	advNotReviewed = adv.Flag("not-reviewed",
		"List only not reviewed advisory.").Bool()

	advReview   = adv.Command("review", "Review advisory.")
	advReviewId = advReview.Arg("id", "advisory id").Required().Int()
	advScore    = advReview.Arg("score", "advisory id").Required().Int()

	advHide   = adv.Command("hide", "Hide advisory.")
	advHideId = advHide.Arg("id", "advisory id").Required().Int()

	advUnhide   = adv.Command("unhide", "Unhide advisory.")
	advUnhideId = advUnhide.Arg("id", "advisory id").Required().Int()
)

var (
	COMMIT_ID  string
	BUILD_DATE string
	BUILD_TIME string
)

func buildInfo() (str string) {
	str = fmt.Sprintf("Version: tin_foil_hat %s %s %s\n",
		COMMIT_ID[:7], BUILD_DATE, BUILD_TIME)
	str += "Author: Mikhail Klementyev <jollheef@riseup.net>\n"
	return
}

func main() {

	fmt.Println(buildInfo())

	kingpin.Parse()

	config, err := config.ReadConfig(*config_path)
	if err != nil {
		log.Fatalln("Cannot open config:", err)
	}

	db, err := steward.OpenDatabase(config.Database.Connection)
	if err != nil {
		log.Fatalln("Open database fail:", err)
	}

	defer db.Close()

	db.SetMaxOpenConns(config.Database.MaxConnections)

	switch kingpin.Parse() {
	case "advisory list":
		advisories, err := steward.GetAdvisories(db)
		if err != nil {
			log.Fatalln("Get advisories fail:", err)
		}

		for _, advisory := range advisories {

			if *advNotReviewed && advisory.Reviewed {
				continue
			}

			fmt.Printf(">>> Advisory: id %d <<<\n", advisory.Id)
			fmt.Printf("(Score: %d, Reviewed: %t, Timestamp: %s)\n",
				advisory.Score, advisory.Reviewed,
				advisory.Timestamp.String())

			fmt.Println(advisory.Text)
		}

	case "advisory review":
		err := steward.ReviewAdvisory(db, *advReviewId, *advScore)
		if err != nil {
			log.Fatalln("Advisory review fail:", err)
		}

	case "advisory hide":
		err := steward.HideAdvisory(db, *advHideId, true)
		if err != nil {
			log.Fatalln("Advisory hide fail:", err)
		}

	case "advisory unhide":
		err := steward.HideAdvisory(db, *advUnhideId, false)
		if err != nil {
			log.Fatalln("Advisory unhide fail:", err)
		}

	case "scoreboard":
		res, err := scoreboard.CollectLastResult(db)
		if err != nil {
			log.Fatalln("Get last result fail:", err)
		}

		scoreboard.CountScoreAndSort(&res)

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Rank", "Name", "Score", "Attack",
			"Defence", "Advisory"})

		for _, tr := range res.Teams {

			var row []string

			row = append(row, fmt.Sprintf("%d", tr.Rank))
			row = append(row, tr.Name)
			row = append(row, fmt.Sprintf("%05.2f%%", tr.ScorePercent))
			row = append(row, fmt.Sprintf("%.3f", tr.Attack))
			row = append(row, fmt.Sprintf("%.3f", tr.Defence))
			row = append(row, fmt.Sprintf("%d", tr.Advisory))

			table.Append(row)
		}

		table.Render()
	}
}
