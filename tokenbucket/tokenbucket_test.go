package tokenbucket

import (
	"time"
	"testing"
)

func TestTokenBucketNew(t *testing.T) {
	rate := uint64(1000)
	size := uint64(10000)
	tb := New(rate, size)

	if tb.rate != rate {
		t.Errorf("Actual: %v, Expected: %v\n", tb.rate, rate)
	}

	if tb.size != size {
		t.Errorf("Actual: %v, Expected: %v\n", tb.size, size)
	}
}

func TestTokenBucketRemove(t *testing.T) {
	tb1 := New(1000, 10000)

	start := time.Now()
	rv := tb1.Remove(10)
	elapsed := time.Since(start)

	if rv != 10 {
		t.Errorf("Actual: %v, Expected: %v\n", rv, 10)
	}

	if elapsed < (time.Millisecond * 10) {
		t.Errorf("Actual: %v, Expected: %v\n", rv, 10)
	}

	tb2 := New(0, 10000)
	start = time.Now()
	rv = tb2.Remove(10)
	elapsed = time.Since(start)

	if rv != 10 {
		t.Errorf("Actual: %v, Expected: %v\n", rv, 10)
	}

	if elapsed > (time.Millisecond * 1) {
		t.Errorf("Actual: %v, Expected: %v\n", rv, 10)
	}
}

func TestTokenBucketRequest(t *testing.T) {
	tb1 := New(1000, 10000)

	rv := tb1.Request(500)
	if rv != 0 {
		t.Errorf("Actual: %v, Expected: %v\n", rv, 0)
	}

	tb2 := New(0, 10000)
	rv = tb2.Request(tb2.size)
	if rv != tb2.size {
		t.Errorf("Actual: %v, Expected: %v\n", rv, tb2.size)
	}

	tb3 := New(0, 10000)
	rv = tb3.Request(20000)
	if rv != tb3.size {
		t.Errorf("Actual: %v, Expected: %v\n", rv, tb3.size)
	}
}

func TestTokenBucketReturn(t *testing.T) {
	tb1 := New(0, 10000)

	rv := tb1.Return(20000)
	if rv != tb1.size {
		t.Errorf("Actual: %v, Expected: %v\n", rv, tb1.size)
	}

	tb2 := New(1000, 10000)
	rv = tb2.Return(20000)
	if rv != tb2.size {
		t.Errorf("Actual: %v, Expected: %v\n", rv, tb2.size)
	}
}
