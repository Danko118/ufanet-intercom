package main

import (
	"database/sql"
	"fmt"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var db *sql.DB

func PSQLInit() {
	// Временное решение для удобства разработки
	connStr := "user=postgres dbname=ufanet sslmode=disable password=12345"
	dbc, err := sql.Open("postgres", connStr)

	if err != nil {
		logger.WithFields(logrus.Fields{
			"state":   "Init",
			"status":  "Error",
			"service": "DB",
			"error":   err.Error(),
		}).Fatal("Не удалось подключиться к psql")
	}

	db = dbc
	logger.WithFields(logrus.Fields{
		"state":   "Init",
		"status":  "Success",
		"service": "DB",
	}).Info("Успешно подклченно к psql")
}

func QueryToStruct(query string, args []interface{}, destFunc func(rows *sql.Rows) error) error {
	rows, err := db.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	return destFunc(rows)
}

func FecthEvents(mac string) ([]Event, error) {
	var events []Event
	query := fmt.Sprintf(`SELECT mac, event_name, event_args, event_desc, event_time
	          FROM events 
              WHERE mac='%s'
			  ORDER BY event_time DESC;`, mac)

	err := QueryToStruct(query, nil, func(rows *sql.Rows) error {
		for rows.Next() {
			var evnt Event
			if err := rows.Scan(&evnt.MAC, &evnt.Name, &evnt.Args, &evnt.Desc, &evnt.Time); err != nil {
				return err
			}
			events = append(events, evnt)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return events, nil

}

func FecthIntercom(intercomMAC string) (*Intercom, error) {
	var intercom *Intercom
	query := fmt.Sprintf(`SELECT i.mac,i.address,i.vendor,i.aparts,s.status
	          FROM intercoms i 
			  JOIN intercomStatus s ON i.mac = s.mac
              WHERE i.mac='%s'
			  ORDER BY s.time DESC LIMIT 1;`, intercomMAC)

	err := QueryToStruct(query, nil, func(rows *sql.Rows) error {
		for rows.Next() {
			var intr Intercom
			arr := pq.Int64Array{}
			if err := rows.Scan(&intr.MAC, &intr.Address, &intr.Vendor, &arr, &intr.Status); err != nil {
				return err
			}
			intr.Aparts = arr
			intercom = &intr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return intercom, nil

}

func FecthIntercoms() ([]Intercom, error) {
	var intercoms []Intercom
	query := `SELECT DISTINCT ON (i.mac)
       i.mac,
       i.address,
	   i.vendor,
       s.status
FROM intercoms i
JOIN intercomStatus s ON i.mac = s.mac
ORDER BY i.mac, s.time DESC;`

	err := QueryToStruct(query, nil, func(rows *sql.Rows) error {
		for rows.Next() {
			var intr Intercom
			if err := rows.Scan(&intr.MAC, &intr.Address, &intr.Vendor, &intr.Status); err != nil {
				return err
			}
			intercoms = append(intercoms, intr)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return intercoms, nil

}

func IntercomAppend(msg mqtt.Message) error {
	var intrcom *Intercom
	var err error

	intrcom, err = UnmarshallIntercom(string(msg.Payload()))
	if err != nil {
		return err
	}

	intrcom.MAC = string(strings.Split(msg.Topic(), "/")[1])
	query := `
        INSERT INTO intercoms (address, aparts, vendor, mac)
        VALUES ($1, $2, $3, $4)
		ON CONFLICT (mac) DO NOTHING
    `

	// Выполняем запрос
	_, errr := db.Exec(query, intrcom.Address, pq.Array(intrcom.Aparts), intrcom.Vendor, intrcom.MAC)
	if errr != nil {
		return errr
	}

	return nil
}

func EventAppend(payload string, mac string) error {
	var event *Event
	var err error

	event, err = UnmarshallEvent(payload)
	if err != nil {
		return err
	}

	query := `
        INSERT INTO events (mac, event_name, event_args, event_desc)
        VALUES ($1, $2, $3, $4)
    `

	// Выполняем запрос
	_, errr := db.Exec(query, mac, event.Name, event.Args, event.Desc)
	if errr != nil {
		return errr
	}

	return nil
}

func AppendState(msg mqtt.Message) error {

	query := `
        INSERT INTO intercomStatus (mac, status)
        VALUES ($1, $2)
    `

	// Выполняем запрос
	_, err := db.Exec(query, string(strings.Split(msg.Topic(), "/")[1]), string(msg.Payload()))
	if err != nil {
		return err
	}

	return nil
}
