package tokenbucket

import (
	"context"
	"time"
)

// TokenBucket implements a non thread-safe token bucket algorithm.
type TokenBucket struct {
	// The number of tokens currently in the bucket.
	fill uint64

	// The rate at which tokens are added to the bucket measured in tokens per
	// second.
	rate uint64

	// The capacity of the bucket measured in tokens.
	size uint64

	// Last time that tokens were added to the bucket.
	time time.Time

	wait *time.Timer
}

// New creates and returns a new TokenBucket instance.
func New(rate uint64, size uint64) *TokenBucket {
	if size == 0 {
		panic("bucket size must be greater than value")
	}

	start := time.Now()
	t := time.NewTimer(0)
	<-t.C
	t.Stop()

	return &TokenBucket{
		fill: size,
		rate: rate,
		size: size,
		time: start,
		wait: t,
	}
}

// Remove blocks until all requested tokens are available to be removed from
// the bucket.
func (tb *TokenBucket) Remove(tokens uint64) uint64 {
	n, _ := tb.RemoveWithContext(context.Background(), tokens)
	return n
}

// RemoveWithContext blocks until all requested tokens are available to be
// removed from the bucket or the context is canceled.
func (tb *TokenBucket) RemoveWithContext(ctx context.Context, tokens uint64) (uint64, error) {
	n := tb.request(tokens)

	switch {
	case ctx.Err() != nil:
		return 0, ctx.Err()
	case n >= tokens || tb.rate == 0:
		return tokens, nil
	default:
		deadline := time.Unix(0, int64(tokens-n)*int64(time.Second)/int64(tb.rate)+tb.time.UnixNano())
		duration := time.Until(deadline)

		tb.wait.Reset(duration)
		defer tb.wait.Stop()

		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case <-tb.wait.C:
			n += tb.request(tokens - n)
			if n > tokens {
				tb.Return(n - tokens)
			}
			return tokens, nil
		}
	}
}

func (tb *TokenBucket) request(tokens uint64) uint64 {
	switch {
	case tb.rate == 0:
		return min(tb.size, tokens)
	case tb.fill >= tokens:
		tb.fill -= tokens
		return tokens
	default:
		now := time.Now()

		if newTokens := tb.rate * uint64(now.Sub(tb.time).Nanoseconds()) / uint64(time.Second); newTokens > 0 {
			tb.fill += newTokens
			tb.time = now
		}

		tb.fill = min(tb.size, tb.fill)

		if tb.fill >= tokens {
			tb.fill -= tokens
			return tokens
		}

		tokens = tb.fill
		tb.fill = 0
		return tokens
	}
}

// Return allows unused tokens retrieved with Remove or RemoveWithContext to
// refill the bucket.
func (tb *TokenBucket) Return(tokens uint64) uint64 {
	switch {
	case tokens == 0 || tb.fill == tb.size:
		return 0
	case tokens+tb.fill > tb.size:
		tokens = tb.size - tb.fill
		tb.fill = tb.size
		return tokens
	default:
		tb.fill += tokens
		return tokens
	}
}
