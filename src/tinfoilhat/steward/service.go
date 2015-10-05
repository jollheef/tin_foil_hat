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
	Id          int
	Name        string
	Port        int
	CheckerPath string
	Udp         bool
}

func createServiceTable(db *sql.DB) (err error) {

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS "service" (
		id	SERIAL PRIMARY KEY,
		name	TEXT NOT NULL,
		port	INTEGER NOT NULL,
		checker_path	TEXT NOT NULL,
		udp	BOOLEAN NOT NULL
	)`)

	return
}

func AddService(db *sql.DB, svc Service) error {

	stmt, err := db.Prepare(
		"INSERT INTO service (name, port, checker_path, udp) " +
			"VALUES ($1, $2, $3, $4)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(svc.Name, svc.Port, svc.CheckerPath, svc.Udp)

	if err != nil {
		return err
	}

	return nil
}

func GetServices(db *sql.DB) (services []Service, err error) {

	rows, err := db.Query("SELECT id,name, port, checker_path, udp " +
		"FROM service ")
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var svc Service

		err = rows.Scan(&svc.Id, &svc.Name, &svc.Port, &svc.CheckerPath,
			&svc.Udp)
		if err != nil {
			return
		}

		services = append(services, svc)
	}

	return
}
