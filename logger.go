package main

import (
	"bytes"
	"encoding/json"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	buffer   bytes.Buffer
	logFile  *os.File
	filePath string
)

type bufferHook struct{}

func (h *bufferHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *bufferHook) Fire(entry *logrus.Entry) error {
	line, err := entry.Logger.Formatter.Format(entry)
	if err != nil {
		return err
	}
	buffer.Write(line)
	return nil
}

func LoggerInit() func() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	logrus.SetOutput(os.Stdout)

	// Добавляем хук в логрус
	logrus.AddHook(&bufferHook{})

	// Возврат функции сохранения JSON-массива
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filePath = "./logs/log_" + timestamp + ".json"

	var err error
	logFile, err = os.Create(filePath)
	if err != nil {
		panic(err)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	logrus.SetOutput(os.Stdout) // stdout для разработчика

	logrus.AddHook(&bufferHook{}) // буферизация

	logrus.Info("Логгер инициализирован")

	return func() {
		defer logFile.Close()

		lines := bytes.Split(buffer.Bytes(), []byte{'\n'})
		var logs []json.RawMessage
		for _, line := range lines {
			if len(line) == 0 {
				continue
			}
			logs = append(logs, json.RawMessage(line))
		}

		wrapped := map[string]interface{}{
			"logs": logs,
		}

		enc := json.NewEncoder(logFile)
		enc.SetIndent("", "  ")
		if err := enc.Encode(wrapped); err != nil {
			panic(err)
		}
	}
}
