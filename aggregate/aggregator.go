// Inspired by MySQL aggregate functions
package aggregator

import (
	"reflect"
	"syscall"
	"time"
)

// Calculate slope: dY/dX, Calculate Rate: dY/dt?
type record struct {
	count   int64
	max     interface{}
	min     interface{}
	sum     interface{}
	time0   time.Time
	timeN   time.Time
	value0  interface{}
}

type Aggregator struct {
	db map[string]*record
}

type Aggregate struct {
	Avg interface{}
	Cnt interface{}
	Max interface{}
	Min interface{}
	Sum interface{} // Product as well?
}

func NewAggregator() *Aggregator {
	return &Aggregator{db: make(map[string]*record)}
}

func (a *Aggregator) add(x interface{}, y interface{}) interface{} {
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
	case float64:
		if x.(float64) < y.(float64) {
			res = -1
		} else if x.(float64) > y.(float64) {
			res = 1
		}
	case int64:
		if x.(int64) < y.(int64) {
			res = -1
		} else if x.(int64) > y.(int64) {
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
		f = int64(value.(int64))
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
		err = syscall.ENOTSUP // Unsupported type
	}

	return f, err
}

func (a *Aggregator) Delete(key string) error {
	var err error

	if _, err = a.findRecord(key); err == nil {
		delete(a.db, key)
	}

	return err
}

func (a *Aggregator) findRecord(key string) (*record, error) {
	var err error

	rec, ok := a.db[key]
	if !ok {
		err = syscall.ENOENT
	}

	return rec, err
}

func (a *Aggregator) GetAverage(key string) (interface{}, error) {
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
			return rec.sum, err
		}
	}

	return nil, err
}

func (a *Aggregator) GetCount(key string) (int64, error) {
	rec, err := a.findRecord(key)
	if err == nil {
		return rec.count, nil
	}

	return -1, err
}

func (a *Aggregator) GetMaximum(key string) (interface{}, error) {
	var max interface{}
	rec, err := a.findRecord(key)
	if err == nil {
		max = rec.max
	}

	return max, err

}

func (a *Aggregator) GetMinimum(key string) (interface{}, error) {
	var min interface{}
	rec, err := a.findRecord(key)
	if err == nil {
		min = rec.min
	}

	return min, err

}

func (a *Aggregator) GetSum(key string) (interface{}, error) {
	var sum interface{}
	rec, err := a.findRecord(key)
	if err == nil {
		sum = rec.sum
	}

	return sum, err
}

// Protect against overflow?
// Allow custom add/product functions for structs (e.g., complex numbers)?
func (a *Aggregator) Insert(key string, value interface{}) error {
	var err error

	now := time.Now()

	if entry, ok := a.db[key]; ok {
		if reflect.TypeOf(value) == reflect.TypeOf(entry.value0) {
			var newValue interface{}
			if newValue, err = a.convert(value); err == nil {
				entry.count += 1
				entry.sum = a.add(entry.sum, newValue)
				entry.timeN = now

				if a.compare(newValue, entry.min) == -1 {
					entry.min = newValue
				}

				if a.compare(newValue, entry.max) == 1 {
					entry.max = newValue
				}
			}
		} else {
			err = syscall.EINVAL  // Invalid type
		}
	} else {
		var newValue interface{}
		if newValue, err = a.convert(value); err == nil {
			a.db[key] = &record{
				count: 1,
				max: newValue,
				min: newValue,
				sum: newValue,
				time0: now,
				timeN: now,
				value0: value,
			}
		}
	}

	return err
}
