// Inspired by MySQL aggregate functions
package aggregator

import (
	"reflect"
	"syscall"
	"time"
)

type FuncAdd func(*Aggregator, interface{}, interface{}) interface{}

type table struct {
	count   uint64
	funcAdd FuncAdd
	product interface{}
	sum     interface{}
	time0   time.Time
	timeN   time.Time
	value0  interface{}
}

type Aggregator struct {
	db map[string]table
}

type Aggregate struct {
	Avg interface{}
	Cnt interface{}
	Max interface{}
	Min interface{}
	Sum interface{} // Product as well?
}

func NewAggregator() *Aggregator {
	return &Aggregator{db: make(map[string]table)}
}

// Protect against overflow?
// Allow custom add/product functions for structs (e.g., complex numbers)?
func (a *Aggregator) Insert(key string, value interface{}) error {
	var err error

	now := time.Now()

	if entry, ok := a.db[key]; ok {
		if reflect.TypeOf(value) == reflect.TypeOf(entry.value0) {
			entry.count += 1
			entry.sum = entry.funcAdd(a, entry.sum, value)
			entry.timeN = now
		} else {
			err = syscall.EINVAL  // Type mismatch
		}
	} else {
		if adder := a.getAdder(value); adder == nil {
			err = syscall.ENOTSUP
		} else {
			a.db[key] = table{count: 1, funcAdd: adder, sum: value, time0: now, timeN: now, value0: value}
		}
	}

	return err
}

func (a *Aggregator) addFloat32(x, y interface{}) interface{} {
	return x.(float32) + y.(float32)
}

func (a *Aggregator) addFloat64(x, y interface{}) interface{} {
	return x.(float64) + y.(float64)
}

func (a *Aggregator) addInt(x, y interface{}) interface{} {
	return x.(int) + y.(int)
}

func (a *Aggregator) addInt8(x, y interface{}) interface{} {
	return x.(int8) + y.(int8)
}

func (a *Aggregator) addInt16(x, y interface{}) interface{} {
	return x.(int16) + y.(int16)
}

func (a *Aggregator) addInt32(x, y interface{}) interface{} {
	return x.(int32) + y.(int32)
}

func (a *Aggregator) addInt64(x, y interface{}) interface{} {
	return x.(int64) + y.(int64)
}

func (a *Aggregator) addUint(x, y interface{}) interface{} {
	return x.(uint) + y.(uint)
}

func (a *Aggregator) addUint8(x, y interface{}) interface{} {
	return x.(uint8) + y.(uint8)
}

func (a *Aggregator) addUint16(x, y interface{}) interface{} {
	return x.(uint16) + y.(uint16)
}

func (a *Aggregator) addUint32(x, y interface{}) interface{} {
	return x.(uint32) + y.(uint32)
}

func (a *Aggregator) addUint64(x, y interface{}) interface{} {
	return x.(uint64) + y.(uint64)
}

func (a *Aggregator) getAdder(f interface{}) FuncAdd {
	var addFunc FuncAdd

	switch f.(type) {
	case float32:
		addFunc = (*Aggregator).addFloat32
	case float64:
		addFunc = (*Aggregator).addFloat64
	case int:
		addFunc = (*Aggregator).addInt
	case int8:
		addFunc = (*Aggregator).addInt8
	case int16:
		addFunc = (*Aggregator).addInt16
	case int32:
		addFunc = (*Aggregator).addInt32
	case int64:
		addFunc = (*Aggregator).addInt64
	case uint:
		addFunc = (*Aggregator).addUint
	case uint8:
		addFunc = (*Aggregator).addUint8
	case uint16:
		addFunc = (*Aggregator).addUint16
	case uint32:
		addFunc = (*Aggregator).addUint32
	case uint64:
		addFunc = (*Aggregator).addUint64
	}

	return addFunc
}
