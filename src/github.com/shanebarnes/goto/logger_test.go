package main

import (
    "bytes"
    "testing"

    "github.com/shanebarnes/goto/logger"
    "github.com/stretchr/testify/assert"
)

func TestLoggerPrintln(t *testing.T) {
    assert := assert.New(t)
    buf := new(bytes.Buffer)

    logger.Init(buf, 0)

    logger.Println(logger.Error, "Hello, world!")
    assert.Equal("[ERR] Hello, world!\n", buf.String())
    buf.Reset()

    logger.Println(logger.Info, "Hello, world!")
    assert.Equal("[INF] Hello, world!\n", buf.String())
    buf.Reset()
}

func TestLoggerSetLevel(t *testing.T) {
    assert := assert.New(t)
    buf := new(bytes.Buffer)

    logger.Init(buf, 0)

    assert.Equal(logger.Info, logger.GetLevel())

    logger.Println(logger.Debug, "Hello, world!")
    assert.Equal("", buf.String())
    buf.Reset()

    logger.SetLevel(logger.Debug)
    assert.Equal(logger.Debug, logger.GetLevel())

    logger.Println(logger.Debug, "Hello, world!")
    assert.Equal("[DBG] Hello, world!\n", buf.String())
}
