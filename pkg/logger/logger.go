package logger

import (
	"io"
	logging "log"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	Log *logrus.Logger // share will all packages
)

func init() {
	// the file needs to exist prior
	f, err := os.OpenFile("_log/application.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		// use go's logger, while we configure Logrus
		logging.Fatalf("error opening file: %v", err)
	}

	// configure Logrus
	Log = logrus.New()
	Log.Formatter = &logrus.JSONFormatter{}
	Log.SetReportCaller(true)

	mw := io.MultiWriter(os.Stdout, f)
	Log.SetOutput(mw)
}
