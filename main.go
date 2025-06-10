package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/sirupsen/logrus"
)

func main() {
	defer LoggerInit()()

	// defer db.Close()

	// mqttClient.Subscribe("#", 0, func(client mqtt.Client, msg mqtt.Message) {
	// 	log.Print("[Информация] {Новое сообщение} MQTT")
	// 	switch msg.Topic()[:strings.Index(msg.Topic(), "/")] {
	// 	case "Devices":
	// 		if err := ProcessMQTTDeviceData(strings.Split(msg.Topic(), "/")[1], string(msg.Payload())); err == nil {
	// 			log.Print("[Успешно] {Сообщение обработано} MQTT")
	// 		} else {
	// 			log.Printf("[Ошибка] {Этап: Обработка сообщения с MQTT} MQTT: %s", err)
	// 		}
	// 	case "Values":
	// 		err := InsertValue(strings.Split(msg.Topic(), "/")[1], string(msg.Payload()))
	// 		if err != nil {
	// 			log.Fatalf("[Ошибка] {Этап: Получение данных из DB}: %v", err)
	// 		} else {
	// 			log.Print("[Успешно] {Сообщение обработано} MQTT")
	// 		}
	// 	default:
	// 		log.Fatalf("[Ошибка] {Неверное сообщение MQTT} MQTT")

	// 	}
	// 	sendMessageToAllClients()
	// })

	// WebSocket обработчик
	// http.HandleFunc("/", websocketConnect)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, Go!")
		logger.WithFields(logrus.Fields{
			"state":    "Response",
			"status":   "Responded",
			"service":  "Web-server",
			"endpoint": "/",
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

	select {}
}
