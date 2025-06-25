package main

import (
	"errors"
	"fmt"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

var mqttClient mqtt.Client

func MqttInit() {
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883").SetClientID("mqtt-to-ws")
	mqttc := mqtt.NewClient(opts)

	if token := mqttc.Connect(); token.Wait() && token.Error() != nil {
		logger.WithFields(logrus.Fields{
			"state":   "Init",
			"status":  "Error",
			"service": "Mqtt-client",
			"error":   token.Error(),
		}).Fatal("Не удалось подключиться к mqtt")
	}
	mqttClient = mqttc
	logger.WithFields(logrus.Fields{
		"state":   "Init",
		"status":  "Success",
		"service": "Mqtt-client",
	}).Info("Успешно подклченно к MQTT")
}

func TopicResolve(msg mqtt.Message) {
	var mqttMessage = strings.Split(string(msg.Topic()), "/")

	switch mqttMessage[2] {
	case "init":
		err := IntercomAppend(msg)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"state":   "Unmarhall-JSON",
				"status":  "Error",
				"service": "Mqtt-client",
				"error":   err.Error(),
			}).Error("Не удалось обработать mqtt сообщение")
		} else {
			logger.WithFields(logrus.Fields{
				"state":   "Message",
				"status":  "Success",
				"service": "Mqtt-client",
			}).Info("Сообщение успешно обработано")
		}
	case "status":
		err := AppendState(msg)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"state":   "Unmarhall-JSON",
				"status":  "Error",
				"service": "Mqtt-client",
				"error":   err.Error(),
			}).Error("Не удалось обработать mqtt сообщение")
		} else {
			logger.WithFields(logrus.Fields{
				"state":   "Message",
				"status":  "Success",
				"service": "Mqtt-client",
			}).Info("Сообщение успешно обработано")
		}
	case "events":
		err := EventResolve(mqttMessage[3], string(msg.Payload()), mqttMessage[1])
		if err != nil {
			logger.WithFields(logrus.Fields{
				"state":   "Unmarhall-JSON",
				"status":  "Error",
				"service": "Mqtt-client",
				"error":   err.Error(),
			}).Error("Не удалось обработать mqtt сообщение")
		} else {
			logger.WithFields(logrus.Fields{
				"state":   "Message",
				"status":  "Success",
				"service": "Mqtt-client",
			}).Info("Ивент успешно обработан")
		}
	default:
		logger.WithFields(logrus.Fields{
			"state":   "Message",
			"status":  "Error",
			"service": "Mqtt-client",
		}).Error("Не удалось обработать mqtt сообщение, неизвестный топик третьего уровня")
	}

}

func OpenEvent(mac string) {
	var topic = fmt.Sprintf("intercoms/%s/control/door", mac)
	token := mqttClient.Publish(topic, 0, false, "1")
	token.Wait()
}

func CallEvent(mac string, id string) {
	var topic = fmt.Sprintf("intercoms/%s/control/call", mac)
	token := mqttClient.Publish(topic, 0, false, id)
	token.Wait()
}

func EventResolve(event string, payload string, mac string) error {

	switch event {
	case "door":
		err := EventAppend(payload, mac)
		if err != nil {
			return err
		}
	case "call":
		err := EventAppend(payload, mac)
		if err != nil {
			return err
		}
	default:
		return errors.New("неизвестный топик четвертого уровня")
	}
	return nil
}
