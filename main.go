/**
 * @file main.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date September, 2015
 * @brief contest checking system daemon
 *
 * Entry point for contest checking system daemon
 */

package main

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"time"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

import (
	"github.com/jollheef/tin_foil_hat/checker"
	"github.com/jollheef/tin_foil_hat/config"
	"github.com/jollheef/tin_foil_hat/pulse"
	"github.com/jollheef/tin_foil_hat/receiver"
	"github.com/jollheef/tin_foil_hat/scoreboard"
	"github.com/jollheef/tin_foil_hat/steward"
	"github.com/jollheef/tin_foil_hat/vexillary"
)

var (
	configPath = kingpin.Arg("config",
		"Path to configuration file.").Required().String()

	dbReinit = kingpin.Flag("reinit", "Reinit database.").Bool()
)

var (
	commitID  string
	buildDate string
	buildTime string
)

func buildInfo() (str string) {

	if len(commitID) > 7 {
		commitID = commitID[:7] // abbreviated commit hash
	}

	str = fmt.Sprintf("Version: tin_foil_hat %s %s %s\n",
		commitID, buildDate, buildTime)
	str += "Author: Mikhail Klementyev <jollheef@riseup.net>\n"
	return
}

func main() {

	fmt.Println(buildInfo())

	kingpin.Parse()

	config, err := config.ReadConfig(*configPath)
	if err != nil {
		log.Fatalln("Cannot open config:", err)
	}

	logFile, err := os.OpenFile(config.LogFile,
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Cannot open file:", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	log.Println(buildInfo())

	var rlim syscall.Rlimit
	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlim)
	if err != nil {
		log.Fatalln("Getrlimit fail:", err)
	}

	log.Println("RLIMIT_NOFILE CUR:", rlim.Cur, "MAX:", rlim.Max)

	db, err := steward.OpenDatabase(config.Database.Connection)
	if err != nil {
		log.Fatalln("Open database fail:", err)
	}

	defer db.Close()

	db.SetMaxOpenConns(config.Database.MaxConnections)

	if *dbReinit {

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

			var network string
			if svc.UDP {
				network = "udp"
			} else {
				network = "tcp"
			}

			log.Printf("Add service %s (%s)\n", svc.Name, network)

			err = steward.AddService(db, svc)
			if err != nil {
				log.Fatalln("Add service failed:", err)
			}
		}
	}

	checker.SetTimeout(config.CheckerTimeout.Duration)

	priv, err := vexillary.GenerateKey()
	if err != nil {
		log.Fatalln("Generate key fail:", err)
	}

	go receiver.FlagReceiver(db, priv, config.FlagReceiver.Addr,
		config.FlagReceiver.ReceiveTimeout.Duration,
		config.FlagReceiver.SocketTimeout.Duration)

	go receiver.AdvisoryReceiver(db, config.AdvisoryReceiver.Addr,
		config.AdvisoryReceiver.ReceiveTimeout.Duration,
		config.AdvisoryReceiver.SocketTimeout.Duration)

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
