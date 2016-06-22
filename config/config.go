/**
 * @file config.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date September, 2015
 * @brief read configuration
 *
 * Contain functions for read configuration file
 */

package config

import (
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/naoina/toml"
)

import "github.com/jollheef/tin_foil_hat/steward"

// Pulse config
type Pulse struct {
	Start        Time
	Half         Duration
	Lunch        Duration
	RoundLen     Duration
	CheckTimeout Duration
	DarkestTime  Duration
}

// FlagReceiver config
type FlagReceiver struct {
	Addr           string
	ReceiveTimeout Duration
	SocketTimeout  Duration
}

// AdvisoryReceiver config
type AdvisoryReceiver struct {
	Addr           string
	ReceiveTimeout Duration
	SocketTimeout  Duration
	Disabled       bool
}

// Config config
type Config struct {
	LogFile        string
	CheckerTimeout Duration
	Database       struct {
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

// Duration type with toml unmarshalling support
type Duration struct {
	time.Duration
}

// UnmarshalTOML for Duration
func (d *Duration) UnmarshalTOML(data []byte) (err error) {
	duration := strings.Replace(string(data), "\"", "", -1)
	d.Duration, err = time.ParseDuration(duration)
	return
}

// Time type with toml unmarshalling support
type Time struct {
	time.Time
}

// UnmarshalTOML for Time
func (t *Time) UnmarshalTOML(data []byte) (err error) {

	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return
	}

	rawTime := strings.Replace(string(data), "\"", "", -1)

	layout := "Jan _2 15:04 2006"
	t.Time, err = time.ParseInLocation(layout, rawTime, loc)
	if err != nil {
		return
	}

	return
}

// ReadConfig for tfh
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
