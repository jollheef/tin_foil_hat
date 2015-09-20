/**
 * @file advisory_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief test work with advisory table
 */

package steward_test

import (
	"fmt"
	"log"
	"testing"
	"time"
)

import "tinfoilhat/steward"

func TestAddAdvisory(t *testing.T) {

	db, err := openDB()
	if err != nil {
		log.Fatalln("Open database failed:", err)
	}

	defer db.Close()

	_, err = steward.AddAdvisory(db.db, 10, "ololo")
	if err != nil {
		log.Fatalln("Add advisory failed:", err)
	}
}

func TestReviewAdvisory(t *testing.T) {

	db, err := openDB()

	defer db.Close()

	team_id := 10
	advisory_text := "advisory text"

	id, _ := steward.AddAdvisory(db.db, team_id, advisory_text)

	err = steward.ReviewAdvisory(db.db, id, 20)
	if err != nil {
		log.Fatalln("Review advisory fail:", err)
	}
}

func TestGetAdvisoryScore(t *testing.T) {

	db, err := openDB()

	defer db.Close()

	team_id := 10
	advisory_text := "advisory text"
	score := 40
	amount := 5

	var ids []int

	for i := 0; i < amount; i++ {
		entry_text := fmt.Sprintf("%s%d", advisory_text, i)
		id, err := steward.AddAdvisory(db.db, team_id, entry_text)
		if err != nil {
			log.Fatalln("Add advisory failed:", err)
		}
		ids = append(ids, id)
		steward.ReviewAdvisory(db.db, id, score)
	}

	team_score, err := steward.GetAdvisoryScore(db.db, team_id)
	if err != nil {
		log.Fatalln("Get team advisory score fail:", err)
	}

	if team_score != score*amount {
		log.Fatalf("Team advisory score (%d) not equal of sum of "+
			"entries (%d)", team_score, score*amount)
	}
}

func TestGetAdvisories(t *testing.T) {

	db, err := openDB()

	defer db.Close()

	team_id := 10

	var adv1, adv2 steward.Advisory

	adv1.Text = "adv1_test"
	adv1.Score = 10
	adv1.Reviewed = true

	adv2.Text = "adv2_test"
	adv2.Score = 20
	adv2.Reviewed = true

	adv1.Id, _ = steward.AddAdvisory(db.db, team_id, adv1.Text)
	adv2.Id, _ = steward.AddAdvisory(db.db, team_id, adv2.Text)

	steward.ReviewAdvisory(db.db, adv1.Id, adv1.Score)
	steward.ReviewAdvisory(db.db, adv2.Id, adv2.Score)

	advisories, err := steward.GetAdvisories(db.db)
	if err != nil {
		log.Fatalln("Get all advisories failed:", err)
	}

	if len(advisories) != 2 {
		log.Fatalln("Get advisories more than added")
	}

	if time.Now().Sub(advisories[0].Timestamp) > 5*time.Second {
		log.Fatalln("Time must be ~ current:", advisories[0].Timestamp)
	}

	// No timestamp check
	adv1.Timestamp = advisories[0].Timestamp
	adv2.Timestamp = advisories[1].Timestamp

	if advisories[0] != adv1 || advisories[1] != adv2 {
		log.Fatalln("Added advisories broken")
	}

}
