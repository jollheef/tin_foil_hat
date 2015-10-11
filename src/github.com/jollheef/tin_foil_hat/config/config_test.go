/**
 * @file config_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief test config package
 */

package config_test

import (
	"log"
	"testing"
)

import "github.com/jollheef/tin_foil_hat/config"

func bug_on_invalid(real, parsed string) {
	if real != parsed {
		log.Fatalln("Parsed", parsed, "instead", real)
	}
}

func TestReadConfig(*testing.T) {

	cfg, err := config.ReadConfig("tinfoilhat.toml")
	if err != nil {
		log.Fatalln("Read config error:", err)
	}

	bug_on_invalid("2015-08-02 15:04:00 +0300 MSK", cfg.Pulse.Start.String())

	bug_on_invalid("4h0m0s", cfg.Pulse.Half.String())

	bug_on_invalid("1h0m0s", cfg.Pulse.Lunch.String())

	bug_on_invalid("2m0s", cfg.Pulse.RoundLen.String())

	bug_on_invalid("30s", cfg.Pulse.CheckTimeout.String())

	// other values has built-in types
}
