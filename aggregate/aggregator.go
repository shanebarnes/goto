// Inspired by MySQL aggregate functions
package aggregator

import (
	"reflect"
	"sync"
	"syscall"
	"time"
)

type record struct {
	count  int64
	max    interface{}
	min    interface{}
	sum    interface{}
	time0  time.Time
	timeN  time.Time
	value0 interface{}
}

type Aggregator struct {
	db map[interface{}]*record
	mu *sync.RWMutex
}

type Aggregate struct {
	Avg interface{}
	Cnt interface{}
	Max interface{}
	Min interface{}
	Sum interface{}
}

func NewAggregator() *Aggregator {
	return &Aggregator{db: make(map[interface{}]*record), mu: &sync.RWMutex{}}
}

func (a *Aggregator) add(x, y interface{}) interface{} {
	var sum interface{}

	switch x.(type) {
	case float64:
		sum = x.(float64) + y.(float64)
	case int64:
		sum = x.(int64) + y.(int64)
	case uint64:
		sum = x.(uint64) + y.(uint64)
	}

	return sum
}

// Compare two integer interfaces values:
//   -1: x <  y
//    0: x == y
//    1: x >  y
func (a *Aggregator) compare(x interface{}, y interface{}) int {
	res := 0

	switch x.(type) {
	case float32:
		if x.(float32) < y.(float32) {
			res = -1
		} else if x.(float32) > y.(float32) {
			res = 1
		}
	case float64:
		if x.(float64) < y.(float64) {
			res = -1
		} else if x.(float64) > y.(float64) {
			res = 1
		}
	case int:
		if x.(int) < y.(int) {
			res = -1
		} else if x.(int) > y.(int) {
			res = 1
		}
	case int8:
		if x.(int8) < y.(int8) {
			res = -1
		} else if x.(int8) > y.(int8) {
			res = 1
		}
	case int16:
		if x.(int16) < y.(int16) {
			res = -1
		} else if x.(int16) > y.(int16) {
			res = 1
		}
	case int32:
		if x.(int32) < y.(int32) {
			res = -1
		} else if x.(int32) > y.(int32) {
			res = 1
		}
	case int64:
		if x.(int64) < y.(int64) {
			res = -1
		} else if x.(int64) > y.(int64) {
			res = 1
		}
	case uint:
		if x.(uint) < y.(uint) {
			res = -1
		} else if x.(uint) > y.(uint) {
			res = 1
		}
	case uint8:
		if x.(uint8) < y.(uint8) {
			res = -1
		} else if x.(uint8) > y.(uint8) {
			res = 1
		}
	case uint16:
		if x.(uint16) < y.(uint16) {
			res = -1
		} else if x.(uint16) > y.(uint16) {
			res = 1
		}
	case uint32:
		if x.(uint32) < y.(uint32) {
			res = -1
		} else if x.(uint32) > y.(uint32) {
			res = 1
		}
	case uint64:
		if x.(uint64) < y.(uint64) {
			res = -1
		} else if x.(uint64) > y.(uint64) {
			res = 1
		}
	}

	return res
}

func (a *Aggregator) convert(value interface{}) (interface{}, error) {
	var err error
	var f interface{}

	switch value.(type) {
	case float32:
		f = float64(value.(float32))
	case float64:
		f = value.(float64)
	case int:
		f = int64(value.(int))
	case int8:
		f = int64(value.(int8))
	case int16:
		f = int64(value.(int16))
	case int32:
		f = int64(value.(int32))
	case int64:
		f = value.(int64)
	case uint:
		f = uint64(value.(uint))
	case uint8:
		f = uint64(value.(uint8))
	case uint16:
		f = uint64(value.(uint16))
	case uint32:
		f = uint64(value.(uint32))
	case uint64:
		f = uint64(value.(uint64))
	default:
		err = syscall.ENOTSUP
	}

	return f, err
}

func (a *Aggregator) Delete(key interface{}) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	_, err := a.findRecord(key)
	if err == nil {
		delete(a.db, key)
	}

	return err
}

func (a *Aggregator) findRecord(key interface{}) (*record, error) {
	var err error

	rec, ok := a.db[key]
	if !ok {
		err = syscall.ENOENT
	}

	return rec, err
}

