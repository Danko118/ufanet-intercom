package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

func main() {

	defer LoggerInit()()
	MqttInit()
	defer mqttClient.Disconnect(250)
	PSQLInit()
	defer db.Close()

	mqttClient.Subscribe("#", 0, func(client mqtt.Client, msg mqtt.Message) {
		switch msg.Topic()[:strings.Index(msg.Topic(), "/")] {
		case "intercoms":
			logger.WithFields(logrus.Fields{
				"state":    "Msg",
				"status":   "Received",
				"service":  "Mqtt-client",
				"topic":    msg.Topic(),
				"mqtt-msg": string(msg.Payload()),
			}).Info("Получено новое сообщение от Mqtt клиента")
			err := ProcessMQTTData(string(msg.Payload()))
			if err != nil {
				logger.WithFields(logrus.Fields{
					"state":   "Unmarhall-JSON",
					"status":  "Error",
					"service": "Mqtt-client",
					"error":   err.Error(),
				}).Error("Не удалось обработать mqtt сообщение")
			}
		default:
			logger.WithFields(logrus.Fields{
				"state":    "Msg",
				"status":   "Received",
				"service":  "Mqtt-client",
				"topic":    msg.Topic(),
				"mqtt-msg": string(msg.Payload()),
			}).Warn("Получено новое сообщение от Mqtt клиента, но неизвестный топик")

		}
	})

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
