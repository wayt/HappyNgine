package log

import (
	"github.com/wayt/happyngine/env"
	golog "log"
	"os"
)

var (
	debug    *golog.Logger
	info     *golog.Logger
	warning  *golog.Logger
	err      *golog.Logger
	critical *golog.Logger
)

func init() {

	logger := os.Stdout
	logFile := env.Get("LOG_FILE")
	if logFile != "" {
		f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			panic(err)
		}
		// defer f.Close()

		logger = f
	}

	debug = golog.New(logger, "DEBUG ", golog.LstdFlags)
	info = golog.New(logger, "INFO ", golog.LstdFlags)
	warning = golog.New(logger, "\033[33mWARNING ", golog.LstdFlags)
	err = golog.New(logger, "\033[41mERROR ", golog.LstdFlags)
	critical = golog.New(logger, "\033[41mCRITICAL ", golog.LstdFlags)
}

func Debugln(args ...interface{}) {
	if env.Get("DEBUG") == "1" {
		debug.Println(args...)
	}
}
func Infoln(args ...interface{}) {
	info.Println(args...)
}

func Warningln(args ...interface{}) {
	warning.Println(append(args, "\x1B[0m")...)
}

func Errorln(args ...interface{}) {
	err.Println(append(args, "\x1B[0m")...)
}

func Criticalln(args ...interface{}) {
	critical.Println(append(args, "\x1B[0m")...)
}
