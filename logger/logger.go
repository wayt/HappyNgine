package logger

import (
    "time"
    "os"
    "bufio"
    "fmt"
)

type writer struct {

    Out *bufio.Writer
    Flags int
}

const (
    LOG_DEBUG       = 1
    LOG_FATAL       = 2
    LOG_ERROR       = 4
    LOG_INFO        = 8
)

type Logger struct{

    Writers []writer
}

func New() *Logger {

    this := &Logger{}

    return this
}

func (this *Logger) LogToFilename(filename string, flags int) {

    f, err := os.OpenFile(filename, os.O_WRONLY | os.O_CREATE, 0666)
    if err != nil {
        panic(err)
    }

    w := bufio.NewWriter(f)
    this.AddOutput(w, flags)
}

func (this *Logger) LogToFile(file *os.File, flags int) {

    w := bufio.NewWriter(file)
    this.AddOutput(w, flags)
}

func (this *Logger) AddOutput(out *bufio.Writer, flags int) {

    this.Writers = append(this.Writers, writer{out, flags})
}

func (this *Logger) log(flags int, prefix string, msg string) {

    now := time.Now().Format(time.RFC3339)

    for _, w := range this.Writers {

        if (w.Flags & flags) != 0 {

            _, err := w.Out.WriteString("[" + now + "][" + prefix + "]: " + msg + "\n")

            if err != nil {
                fmt.Println("LoggerError: ", err)
            }
            w.Out.Flush()
        }
    }
}

func (this *Logger) Debug(msg string) {

    this.log(LOG_DEBUG, "DEBUG", msg)
}

func (this *Logger) Debugf(msg string, v ...interface{}) {

    this.Debug(fmt.Sprintf(msg, v...))
}

func (this *Logger) Debugln(v ...interface{}) {

    this.Debug(fmt.Sprintln(v...))
}

func (this *Logger) Fatal(msg string) {

    this.log(LOG_FATAL, "FATAL", msg)
}

func (this *Logger) Fatalf(msg string, v ...interface{}) {

    this.Fatal(fmt.Sprintf(msg, v...))
}

func (this *Logger) Fatalln(v ...interface{}) {

    this.Fatal(fmt.Sprintln(v...))
}

func (this *Logger) Error(msg string) {

    this.log(LOG_ERROR, "ERROR", msg)
}

func (this *Logger) Errorf(msg string, v ...interface{}) {

    this.Error(fmt.Sprintf(msg, v...))
}

func (this *Logger) Errorln(v ...interface{}) {

    this.Error(fmt.Sprintln(v...))
}

func (this *Logger) Info(msg string) {

    this.log(LOG_INFO, "INFO", msg)
}

func (this *Logger) Infof(msg string, v ...interface{}) {

    this.Info(fmt.Sprintf(msg, v...))
}

func (this *Logger) Infoln(v ...interface{}) {

    this.Info(fmt.Sprintln(v...))
}


