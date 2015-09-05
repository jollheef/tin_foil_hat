/**
 * @file service.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date September, 2015
 * @brief queries for service table
 */

package steward

import "database/sql"

type Service struct {
	Name        string
	Port        int
	CheckerPath string
}

func createServiceTable(db *sql.DB) (err error) {

	_, err = db.Exec(`
	CREATE TABLE "service" (
		id	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
		name	TEXT NOT NULL,
		port	INTEGER NOT NULL,
		checker_path	TEXT NOT NULL
	)`)

	return
}

func AddService(db *sql.DB, svc Service) error {

	stmt, err := db.Prepare(
		"INSERT INTO `service` (`name`, `port`, `checker_path`) " +
			"VALUES (?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(svc.Name, svc.Port, svc.CheckerPath)

	if err != nil {
		return err
	}

	return nil
}

func GetServices(db *sql.DB) (services []Service, err error) {

	rows, err := db.Query("SELECT `name`, `port`, `checker_path` " +
		"FROM `service` ")
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var svc Service

		err = rows.Scan(&svc.Name, &svc.Port, &svc.CheckerPath)
		if err != nil {
			return
		}

		services = append(services, svc)
	}

	return
}
