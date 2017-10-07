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
    Always
)

var loggerFlags     int = log.Ldate | log.Ltime | log.Lmicroseconds
var loggerLevel     int32 = int32(Info)
var loggerInstance *log.Logger = nil
var loggerOnce      sync.Once
var loggerWriter    io.Writer = os.Stdout

func GetLevel() Level {
    return Level(atomic.LoadInt32(&loggerLevel))
}

func getLogger() *log.Logger {
    loggerOnce.Do(func() {
        loggerInstance = log.New(loggerWriter, "", loggerFlags)
    })
    return loggerInstance
}

func Init(writer io.Writer, flags int) {
    loggerWriter = writer
    loggerFlags = flags
}

func SetLevel(level Level) {
    atomic.StoreInt32(&loggerLevel, int32(level))
}

func println(level Level, msg *string) {
    if level >= GetLevel() {
        prefix := ""

        switch level {
        case Debug:
            prefix = "[DBG]"
        case Info:
            prefix = "[INF]"
        case Error:
            prefix = "[ERR]"
        case Always:
            prefix = "[ALW]"
        }

        getLogger().Println(prefix, *msg)
    }
}

func PrintlnAlways(msg string) {
    println(Always, &msg)
}

func PrintlnDebug(msg string) {
    println(Debug, &msg)
}

func PrintlnError(msg string) {
    println(Error, &msg)
}

func PrintlnInfo(msg string) {
    println(Info, &msg)
}
