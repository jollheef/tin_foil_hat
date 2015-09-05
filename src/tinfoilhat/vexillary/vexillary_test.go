/**
 * @file vexillary_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief test vexillary package
 */

package vexillary_test

import (
	"log"
	"testing"
)

import "tinfoilhat/vexillary"

func TestGenerateKey(t *testing.T) {

	_, err := vexillary.GenerateKey()
	if err != nil {
		log.Fatalln("Generate key error:", err)
	}
}

func TestGenerateFlag(t *testing.T) {

	priv, _ := vexillary.GenerateKey()

	_, err := vexillary.GenerateFlag(priv)
	if err != nil {
		log.Fatalln("Generate flag error:", err)
	}
}

func TestValidFlag(t *testing.T) {

	priv, _ := vexillary.GenerateKey()
	flag, _ := vexillary.GenerateFlag(priv)

	// Check validation of valid flag
	valid, err := vexillary.ValidFlag(flag, priv.PublicKey)
	if !valid {
		log.Fatalln("Valid flag is invalid:", err)
	}

	// Check validation of invalid flag
	invalid_flag := "aaaaaaa6a0993562af00d027aff63e9502754018="
	valid, err = vexillary.ValidFlag(invalid_flag, priv.PublicKey)
	if valid {
		log.Fatalln("Invalid flag is valid:", err)
	}
}
