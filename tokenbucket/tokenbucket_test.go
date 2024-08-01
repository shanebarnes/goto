package tokenbucket

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func BenchmarkTokenBucket_Remove(b *testing.B) {
	b.Run("burst=1,rate=0", func(b *testing.B) {
		tb := New(0, 1)
		for n := 0; n < b.N; n++ {
			assert.EqualValues(b, 1, tb.Remove(1))
		}
	})

	b.Run("burst=100k,rate=1M", func(b *testing.B) {
		tb := New(1000000, 100000)
		for n := 0; n < b.N; n++ {
			assert.EqualValues(b, 100000, tb.Remove(100000))
		}
	})
}

func TestNew(t *testing.T) {
	tb := New(1000, 10000)
	require.NotNil(t, tb)
}

func TestTokenBucketRemove(t *testing.T) {
	tests := []struct {
		burst uint64
		rate  uint64
	}{
		{burst: 10, rate: 0},
		{burst: 100, rate: 0},
		{burst: 200, rate: 0},
		{burst: 10, rate: 100},
		{burst: 100, rate: 100},
		{burst: 200, rate: 100},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("burst=%d,rate=%d", test.burst, test.rate), func(t *testing.T) {
			start := time.Now()
			tb := New(test.rate, test.burst)
			assert.EqualValues(t, test.burst, tb.Remove(test.burst))
			switch {
			case test.rate == 0:
				assert.Less(t, time.Since(start), time.Millisecond)
			default:
				assert.Less(t, time.Since(start), time.Second*time.Duration(test.burst)/time.Duration(test.rate))
			}

			start = time.Now()
			assert.Equal(t, test.burst, tb.Remove(test.burst))
			switch {
			case test.rate == 0:
				assert.Less(t, time.Since(start), time.Millisecond)
			default:
				assert.GreaterOrEqual(t, time.Since(start), time.Second*time.Duration(test.burst)/time.Duration(test.rate))
			}
		})
	}
}

func TestTokenBucketRemoveWithContext(t *testing.T) {
	t.Run("contextCanceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		tb := New(1, 1)

		tokens, err := tb.RemoveWithContext(ctx, 1)
		assert.ErrorIs(t, err, context.Canceled)
		assert.Zero(t, tokens)
	})

	t.Run("contextDeadlineExceeded", func(t *testing.T) {
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second))
		defer cancel()
		tb := New(1, 1)

		// Will block here for a long time if context is not canceled.
		tokens, err := tb.RemoveWithContext(ctx, 1000)
		assert.ErrorIs(t, err, context.DeadlineExceeded)
		assert.Zero(t, tokens)
	})

	tests := []struct {
		burst uint64
		rate  uint64
	}{
		{burst: 10, rate: 0},
		{burst: 100, rate: 0},
		{burst: 200, rate: 0},
		{burst: 10, rate: 100},
		{burst: 100, rate: 100},
		{burst: 200, rate: 100},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("burst=%d,rate=%d", test.burst, test.rate), func(t *testing.T) {
			start := time.Now()
			tb := New(test.rate, test.burst)
			assert.EqualValues(t, test.burst, tb.Remove(test.burst))
			switch {
			case test.rate == 0:
				assert.Less(t, time.Since(start), time.Millisecond)
			default:
				assert.Less(t, time.Since(start), time.Second*time.Duration(test.burst)/time.Duration(test.rate))
			}

			start = time.Now()
			tokens, err := tb.RemoveWithContext(context.Background(), test.burst)
			assert.NoError(t, err)
			assert.Equal(t, test.burst, tokens)
			switch {
			case test.rate == 0:
				assert.Less(t, time.Since(start), time.Millisecond)
			default:
				assert.GreaterOrEqual(t, time.Since(start), time.Second*time.Duration(test.burst)/time.Duration(test.rate))
			}
		})
	}
}

func TestTokenBucketReturn(t *testing.T) {
	tests := []struct {
		burst  uint64
		rate   uint64
		unused uint64
	}{
		{burst: 10, rate: 100, unused: 0},
		{burst: 10, rate: 100, unused: 10},
		{burst: 100, rate: 100, unused: 100},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("burst=%d,rate=%d,unused=%d", test.burst, test.rate, test.unused), func(t *testing.T) {
			tb := New(test.rate, test.burst)
			assert.Equal(t, test.burst, tb.Remove(test.burst))
			assert.Equal(t, test.unused, tb.Return(test.unused))
		})
	}
}
