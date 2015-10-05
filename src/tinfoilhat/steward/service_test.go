/**
 * @file service_test.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief test work with service table
 */

package steward_test

import (
	"log"
	"testing"
)

import "tinfoilhat/steward"

func TestAddService(t *testing.T) {

	db, err := openDB()

	defer db.Close()

	svc := steward.Service{-1, "lol", 10, "/test", false}

	err = steward.AddService(db.db, svc)
	if err != nil {
		log.Fatalln("Add service fail:", err)
	}
}

func TestGetServices(t *testing.T) {

	db, err := openDB()

	defer db.Close()

	svc := steward.Service{-1, "lol", 10, "/test", false}

	const services_amount int = 5

	for i := 0; i < services_amount; i++ {
		svc.Port = i
		err = steward.AddService(db.db, svc)
	}

	services, err := steward.GetServices(db.db)
	if err != nil {
		log.Fatalln("Get services fail:", err)
	}

	if len(services) != services_amount {
		log.Fatalln("Get services moar than add")
	}

	for i := 0; i < len(services); i++ {
		svc.Id = i + 1
		svc.Port = i
		if services[i] != svc {
			log.Fatalln("Get service", services[i], "instead", svc)
		}
	}
}