func (a *Aggregator) GetAverage(key interface{}) (interface{}, error) {
	var avg interface{}
	a.mu.RLock()
	defer a.mu.RUnlock()

	rec, err := a.findRecord(key)
	if err == nil {
		if rec.count > 0 {
			switch rec.value0.(type) {
			case float32, float64:
				return rec.sum.(float64) / float64(rec.count), err
			case int, int8, int16, int32, int64:
				return rec.sum.(int64) / rec.count, err
			case uint, uint8, uint16, uint32, uint64:
				return rec.sum.(uint64) / uint64(rec.count), err
			default:
				err = syscall.ENOTSUP
			}
		} else {
			avg = rec.sum
		}
	}

	return avg, err
}

func (a *Aggregator) GetCount(key interface{}) (int64, error) {
	var cnt int64 = -1
	a.mu.RLock()
	defer a.mu.RUnlock()

	rec, err := a.findRecord(key)
	if err == nil {
		cnt = rec.count
	}

	return cnt, err
}

func (a *Aggregator) GetDuration(key interface{}) (time.Duration, error) {
	var dur time.Duration
	a.mu.RLock()
	defer a.mu.RUnlock()

	rec, err := a.findRecord(key)
	if err == nil {
		dur = rec.timeN.Sub(rec.time0)
	}

	return dur, err
}

func (a *Aggregator) GetMaximum(key interface{}) (interface{}, error) {
	var max interface{}
	a.mu.RLock()
	defer a.mu.RUnlock()

	rec, err := a.findRecord(key)
	if err == nil {
		max = rec.max
	}

	return max, err
}

func (a *Aggregator) GetMinimum(key interface{}) (interface{}, error) {
	var min interface{}
	a.mu.RLock()
	defer a.mu.RUnlock()

	rec, err := a.findRecord(key)
	if err == nil {
		min = rec.min
	}

	return min, err
}

func (a *Aggregator) GetRate(key interface{}, dur time.Duration) (interface{}, error) {
	var rate interface{}
	a.mu.RLock()
	defer a.mu.RUnlock()

	rec, err := a.findRecord(key)
	if err == nil {
		elapsedNsec := rec.timeN.Sub(rec.time0).Nanoseconds()
		switch rec.value0.(type) {
		case float32, float64:
			if elapsedNsec != 0 {
				// Protect against overflow on duration multiplier?
				rate = rec.sum.(float64) * float64(dur.Nanoseconds()) / float64(elapsedNsec)
			} else {
				rate = float64(0)
			}
		case int, int8, int16, int32, int64:
			if elapsedNsec != 0 {
				// Protect against overflow on duration multiplier?
				rate = rec.sum.(int64) * dur.Nanoseconds() / elapsedNsec
			} else {
				rate = int64(0)
			}
		case uint, uint8, uint16, uint32, uint64:
			if dur >= 0 {
				if elapsedNsec > 0 {
					// Protect against overflow on duration multiplier?
					rate = rec.sum.(uint64) * uint64(dur.Nanoseconds()) / uint64(elapsedNsec)
				} else if elapsedNsec == 0 {
					rate = uint64(0)
				} else {
					err = syscall.EINVAL
				}
			} else {
				err = syscall.EINVAL
			}
		default:
			err = syscall.ENOTSUP
		}
	}

	return rate, err
}

func (a *Aggregator) GetSum(key interface{}) (interface{}, error) {
	var sum interface{}
	a.mu.RLock()
	defer a.mu.RUnlock()

	rec, err := a.findRecord(key)
	if err == nil {
		sum = rec.sum
	}

	return sum, err
}

// Protect against overflow?
// Allow custom add/product functions for structs (e.g., complex numbers)?
func (a *Aggregator) Insert(key interface{}, value interface{}) error {
	var err error
	now := time.Now()

	a.mu.Lock()
	defer a.mu.Unlock()
	if entry, ok := a.db[key]; ok {
		if reflect.TypeOf(value) == reflect.TypeOf(entry.value0) {
			var newValue interface{}
			if newValue, err = a.convert(value); err == nil {
				entry.count += 1
				entry.sum = a.add(entry.sum, newValue)
				entry.timeN = now

				if a.compare(value, entry.min) == -1 {
					entry.min = value
				}

				if a.compare(value, entry.max) == 1 {
					entry.max = value
				}
			}
		} else {
			err = syscall.EINVAL
		}
	} else {
		var newValue interface{}
		if newValue, err = a.convert(value); err == nil {
			a.db[key] = &record{
				count:  1,
				max:    value,
				min:    value,
				sum:    newValue,
				time0:  now,
				timeN:  now,
				value0: value,
			}
		}
	}

	return err
}
