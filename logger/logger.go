package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

// In go, init functions in a package get called only once.
// So, no singleton initialization is required.

// TODO
// Make this as singleton
func Init(debug bool, logFile string) {
	log = logrus.New()

	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)

	if err != nil {
		log.Fatalf("Unable to open log file: %s", err.Error())
	}
	log.SetOutput(f)

	if debug {
		log.SetLevel(logrus.DebugLevel)
	} else {
		log.SetLevel(logrus.InfoLevel)
	}
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		DisableQuote:  true,
	})
}

func Get() *logrus.Logger {
	return log
}
