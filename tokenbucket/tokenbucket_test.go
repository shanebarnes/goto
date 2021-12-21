package tokenbucket

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTokenBucketNew(t *testing.T) {
	rate := uint64(1000)
	size := uint64(10000)
	tb := New(rate, size)

	assert.Equal(t, rate, tb.rate)
	assert.Equal(t, size, tb.size)
}

func TestTokenBucketRemove(t *testing.T) {
	tb1 := New(1000, 10)

	start := time.Now()
	rv := tb1.Remove(10)
	rv = tb1.Remove(10)
	elapsed := time.Since(start)

	assert.Equal(t, uint64(10), rv)
	assert.True(t, elapsed >= (time.Millisecond*10))

	tb2 := New(0, 10000)
	start = time.Now()
	rv = tb2.Remove(10)
	elapsed = time.Since(start)

	assert.Equal(t, uint64(10), rv)
	assert.True(t, elapsed <= (time.Millisecond*1))
}

func TestTokenBucketRequest(t *testing.T) {
	tb1 := New(1000, 100)

	rv := tb1.Request(500)
	assert.NotEqual(t, 500, rv)

	tb2 := New(0, 10000)
	rv = tb2.Request(tb2.size)
	assert.Equal(t, tb2.size, rv)

	tb3 := New(0, 10000)
	rv = tb3.Request(20000)
	assert.Equal(t, tb3.size, rv)
}

func TestTokenBucketReturn(t *testing.T) {
	tb1 := New(1000, 10000)

	tb1.Remove(10000)
	rv := tb1.Return(20000)
	assert.Equal(t, uint64(10000), rv)
}
