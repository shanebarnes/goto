package logger

import (
    "log"
    "os"
    "sync"
)

var writer io.Writer = os.Stdout
var logger *log.Logger = nil
var once sync.Once

func getLoggerInstance() *log.Logger {
    once.Do(func() {
        flags := log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile
        file, err := os.OpenFile("scout.log", os.O_APPEND | os.O_CREATE | os.O_RDWR, 0644)

        if err == nil {
            logger = log.New(file, "", flags)
        } else {
            logger = log.New(writer, "", flags)
        }
        //defer file.Close()
    })
    return logger
}

func LoggerDebug(msg string) {
    getLogInstance().Println("[DEBUG] " + msg)
}

func LoggerError(msg string) {
    getLogInstance().Println("[ERROR] " + msg)
}

func LoggerInfo(msg string) {
    getLogInstance().Println("[INFO] " + msg)
}
