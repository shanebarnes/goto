package main

import (
    "bytes"
    "testing"

    "github.com/shanebarnes/goto/logger"
    "github.com/stretchr/testify/assert"
)

var output *bytes.Buffer = new(bytes.Buffer)

func TestLoggerPrintln(t *testing.T) {
    assert := assert.New(t)

    logger.Init(output, 0)

    logger.Println(logger.Error, "Hello, world!")
    assert.Equal("[ERR] Hello, world!\n", output.String())
    output.Reset()

    logger.Println(logger.Info, "Hello, world!")
    assert.Equal("[INF] Hello, world!\n", output.String())
    output.Reset()
}

func TestLoggerSetLevel(t *testing.T) {
    assert := assert.New(t)

    logger.Init(output, 0)

    assert.Equal(logger.Info, logger.GetLevel())

    logger.Println(logger.Debug, "Hello, world!")
    assert.Equal("", output.String())
    output.Reset()

    logger.SetLevel(logger.Debug)
    assert.Equal(logger.Debug, logger.GetLevel())

    logger.Println(logger.Debug, "Hello, world!")
    assert.Equal("[DBG] Hello, world!\n", output.String())
    output.Reset()
}
