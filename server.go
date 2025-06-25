package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func intercomHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mac := vars["mac"]
	var intr *Intercom
	var err error

	intr, err = FecthIntercom(mac)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"state":   "Fetching",
			"status":  "Error",
			"service": "DB",
			"error":   err.Error(),
		}).Fatal("Не удалось получить данные из базы данных")
	}
	var intrJSON string
	intrJSON, err = interfaceToData(intr)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"state":   "Unmarshling",
			"status":  "Error",
			"service": "DB",
			"error":   err.Error(),
		}).Fatal("Не удалось обработать данные")
	}

	fmt.Fprint(w, intrJSON)
	logger.WithFields(logrus.Fields{
		"state":    "Response",
		"status":   "Responded",
		"service":  "Web-server",
		"endpoint": "/intercom",
	}).Info("Ответ отправлен клиенту")
}

func openIntercom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mac := vars["mac"]
	OpenEvent(mac)
	fmt.Fprint(w, "Дверь открыта")
}

func callIntercom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mac := vars["mac"]
	id := vars["id"]
	CallEvent(mac, id)
	fmt.Fprint(w, "Звонок запущен")
}

func eventsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mac := vars["mac"]
	var events []Event
	var err error

	events, err = FecthEvents(mac)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"state":   "Fetching",
			"status":  "Error",
			"service": "DB",
			"error":   err.Error(),
		}).Fatal("Не удалось получить данные из базы данных")
	}
	var eventJSON string
	eventJSON, err = interfaceToData(events)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"state":   "Unmarshling",
			"status":  "Error",
			"service": "DB",
			"error":   err.Error(),
		}).Fatal("Не удалось обработать данные")
	}

	fmt.Fprint(w, eventJSON)
	logger.WithFields(logrus.Fields{
		"state":    "Response",
		"status":   "Responded",
		"service":  "Web-server",
		"endpoint": "/intercom",
	}).Info("Ответ отправлен клиенту")
}

func HttpInit() {
	router := mux.NewRouter()
	var err error

	router.HandleFunc("/intercom/{mac}", intercomHandler)
	router.HandleFunc("/intercom/{mac}/open", openIntercom)
	router.HandleFunc("/intercom/{mac}/call/{id}", callIntercom)
	router.HandleFunc("/events/{mac}", eventsHandler)

	http.Handle("/", router)
	http.HandleFunc("/intercoms", func(w http.ResponseWriter, r *http.Request) {
		var fetchedintr []Intercom
		fetchedintr, err = FecthIntercoms()
		if err != nil {
			logger.WithFields(logrus.Fields{
				"state":   "Fetching",
				"status":  "Error",
				"service": "DB",
				"error":   err.Error(),
			}).Fatal("Не удалось получить данные из базы данных")
		}

		var intr string
		intr, err = interfaceToData(fetchedintr)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"state":   "Unmarshling",
				"status":  "Error",
				"service": "DB",
				"error":   err.Error(),
			}).Fatal("Не удалось обработать данные")
		}
		fmt.Fprintf(w, intr)
		logger.WithFields(logrus.Fields{
			"state":    "Response",
			"status":   "Responded",
			"service":  "Web-server",
			"endpoint": "/intercoms",
		}).Info("Ответ отправлен клиенту")
	})

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		logger.WithFields(logrus.Fields{
			"state":   "Init",
			"status":  "Error",
			"service": "Web-server",
			"error":   err.Error(),
		}).Fatal("Не удалось занять порт 8080")
		return
	}

	go func() {
		err := http.Serve(listener, nil)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"state":   "Runtime",
				"status":  "Error",
				"service": "Web-server",
				"error":   err.Error(),
			}).Error("Сервер остановлен с ошибкой")
		}
	}()

	logger.WithFields(logrus.Fields{
		"state":   "Init",
		"status":  "Success",
		"service": "Web-server",
	}).Info("Web-сервер запущен на :8080")

}
