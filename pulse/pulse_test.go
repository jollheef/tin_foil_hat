/**
 * @file pulse_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date September, 2015
 * @brief test pulse package
 */

package pulse_test

import (
	"log"
	"testing"
	"time"
)

import "github.com/jollheef/tin_foil_hat/pulse"

func TestWaitWithWait(*testing.T) {

	end_time := time.Now().Add(time.Second / 10)

	waited := pulse.Wait(end_time, time.Nanosecond)

	time_diff := time.Now().Sub(end_time)

	if time_diff > time.Second/100 {
		log.Fatalln("Too long wait time diff:", time_diff)
	}

	if !waited {
		log.Fatalln("Fail: no wait")
	}
}

func TestWaitWithoutWait(*testing.T) {

	end_time := time.Now()

	waited := pulse.Wait(end_time, time.Nanosecond)

	time_diff := time.Now().Sub(end_time)

	if time_diff > time.Second/100 {
		log.Fatalln("Too long wait time diff:", time_diff)
	}

	if waited {
		log.Fatalln("Fail: has wait")
	}
}
