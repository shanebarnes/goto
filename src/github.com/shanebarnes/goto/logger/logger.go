package logger

import (
    "io"
    "log"
    "os"
    "sync"
    "sync/atomic"
)

type Level int32

const (
    Debug Level = iota
    Info
    Error
)

var flags int = log.Ldate | log.Ltime | log.Lmicroseconds
var level int32 = int32(Info)
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

func SetLevel(l Level) {
    val := int32(l)
    atomic.StoreInt32(&level, val)
}

func GetLevel() Level {
    val := atomic.LoadInt32(&level)
    return Level(val)
}

func Println(lev Level, msg string) {
    if lev >= GetLevel() {
        prefix := ""

        switch lev {
        case Debug:
            prefix = "[DBG]"
        case Info:
            prefix = "[INF]"
        case Error:
            prefix = "[ERR]"
        }

        getLogger().Println(prefix, msg)
    }
}
