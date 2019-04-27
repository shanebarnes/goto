package aggregator

import (
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAggregator_Delete(t *testing.T) {
	a := NewAggregator()

	assert.Equal(t, syscall.ENOENT, a.Delete("badKey"))

	assert.Nil(t, a.Insert("key1", 1))
	assert.Equal(t, len(a.db), 1)

	assert.Nil(t, a.Delete("key1"))
	assert.Equal(t, len(a.db), 0)

	assert.Nil(t, a.Insert("key2", 2))
	assert.Nil(t, a.Insert("key3", 3))
	assert.Nil(t, a.Insert("key4", 4))

	assert.Nil(t, a.Delete("key3"))
	assert.Equal(t, len(a.db), 2)
}

func TestAggregator_GetAverage(t *testing.T) {
	a := NewAggregator()
	key1 := "int"
	key2 := "float32"

	avg, err := a.GetAverage("badKey")
	assert.Equal(t, syscall.ENOENT, err)

	assert.Nil(t, a.Insert(key1, 0))
	assert.Nil(t, a.Insert(key1, 10))
	assert.Nil(t, a.Insert(key1, 20))
	assert.Nil(t, a.Insert(key2, float32(15.5)))

	avg, err = a.GetAverage(key1)
	assert.Nil(t, err)
	assert.Equal(t, int64(10), avg)

	avg, err = a.GetAverage(key2)
	assert.Nil(t, err)
	assert.Equal(t, float64(15.5), avg)
}

func TestAggregator_GetCount(t *testing.T) {
	a := NewAggregator()
	key1 := "int16"
	key2 := "uint32"

	count, err := a.GetCount("badKey")
	assert.Equal(t, syscall.ENOENT, err)

	assert.Nil(t, a.Insert(key1, int16(-1)))
	assert.Nil(t, a.Insert(key1, int16(-2)))
	assert.Nil(t, a.Insert(key1, int16(-4)))
	assert.Nil(t, a.Insert(key2, uint32(8)))

	count, err = a.GetCount(key1)
	assert.Nil(t, err)
	assert.Equal(t, int64(3), count)

	count, err = a.GetCount(key2)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), count)
}

func TestAggregator_GetMaximum(t *testing.T) {
	a := NewAggregator()
	key1 := "float64"
	key2 := "uint16"

	max, err := a.GetMaximum("badKey")
	assert.Equal(t, syscall.ENOENT, err)

	assert.Nil(t, a.Insert(key1, 1.))
	max, err = a.GetMaximum(key1)
	assert.Nil(t, err)
	assert.Equal(t, 1., max)

	assert.Nil(t, a.Insert(key1, 20.))
	max, err = a.GetMaximum(key1)
	assert.Nil(t, err)
	assert.Equal(t, 20., max)

	assert.Nil(t, a.Insert(key1, 10.))
	max, err = a.GetMaximum(key1)
	assert.Nil(t, err)
	assert.Equal(t, 20., max)

	assert.Nil(t, a.Insert(key2, uint16(40)))
	max, err = a.GetMaximum(key2)
	assert.Nil(t, err)
	assert.Equal(t, int64(40), max)
}

func TestAggregator_GetMinimum(t *testing.T) {
	a := NewAggregator()
	key1 := "int32"
	key2 := "uint8"

	min, err := a.GetMinimum("badKey")
	assert.Equal(t, syscall.ENOENT, err)

	assert.Nil(t, a.Insert(key1, int32(1)))
	min, err = a.GetMinimum(key1)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), min)

	assert.Nil(t, a.Insert(key1, int32(20)))
	min, err = a.GetMinimum(key1)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), min)

	assert.Nil(t, a.Insert(key1, int32(-30)))
	min, err = a.GetMinimum(key1)
	assert.Nil(t, err)
	assert.Equal(t, int64(-30), min)

	assert.Nil(t, a.Insert(key2, uint8(40)))
	min, err = a.GetMinimum(key2)
	assert.Nil(t, err)
	assert.Equal(t, int64(40), min)
}

func TestAggregator_GetSum(t *testing.T) {
	a := NewAggregator()
	key1 := "float64_1"
	key2 := "float64_2"

	sum, err := a.GetSum("badKey")
	assert.Equal(t, syscall.ENOENT, err)

	assert.Nil(t, a.Insert(key1, float64(0.5)))
	assert.Nil(t, a.Insert(key1, 3.25))
	assert.Nil(t, a.Insert(key2, float64(-33)))

	sum, err = a.GetSum(key1)
	assert.Nil(t, err)
	assert.Equal(t, 3.75, sum)

	sum, err = a.GetSum(key2)
	assert.Nil(t, err)
	assert.Equal(t, -33., sum)
}

func TestAggregator_Insert(t *testing.T) {
	a := NewAggregator()

	assert.Equal(t, syscall.ENOTSUP, a.Insert("int", "0"))
	assert.Nil(t, a.Insert("int", int(1)))
	assert.Nil(t, a.Insert("int", int(2)))
	assert.Equal(t, syscall.EINVAL, a.Insert("int", int8(3)))
	assert.Equal(t, syscall.EINVAL, a.Insert("int", "4"))
}

func TestAggregator_NewAggregator(t *testing.T) {
	a := NewAggregator()

	assert.NotNil(t, a)
	assert.NotNil(t, a.db)
}
