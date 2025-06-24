package main

import (
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
	HttpInit()

	mqttClient.Subscribe("#", 0, func(client mqtt.Client, msg mqtt.Message) {
		switch msg.Topic()[:strings.Index(msg.Topic(), "/")] {
		case "intercoms":
			// получили нужный топик
			logger.WithFields(logrus.Fields{
				"state":    "Msg",
				"status":   "Received",
				"service":  "Mqtt-client",
				"topic":    msg.Topic(),
				"mqtt-msg": string(msg.Payload()),
			}).Info("Получено новое сообщение от Mqtt клиента")
			// попытка обработать
			TopicResolve(msg)
			logger.WithFields(logrus.Fields{
				"state":   "Msg",
				"status":  "Unmarshled",
				"service": "Mqtt-client",
			}).Info("Сообщение успешно обработано")
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

	select {}
}
