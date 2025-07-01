package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func withCORS(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Устанавливаем CORS-заголовки
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Если это preflight-запрос, просто отвечаем 200 OK
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Иначе продолжаем выполнение
		handler(w, r)
	}
}

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

func GETIntercoms(w http.ResponseWriter, r *http.Request) {
	var err error
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
}

func HttpInit() {
	router := mux.NewRouter()
	var err error

	router.HandleFunc("/intercoms", withCORS(GETIntercoms))
	router.HandleFunc("/intercom/{mac}", withCORS(intercomHandler))
	router.HandleFunc("/intercom/{mac}/open", withCORS(openIntercom))
	router.HandleFunc("/intercom/{mac}/call/{id}", withCORS(callIntercom))
	router.HandleFunc("/events/{mac}", withCORS(eventsHandler))

	http.Handle("/", router)

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
