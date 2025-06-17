package main

import (
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
