/**
 * @file vexillary_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief test vexillary package
 */

package vexillary_test

import (
	"fmt"
	"os"
	"testing"
)

import "tinfoilhat/vexillary"

func TestGenerateKey(t *testing.T) {

	fmt.Println("GenerateFlag test")

	_, err := vexillary.GenerateKey()
	if err != nil {
		fmt.Println("\tGenerate key error:", err)
		os.Exit(1)
	}
}

func TestGenerateFlag(t *testing.T) {

	fmt.Println("GenerateFlag test")

	priv, _ := vexillary.GenerateKey()

	_, err := vexillary.GenerateFlag(priv)
	if err != nil {
		fmt.Println("\tGenerate flag error:", err)
		os.Exit(1)
	}
}

func TestValidFlag(t *testing.T) {

	fmt.Println("ValidFlag test")

	priv, _ := vexillary.GenerateKey()
	flag, _ := vexillary.GenerateFlag(priv)

	// Check validation of valid flag
	valid, err := vexillary.ValidFlag(flag, priv.PublicKey)
	if !valid {
		fmt.Println("\tValid flag is invalid:", err)
		os.Exit(1)
	}

	// Check validation of invalid flag
	invalid_flag := "aaaaaaa6a0993562af00d027aff63e9502754018="
	valid, err = vexillary.ValidFlag(invalid_flag, priv.PublicKey)
	if valid {
		fmt.Println("\tInvalid flag is valid:", err)
		os.Exit(1)
	}
}
