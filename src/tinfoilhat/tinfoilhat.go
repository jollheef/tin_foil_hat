/**
 * @file tinfoilhat.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief contest checking system daemon
 *
 * Entry point for contest checking system daemon
 */

package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"time"
)

import (
	"tinfoilhat/config"
	"tinfoilhat/pulse"
	"tinfoilhat/receiver"
	"tinfoilhat/scoreboard"
	"tinfoilhat/steward"
	"tinfoilhat/vexillary"
)

var (
	config_path = kingpin.Arg("config",
		"Path to configuration file.").Required().String()

	db_reinit = kingpin.Flag("reinit", "Reinit database.").Bool()
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

	if *db_reinit {

		log.Println("Reinit database")

		log.Println("Clean database")

		steward.CleanDatabase(db)

		for _, team := range config.Teams {

			log.Println("Add team", team.Name)

			_, err = steward.AddTeam(db, team)
			if err != nil {
				log.Fatalln("Add team failed:", err)
			}
		}

		for _, svc := range config.Services {

			log.Println("Add service", svc.Name)

			err = steward.AddService(db, svc)
			if err != nil {
				log.Fatalln("Add service failed:", err)
			}
		}
	}

	priv, err := vexillary.GenerateKey()
	if err != nil {
		log.Fatalln("Generate key fail:", err)
	}

	go receiver.FlagReceiver(db, priv, config.FlagReceiver.Addr,
		config.FlagReceiver.ReceiveTimeout.Duration)

	go receiver.AdvisoryReceiver(db, config.AdvisoryReceiver.Addr,
		config.AdvisoryReceiver.ReceiveTimeout.Duration)

	go scoreboard.Scoreboard(db, config.Scoreboard.WwwPath,
		config.Scoreboard.Addr,
		config.Scoreboard.UpdateTimeout.Duration,
		config.Pulse.Start.Time,
		config.Pulse.Half.Duration,
		config.Pulse.Lunch.Duration,
		config.Pulse.DarkestTime.Duration)

	err = pulse.Pulse(db, priv,
		config.Pulse.Start.Time,
		config.Pulse.Half.Duration,
		config.Pulse.Lunch.Duration,
		config.Pulse.RoundLen.Duration,
		config.Pulse.CheckTimeout.Duration)
	if err != nil {
		log.Fatalln("Game error:", err)
	}

	log.Println("It's now safe to turn off you computer")
	for {
		time.Sleep(time.Hour)
	}
}
