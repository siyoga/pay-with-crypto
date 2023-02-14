package utility

import (
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	loggerFile := &lumberjack.Logger{
		Filename:   "logs/pay-with-crypto.log",
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     28,
	}

	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true, FullTimestamp: true})
	logrus.SetOutput(loggerFile)
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func Error(e error, data string) {
	logrus.WithField("data", data).Error(e)
}
