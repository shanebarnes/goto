package logger

import (
	"io"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type Level int32

const (
	Debug Level = iota
	Info
	Error
	Always
)

type loggerObject struct {
	Flags  int
	Level  int32
	Logger *log.Logger
	Writer io.Writer
}

type logger struct {
	Jam  []loggerObject
	Once sync.Once
}

var instance logger

func GetLevel(index int) Level {
	var ret Level = Always
	logger := getLogger()

	if index < len(logger.Jam) {
		ret = Level(atomic.LoadInt32(&logger.Jam[index].Level))
	}

	return ret
}

func getLogger() *logger {
	instance.Once.Do(func() {
		for i := range instance.Jam {
			instance.Jam[i].Logger = log.New(instance.Jam[i].Writer,
				"",
				instance.Jam[i].Flags)
		}
	})
	return &instance
}

func Init(flags int, level Level, writer io.Writer) {
	instance.Jam = append(instance.Jam, loggerObject{Flags: flags, Level: int32(level), Writer: writer})
}

func SetLevel(index int, level Level) {
	logger := getLogger()

	if index < len(logger.Jam) {
		atomic.StoreInt32(&logger.Jam[index].Level, int32(level))
	}
}

func println(level Level, v ...interface{}) {
	logger := getLogger()

	for i := range logger.Jam {
		if level >= GetLevel(i) {
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

			m := append([]interface{}{prefix}, v...)
			logger.Jam[i].Logger.Println(m...)
		}
	}
}

func PrintlnAlways(v ...interface{}) {
	println(Always, v...)
}

func PrintlnDebug(v ...interface{}) {
	println(Debug, v...)
}

func PrintlnError(v ...interface{}) {
	println(Error, v...)
}

func PrintlnInfo(v ...interface{}) {
	println(Info, v...)
}

func PrintlnTime(l Level, start time.Time, v ...interface{}) func() {
	if start.IsZero() {
		start = time.Now()
	}

	return func() {
		m := append([]interface{}{"(exec time: " + time.Since(start).String() + ")"}, v...)
		println(l, m...)
	}
}
