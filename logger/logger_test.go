package logger

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

var output *bytes.Buffer = new(bytes.Buffer)

func TestLoggerInit(t *testing.T) {
	message := "Hello, world!"
	expected := "[INF] Hello, world!\n"

	Init(0, Info, output)
	PrintlnInfo(message)

	if output.String() != expected {
		t.Errorf("Actual: %s, Expected: %s\n", output.String(), expected)
	}

	output.Reset()
}

func TestLoggerSetLevel(t *testing.T) {
	levels := []Level{Debug, Info, Error, Always}

	for i := range levels {
		SetLevel(0, levels[i])

		if GetLevel(0) != levels[i] {
			t.Errorf("Actual: %d, Expected: %d\n", GetLevel(0), levels[i])
		}
	}
}

func TestLoggerPrintlnDebug(t *testing.T) {
	levels := []Level{Debug, Info, Error, Always}
	expected := "[DBG] Debug\n"

	for i := range levels {
		SetLevel(0, levels[i])
		PrintlnDebug("Debug")

		if GetLevel(0) > Debug && output.String() != "" {
			t.Errorf("Actual: %s, Expected: %s\n", output.String(), "")
		} else if GetLevel(0) <= Debug && output.String() != expected {
			t.Errorf("Actual: %s, Expected: %s\n", output.String(), expected)
		}

		output.Reset()
	}
}

func TestLoggerPrintlnInfo(t *testing.T) {
	levels := []Level{Debug, Info, Error, Always}
	expected := "[INF] Info\n"

	for i := range levels {
		SetLevel(0, levels[i])
		PrintlnInfo("Info")

		if GetLevel(0) > Info && output.String() != "" {
			t.Errorf("Actual: %s, Expected: %s\n", output.String(), "")
		} else if GetLevel(0) <= Info && output.String() != expected {
			t.Errorf("Actual: %s, Expected: %s\n", output.String(), expected)
		}

		output.Reset()
	}
}

func TestLoggerPrintlnError(t *testing.T) {
	levels := []Level{Debug, Info, Error, Always}
	expected := "[ERR] Error\n"

	for i := range levels {
		SetLevel(0, levels[i])
		PrintlnError("Error")

		if GetLevel(0) > Error && output.String() != "" {
			t.Errorf("Actual: %s, Expected: %s\n", output.String(), "")
		} else if GetLevel(0) <= Error && output.String() != expected {
			t.Errorf("Actual: %s, Expected: %s\n", output.String(), expected)
		}

		output.Reset()
	}
}

func TestLoggerPrintlnAlways(t *testing.T) {
	levels := []Level{Debug, Info, Error, Always}
	expected := "[ALW] Always\n"

	for i := range levels {
		SetLevel(0, levels[i])
		PrintlnAlways("Always")

		if output.String() != expected {
			t.Errorf("Actual: %s, Expected: %s\n", output.String(), expected)
		}

		output.Reset()
	}
}

func TestLoggerPrintlnTime(t *testing.T) {
	defer func() {
		var elapsedMs float64 = 0.
		n, err := fmt.Sscanf(output.String(), "[ALW] (exec time: %fms) Waited for 10ms", &elapsedMs)

		if n != 1 || err != nil || elapsedMs < 10. {
			t.Errorf("Actual: %s %d %s\n", output.String(), n, err)
		}
		output.Reset()
	}()

	defer PrintlnTime(Always, time.Time{}, "Waited for 10ms")()
	time.Sleep(time.Millisecond * 10)
}
