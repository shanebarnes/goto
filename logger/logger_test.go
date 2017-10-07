package logger

import (
    "bytes"
    "testing"
)

var output *bytes.Buffer = new(bytes.Buffer)

func TestLoggerInit(t *testing.T) {
    message := "Hello, world!"
    expected := "[INF] Hello, world!\n"

    Init(output, 0)
    PrintlnInfo(message)

    if output.String() != expected {
        t.Errorf("Actual: %s, Expected: %s\n", output.String(), expected)
    }

    output.Reset()
}

func TestLoggerSetLevel(t *testing.T) {
    levels := []Level{ Debug, Info, Error, Always }

    for i := range levels {
        SetLevel(levels[i])

        if GetLevel() != levels[i] {
            t.Errorf("Actual: %d, Expected: %d\n", GetLevel(), levels[i])
        }
    }
}

func TestLoggerPrintlnDebug(t *testing.T) {
    levels := []Level{ Debug, Info, Error, Always }
    expected := "[DBG] Debug\n"

    for i := range levels {
        SetLevel(levels[i])
        PrintlnDebug("Debug")

        if GetLevel() > Debug && output.String() != "" {
            t.Errorf("Actual: %s, Expected: %s\n", output.String(), "")
        } else if GetLevel() <= Debug && output.String() != expected {
            t.Errorf("Actual: %s, Expected: %s\n", output.String(), expected)
        }

        output.Reset()
    }
}

func TestLoggerPrintlnInfo(t *testing.T) {
    levels := []Level{ Debug, Info, Error, Always }
    expected := "[INF] Info\n"

    for i := range levels {
        SetLevel(levels[i])
        PrintlnInfo("Info")

        if GetLevel() > Info && output.String() != "" {
            t.Errorf("Actual: %s, Expected: %s\n", output.String(), "")
        } else if GetLevel() <= Info && output.String() != expected {
            t.Errorf("Actual: %s, Expected: %s\n", output.String(), expected)
        }

        output.Reset()
    }
}

func TestLoggerPrintlnError(t *testing.T) {
    levels := []Level{ Debug, Info, Error, Always }
    expected := "[ERR] Error\n"

    for i := range levels {
        SetLevel(levels[i])
        PrintlnError("Error")

        if GetLevel() > Error && output.String() != "" {
            t.Errorf("Actual: %s, Expected: %s\n", output.String(), "")
        } else if GetLevel() <= Error && output.String() != expected {
            t.Errorf("Actual: %s, Expected: %s\n", output.String(), expected)
        }

        output.Reset()
    }
}

func TestLoggerPrintlnAlways(t *testing.T) {
    levels := []Level{ Debug, Info, Error, Always }
    expected := "[ALW] Always\n"

    for i := range levels {
        SetLevel(levels[i])
        PrintlnAlways("Always")

        if output.String() != expected {
            t.Errorf("Actual: %s, Expected: %s\n", output.String(), expected)
        }

        output.Reset()
    }
}
