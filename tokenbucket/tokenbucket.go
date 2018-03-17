package tokenbucket

import (
	"time"
)

type TokenBucket struct {
	fill  uint64 // The number of tokens currently in the bucket
	rate  uint64 // The rate at which tokens are added to the bucket measured in tokens per second
	size  uint64 // The capacity of the bucket measured in tokens
	time  int64  // Unix timestamp in nanoseconds indicating the last time that tokens were added to the bucket
	wait *time.Timer
}

func New(rate uint64, size uint64) *TokenBucket {
	if size == 0 {
		panic("size: must be a non-zero value")
	}

	bucket := new(TokenBucket)
	bucket.time = time.Now().UnixNano()
	bucket.rate = rate
	bucket.fill = 0
	bucket.size = size
	bucket.wait = time.NewTimer(0 * time.Second)

	<-bucket.wait.C
	bucket.wait.Stop()

	return bucket
}

func (tb *TokenBucket) Remove(tokens uint64) uint64 {
	// A remove is a request that blocks until are tokens to remove are available
	rv := tb.Request(tokens)

	if rv < tokens {
		if tb.rate > 0 {
			deadline := time.Unix(0, int64(tokens - rv) * int64(time.Second) / int64(tb.rate) + int64(tb.time))
			duration := time.Until(deadline)

			tb.wait.Reset(duration)
			<-tb.wait.C
			tb.wait.Stop()

			rv = rv + tb.Request(tokens - rv)
		}
	}

	return rv
}

func (tb *TokenBucket) Request(tokens uint64) uint64 {
	var rv uint64 = 0

	if tb != nil {
		if tb.rate <= 0 {
			if tb.size < tokens {
				rv = tb.size
			} else {
				rv = tokens
			}
		} else if tb.fill >= tokens {
			tb.fill = tb.fill - tokens
			rv = tokens
		} else {
			now := time.Now().UnixNano()
			newTokens := tb.rate * uint64(now-tb.time) / uint64(time.Second)

			if newTokens > 0 {
				tb.fill = tb.fill + newTokens
				tb.time = now
			}

			if tb.fill > tb.size {
				tb.fill = tb.size
			}

			if tb.fill >= tokens {
				tb.fill = tb.fill - tokens
				rv = tokens
			} else {
				rv = tb.fill
				tb.fill = 0
			}
		}
	}

	return rv
}

func (tb *TokenBucket) Return(tokens uint64) uint64 {
	var rv uint64 = 0

	if tb != nil && tb.fill < tb.size {
		tb.fill = tb.fill + tokens

		if tb.fill > tb.size {
			rv = tokens - (tb.fill - tb.size)
			tb.fill = tb.size
		} else {
			rv = tokens
		}
	}

	return rv
}
