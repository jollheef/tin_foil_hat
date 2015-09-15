/**
 * @file raw_commands.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief functions for run checkers
 *
 * Provide functions for call checker executables
 */

package checker

import (
	"fmt"
	"os/exec"
	"strings"
)

import "tinfoilhat/steward"

const timeout string = "10s" // max checker work time

func exit_status(no int) string {
	return fmt.Sprintf("exit status %d", no)
}

func parseState(err error) (steward.ServiceState, error) {

	if err == nil {
		return steward.STATUS_OK, nil
	}

	switch err.Error() {
	case exit_status(124): // returns by timeout
		return steward.STATUS_DOWN, nil
	case exit_status(1):
		return steward.STATUS_ERROR, nil
	case exit_status(2):
		return steward.STATUS_MUMBLE, nil
	case exit_status(3):
		return steward.STATUS_CORRUPT, nil
	case exit_status(4):
		return steward.STATUS_DOWN, nil
	}

	return steward.STATUS_UNKNOWN, err
}

func put(checker string, ip string, port int, flag string) (cred string,
	state steward.ServiceState, err error) {

	cmd := exec.Command("timeout", timeout, checker, "put", ip,
		fmt.Sprintf("%d", port), flag)

	raw_cred, err := cmd.Output()

	state, err = parseState(err)

	cred = strings.Trim(string(raw_cred), " \n")

	return
}

func get(checker string, ip string, port int, cred string) (flag string,
	state steward.ServiceState, err error) {

	cmd := exec.Command("timeout", timeout, checker, "get", ip,
		fmt.Sprintf("%d", port), cred)

	raw_flag, err := cmd.Output()

	state, err = parseState(err)

	flag = strings.Trim(string(raw_flag), " \n")

	return
}

func check(checker string, ip string, port int) (state steward.ServiceState,
	err error) {

	cmd := exec.Command("timeout", timeout, checker, "chk", ip,
		fmt.Sprintf("%d", port))

	_, err = cmd.Output()

	state, err = parseState(err)

	return
}
