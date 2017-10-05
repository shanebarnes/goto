package logger

import (
    "io"
    "log"
    "os"
    "sync"
)

var flags int = log.Ldate | log.Ltime | log.Lmicroseconds
var logger *log.Logger = nil
var once sync.Once
var writer io.Writer = os.Stdout

func getLogger() *log.Logger {
    once.Do(func() {
        logger = log.New(writer, "", flags)
    })
    return logger
}

func Init(w io.Writer, f int) {
    writer = w
    flags = f
}

func DebugLn(msg string) {
    getLogger().Println("[DBG] " + msg)
}

func ErrorLn(msg string) {
    getLogger().Println("[ERR] " + msg)
}

func InfoLn(msg string) {
    getLogger().Println("[INF] " + msg)
}
