/**
 * @file raw_commands.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date September, 2015
 * @brief functions for run checkers
 *
 * Provide functions for call checker executables
 */

package checker

import (
	"fmt"
	"strings"
	"time"

	system "github.com/jollheef/go-system"
	"github.com/jollheef/tin_foil_hat/steward"
)

var (
	timeout            = "10s" // max checker work time
	connectionAttempts = "2"   // ssh option
	connectTimeout     = "5"   // ssh option
)

// SetTimeout set max checker work time
func SetTimeout(d time.Duration) {
	timeout = fmt.Sprintf("%ds", int(d.Seconds()))
}

func parseState(ret int) steward.ServiceState {

	switch ret {
	case 0:
		return steward.StatusUP
	case 124: // returns by timeout
		return steward.StatusDown
	case 1:
	case 255: // Could not resolve hostname
		return steward.StatusError
	case 2:
		return steward.StatusMumble
	case 3:
		return steward.StatusCorrupt
	case 4:
		return steward.StatusDown
	}

	return steward.StatusUnknown
}

func put(checker, ip string, port int, flag string) (cred, logs string,
	state steward.ServiceState, err error) {

	cred, logs, ret, err := system.System("timeout", timeout, checker, "put", ip,
		fmt.Sprintf("%d", port), flag)

	state = parseState(ret)
	if state != steward.StatusUnknown {
		err = nil
	}

	cred = strings.Trim(cred, " \n")

	return
}

func sshPut(host, checker, ip string, port int, flag string) (cred, logs string,
	state steward.ServiceState, err error) {

	cred, logs, ret, err := system.System("ssh",
		"-o", "ConnectTimeout="+connectTimeout,
		"-o", "ConnectionAttempts="+connectionAttempts,
		host, "timeout", timeout, checker,
		"put", ip, fmt.Sprintf("%d", port), flag)

	state = parseState(ret)
	if state != steward.StatusUnknown {
		err = nil
	}

	cred = strings.Trim(cred, " \n")

	return
}

func get(checker, ip string, port int, cred string) (flag, logs string,
	state steward.ServiceState, err error) {

	flag, logs, ret, err := system.System("timeout", timeout, checker, "get", ip,
		fmt.Sprintf("%d", port), cred)

	state = parseState(ret)
	if state != steward.StatusUnknown {
		err = nil
	}

	flag = strings.Trim(flag, " \n")

	return
}

func sshGet(host, checker, ip string, port int, cred string) (flag, logs string,
	state steward.ServiceState, err error) {

	flag, logs, ret, err := system.System("ssh",
		"-o", "ConnectTimeout="+connectTimeout,
		"-o", "ConnectionAttempts="+connectionAttempts,
		host, "timeout", timeout, checker,
		"get", ip, fmt.Sprintf("%d", port), cred)

	state = parseState(ret)
	if state != steward.StatusUnknown {
		err = nil
	}

	flag = strings.Trim(flag, " \n")

	return
}

func check(checker, ip string, port int) (state steward.ServiceState,
	logs string, err error) {

	_, logs, ret, err := system.System("timeout", timeout, checker, "chk", ip,
		fmt.Sprintf("%d", port))

	state = parseState(ret)
	if state != steward.StatusUnknown {
		err = nil
	}

	return
}

func sshCheck(host, checker, ip string, port int) (state steward.ServiceState,
	logs string, err error) {

	_, logs, ret, err := system.System("ssh",
		"-o", "ConnectTimeout="+connectTimeout,
		"-o", "ConnectionAttempts="+connectionAttempts,
		host, "timeout", timeout, checker,
		"chk", ip, fmt.Sprintf("%d", port))

	state = parseState(ret)
	if state != steward.StatusUnknown {
		err = nil
	}

	return
}
