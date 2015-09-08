package log

import (
	"github.com/wayt/happyngine/env"
	golog "log"
	"os"
)

var debug = golog.New(os.Stdout, "DEBUG ", golog.LstdFlags)
var info = golog.New(os.Stdout, "INFO ", golog.LstdFlags)
var warning = golog.New(os.Stdout, "\033[33mWARNING ", golog.LstdFlags)
var err = golog.New(os.Stdout, "\033[41mERROR ", golog.LstdFlags)
var critical = golog.New(os.Stdout, "\033[41mCRITICAL ", golog.LstdFlags)

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
