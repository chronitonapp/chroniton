package utils

import (
	stdlog "log"
	"os"

	"github.com/op/go-logging"
)

var (
	Log = logging.MustGetLogger("CHRONITON")
)

func init() {
	// logger
	logBackend := logging.NewLogBackend(os.Stderr, "[CHRONITON]: ", stdlog.LstdFlags)
	logBackend.Color = true

	syslogBackend, err := logging.NewSyslogBackend("[CHRONITON]: ")
	if err != nil {
		Log.Fatal(err)
	}

	logging.SetBackend(logBackend, syslogBackend)
	logging.SetLevel(logging.DEBUG, "CHRONITON")
}
