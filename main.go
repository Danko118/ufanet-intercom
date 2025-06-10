package main

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

func main() {
	finalize := LoggerInit()
	logrus.WithFields(logrus.Fields{
		"state":  "Init",
		"status": "Success",
	}).Info("Логгер инициализирован успешно")
	defer finalize()

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
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {

	}
}
