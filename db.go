package main

import (
	"database/sql"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
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
    `

	// Выполняем запрос
	_, errr := db.Exec(query, intrcom.Address, pq.Array(intrcom.Aparts), intrcom.Vendor, intrcom.MAC)
	if errr != nil {
		return errr
	}

	return nil
}
