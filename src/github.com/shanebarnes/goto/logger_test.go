package main

import (
    "bytes"
    "testing"

    "github.com/shanebarnes/goto/logger"
    "github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
    assert := assert.New(t)
    buf := new(bytes.Buffer)

    logger.Init(buf, 0)

    logger.DebugLn("Debug, world!")
    assert.Equal("[DBG] Debug, world!\n", buf.String())
    buf.Reset()

    logger.ErrorLn("Error, world!")
    assert.Equal("[ERR] Error, world!\n", buf.String())
    buf.Reset()

    logger.InfoLn("Hello, world!")
    assert.Equal("[INF] Hello, world!\n", buf.String())
    buf.Reset()
}
