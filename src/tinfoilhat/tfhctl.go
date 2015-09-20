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
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
)

import (
	"tinfoilhat/config"
	"tinfoilhat/steward"
)

var (
	config_path = kingpin.Flag("config",
		"Path to configuration file.").Required().String()

	adv = kingpin.Command("advisory", "Work with advisories.")

	advList = adv.Command("list", "List advisories.")

	advReview   = adv.Command("review", "Review advisory.")
	advReviewId = advReview.Arg("id", "advisory id").Required().Int()
	advScore    = advReview.Arg("score", "advisory id").Required().Int()
)

func main() {

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
	}
}
