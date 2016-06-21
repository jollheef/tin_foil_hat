/**
 * @file service.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date September, 2015
 * @brief queries for service table
 */

package steward

import "database/sql"

// Service contains info about service
type Service struct {
	ID          int
	Name        string
	Port        int
	CheckerPath string
	UDP         bool
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

// AddService add service to database
func AddService(db *sql.DB, svc Service) error {

	stmt, err := db.Prepare(
		"INSERT INTO service (name, port, checker_path, udp) " +
			"VALUES ($1, $2, $3, $4)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(svc.Name, svc.Port, svc.CheckerPath, svc.UDP)

	if err != nil {
		return err
	}

	return nil
}

// GetServices get all services from database
func GetServices(db *sql.DB) (services []Service, err error) {

	rows, err := db.Query("SELECT id,name, port, checker_path, udp " +
		"FROM service ")
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var svc Service

		err = rows.Scan(&svc.ID, &svc.Name, &svc.Port, &svc.CheckerPath,
			&svc.UDP)
		if err != nil {
			return
		}

		services = append(services, svc)
	}

	return
}
