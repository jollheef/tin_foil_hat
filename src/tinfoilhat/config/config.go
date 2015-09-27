/**
 * @file config.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief read configuration
 *
 * Contain functions for read configuration file
 */

package config

import (
	"github.com/naoina/toml"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

import "tinfoilhat/steward"

type Pulse struct {
	Start        Time
	Half         Duration
	Lunch        Duration
	RoundLen     Duration
	CheckTimeout Duration
	DarkestTime  Duration
}

type FlagReceiver struct {
	Addr           string
	ReceiveTimeout Duration
}

type AdvisoryReceiver struct {
	Addr           string
	ReceiveTimeout Duration
}

type Config struct {
	Database struct {
		Connection     string
		MaxConnections int
	}
	Scoreboard struct {
		WwwPath       string
		Addr          string
		UpdateTimeout Duration
	}
	Pulse            Pulse
	FlagReceiver     FlagReceiver
	AdvisoryReceiver AdvisoryReceiver
	Teams            []steward.Team
	Services         []steward.Service
}

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalTOML(data []byte) (err error) {
	duration := strings.Replace(string(data), "\"", "", -1)
	d.Duration, err = time.ParseDuration(duration)
	return
}

type Time struct {
	time.Time
}

func (t *Time) UnmarshalTOML(data []byte) (err error) {

	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return
	}

	raw_time := strings.Replace(string(data), "\"", "", -1)

	layout := "Jan _2 15:04 2006"
	t.Time, err = time.ParseInLocation(layout, raw_time, loc)
	if err != nil {
		return
	}

	return
}

func ReadConfig(path string) (cfg Config, err error) {

	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

	err = toml.Unmarshal(buf, &cfg)
	if err != nil {
		return
	}

	return
}
