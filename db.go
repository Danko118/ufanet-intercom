package main

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var db *sql.DB

func PSQLInit() {
	connStr := "user=postgres dbname=intercom sslmode=disable"
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
