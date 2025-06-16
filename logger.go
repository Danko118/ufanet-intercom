package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	logger  = logrus.New()
	logFile *os.File
	fileMux sync.Mutex
)

type fileHook struct{}

func (h *fileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *fileHook) Fire(entry *logrus.Entry) error {
	line, err := entry.Logger.Formatter.Format(entry)
	if err != nil {
		return err
	}

	fileMux.Lock()
	defer fileMux.Unlock()

	if logFile != nil {
		if _, err := logFile.Write(line); err != nil {
			return err
		}
	}
	return nil
}

func LoggerInit() func() {
	logDir := "./logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, 0755)
		if err != nil {
			panic(fmt.Sprintf("Не удалось создать директорию логов: %v", err))
		}
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := "./logs/" + timestamp + ".log"

	var err error
	logFile, err = os.Create(filename)
	if err != nil {
		panic(err)
	}

	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	logger.SetOutput(os.Stdout)
	logger.AddHook(&fileHook{})

	logger.Info("Логгер инициализирован")

	return func() {
		if logFile != nil {
			logFile.Close()
		}
	}
}
