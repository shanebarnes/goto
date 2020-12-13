package aggregator

import (
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var unityTypes = []interface{}{
	float32(1), float64(1),
	int(1), int8(1), int16(1), int32(1), int64(1),
	uint(1), uint8(1), uint16(1), uint32(1), uint64(1)}

func TestAggregator_Delete_InvalidKey(t *testing.T) {
	a := NewAggregator()
	assert.Equal(t, syscall.ENOENT, a.Delete("invalid"))
}

func TestAggregator_Delete_ValidKey(t *testing.T) {
	a := NewAggregator()

	for i, typ := range unityTypes {
		assert.Nil(t, a.Insert(i, typ))
	}

	for i, _ := range unityTypes {
		assert.Nil(t, a.Delete(i))
	}
}

func TestAggregator_GetAverage_InvalidKey(t *testing.T) {
	a := NewAggregator()
	avg, err := a.GetAverage("invalid")
	assert.Nil(t, avg)
	assert.Equal(t, syscall.ENOENT, err)
}

func TestAggregator_GetAverage_ValidKey(t *testing.T) {
	a := NewAggregator()

	for i, typ := range unityTypes {
		assert.Nil(t, a.Insert(i, typ))
		assert.Nil(t, a.Insert(i, typ))
		assert.Nil(t, a.Insert(i, typ))

		a.db[i].count -= 1 // Mock the insert count
	}

	for i, typ := range unityTypes {
		switch typ.(type) {
		case float32, float64:
			avg, err := a.GetAverage(i)
			assert.Equal(t, float64(1.5), avg)
			assert.Nil(t, err)
		case int, int8, int16, int32, int64:
			avg, err := a.GetAverage(i)
			assert.Equal(t, int64(1), avg)
			assert.Nil(t, err)
		case uint, uint8, uint16, uint32, uint64:
			avg, err := a.GetAverage(i)
			assert.Equal(t, uint64(1), avg)
			assert.Nil(t, err)
		default:
			assert.Fail(t, "unsupported type")
		}
	}
}

func TestAggregator_GetCount_InvalidKey(t *testing.T) {
	a := NewAggregator()
	cnt, err := a.GetCount("invalid")
	assert.Equal(t, int64(-1), cnt)
	assert.Equal(t, syscall.ENOENT, err)
}

func TestAggregator_GetCount_ValidKey(t *testing.T) {
	a := NewAggregator()
	for i, typ := range unityTypes {
		for j := 0; j <= i; j++ {
			assert.Nil(t, a.Insert(i, typ))
		}
	}

	for i, _ := range unityTypes {
		cnt, err := a.GetCount(i)
		assert.Equal(t, int64(i+1), cnt)
		assert.Nil(t, err)
	}
}

func TestAggregator_GetDuration_InvalidKey(t *testing.T) {
	a := NewAggregator()
	dur, err := a.GetDuration("invalid")
	assert.Equal(t, time.Duration(0), dur)
	assert.Equal(t, syscall.ENOENT, err)
}

func TestAggregator_GetDuration_MultiInsert(t *testing.T) {
	a := NewAggregator()
	for i, typ := range unityTypes {
		assert.Nil(t, a.Insert(i, typ))
		assert.Nil(t, a.Insert(i, typ))
	}

	for i, _ := range unityTypes {
		dur, err := a.GetDuration(i)
		assert.True(t, dur > time.Duration(0))
		assert.Nil(t, err)
	}
}

func TestAggregator_GetDuration_SingleInsert(t *testing.T) {
	a := NewAggregator()
	for i, typ := range unityTypes {
		assert.Nil(t, a.Insert(i, typ))
	}

	for i, _ := range unityTypes {
		dur, err := a.GetDuration(i)
		assert.Equal(t, time.Duration(0), dur)
		assert.Nil(t, err)
	}
}

func TestAggregator_GetMaximum_InvalidKey(t *testing.T) {
	a := NewAggregator()
	max, err := a.GetMaximum("invalid")
	assert.Nil(t, max)
	assert.Equal(t, syscall.ENOENT, err)
}

func TestAggregator_GetMaximum_ValidKey(t *testing.T) {
	a := NewAggregator()
	for i, typ := range unityTypes {
		var val1, val2, val3 interface{}
		switch typ.(type) {
		case float32:
			val1 = float32(-2)
			val2 = float32(-3)
			val3 = float32(3)
		case float64:
			val1 = float64(-2)
			val2 = float64(-3)
			val3 = float64(3)
		case int:
			val1 = -2
			val2 = -3
			val3 = 3
		case int8:
			val1 = int8(-2)
			val2 = int8(-3)
			val3 = int8(3)
		case int16:
			val1 = int16(-2)
			val2 = int16(-3)
			val3 = int16(3)
		case int32:
			val1 = int32(-2)
			val2 = int32(-3)
			val3 = int32(3)
		case int64:
			val1 = int64(-2)
			val2 = int64(-3)
			val3 = int64(3)
		case uint:
			val1 = uint(2)
			val2 = uint(1)
			val3 = uint(3)
		case uint8:
			val1 = uint8(2)
			val2 = uint8(1)
			val3 = uint8(3)
		case uint16:
			val1 = uint16(2)
			val2 = uint16(1)
			val3 = uint16(3)
		case uint32:
			val1 = uint32(2)
			val2 = uint32(1)
			val3 = uint32(3)
		case uint64:
			val1 = uint64(2)
			val2 = uint64(1)
			val3 = uint64(3)
		default:
			assert.Fail(t, "unsupported type")
		}

		assert.Nil(t, a.Insert(i, val1))
		max, err := a.GetMaximum(i)
		assert.Equal(t, val1, max)
		assert.Nil(t, err)

		assert.Nil(t, a.Insert(i, val2))
		max, err = a.GetMaximum(i)
		assert.Equal(t, val1, max)
		assert.Nil(t, err)

		assert.Nil(t, a.Insert(i, val3))
		max, err = a.GetMaximum(i)
		assert.Equal(t, val3, max)
		assert.Nil(t, err)
	}
}

func TestAggregator_GetMinimum_InvalidKey(t *testing.T) {
	a := NewAggregator()
	min, err := a.GetMinimum("invalid")
	assert.Nil(t, min)
	assert.Equal(t, syscall.ENOENT, err)
}

func TestAggregator_GetMinimum_ValidKey(t *testing.T) {
	a := NewAggregator()
	for i, typ := range unityTypes {
		var val1, val2, val3 interface{}
		switch typ.(type) {
		case float32:
			val1 = float32(-2)
			val2 = float32(1)
			val3 = float32(-3)
		case float64:
			val1 = float64(-2)
			val2 = float64(1)
			val3 = float64(-3)
		case int:
			val1 = -2
			val2 = 1
			val3 = -3
		case int8:
			val1 = int8(-2)
			val2 = int8(1)
			val3 = int8(-3)
		case int16:
			val1 = int16(-2)
			val2 = int16(1)
			val3 = int16(-3)
		case int32:
			val1 = int32(-2)
			val2 = int32(1)
			val3 = int32(-3)
		case int64:
			val1 = int64(-2)
			val2 = int64(1)
			val3 = int64(-3)
		case uint:
			val1 = uint(2)
			val2 = uint(3)
			val3 = uint(1)
		case uint8:
			val1 = uint8(2)
			val2 = uint8(3)
			val3 = uint8(1)
		case uint16:
			val1 = uint16(2)
			val2 = uint16(3)
			val3 = uint16(1)
		case uint32:
			val1 = uint32(2)
			val2 = uint32(3)
			val3 = uint32(1)
		case uint64:
			val1 = uint64(2)
			val2 = uint64(3)
			val3 = uint64(1)
		default:
			assert.Fail(t, "unsupported type")
		}

		assert.Nil(t, a.Insert(i, val1))
		min, err := a.GetMinimum(i)
		assert.Equal(t, val1, min)
		assert.Nil(t, err)

		assert.Nil(t, a.Insert(i, val2))
		min, err = a.GetMinimum(i)
		assert.Equal(t, val1, min)
		assert.Nil(t, err)

		assert.Nil(t, a.Insert(i, val3))
		min, err = a.GetMinimum(i)
		assert.Equal(t, val3, min)
		assert.Nil(t, err)
	}
}

func TestAggregator_GetRate_DurationNegative(t *testing.T) {
	for _, typ := range unityTypes {
		a := NewAggregator()
		assert.Nil(t, a.Insert(typ, typ))
		assert.Nil(t, a.Insert(typ, typ))

		a.db[typ].timeN = a.db[typ].time0.Add(time.Millisecond*500) // Mock the elapsed time

		rate, err := a.GetRate(typ, -time.Second)
		switch typ.(type) {
		case float32, float64:
			assert.Equal(t, float64(-4), rate)
			assert.Equal(t, nil, err)
		case int, int8, int16, int32, int64:
			assert.Equal(t, int64(-4), rate)
			assert.Equal(t, nil, err)
		case uint, uint8, uint16, uint32, uint64:
			assert.Equal(t, nil, rate)
			assert.Equal(t, syscall.EINVAL, err)
		default:
			assert.Fail(t, "unsupported type")
		}
	}
}

func TestAggregator_GetRate_DurationPositive(t *testing.T) {
	for _, typ := range unityTypes {
		a := NewAggregator()
		assert.Nil(t, a.Insert(typ, typ))
		assert.Nil(t, a.Insert(typ, typ))

		a.db[typ].timeN = a.db[typ].time0.Add(time.Millisecond*500) // Mock the elapsed time

		rate, err := a.GetRate(typ, time.Second)
		switch typ.(type) {
		case float32, float64:
			assert.Equal(t, float64(4), rate)
			assert.Equal(t, nil, err)
		case int, int8, int16, int32, int64:
			assert.Equal(t, int64(4), rate)
			assert.Equal(t, nil, err)
		case uint, uint8, uint16, uint32, uint64:
			assert.Equal(t, uint64(4), rate)
			assert.Equal(t, nil, err)
		default:
			assert.Fail(t, "unsupported type")
		}
	}
}

func TestAggregator_GetRate_DurationZero(t *testing.T) {
	for _, typ := range unityTypes {
		a := NewAggregator()
		assert.Nil(t, a.Insert(typ, typ))
		assert.Nil(t, a.Insert(typ, typ))

		a.db[typ].timeN = a.db[typ].time0.Add(time.Millisecond*500) // Mock the elapsed time

		rate, err := a.GetRate(typ, 0)
		switch typ.(type) {
		case float32, float64:
			assert.Equal(t, float64(0), rate)
			assert.Equal(t, nil, err)
		case int, int8, int16, int32, int64:
			assert.Equal(t, int64(0), rate)
			assert.Equal(t, nil, err)
		case uint, uint8, uint16, uint32, uint64:
			assert.Equal(t, uint64(0), rate)
			assert.Equal(t, nil, err)
		default:
			assert.Fail(t, "unsupported type")
		}
	}
}

func TestAggregator_GetRate_InvalidKey(t *testing.T) {
	a := NewAggregator()
	rate, err := a.GetRate("invalid", time.Second)
	assert.Equal(t, nil, rate)
	assert.Equal(t, syscall.ENOENT, err)
}

func TestAggregator_GetRate_SingleInsert(t *testing.T) {
	for _, typ := range unityTypes {
		a := NewAggregator()
		assert.Nil(t, a.Insert(typ, typ))
		rate, err := a.GetRate(typ, time.Second)
		switch typ.(type) {
		case float32, float64:
			assert.Equal(t, float64(0), rate)
			assert.Equal(t, nil, err)
		case int, int8, int16, int32, int64:
			assert.Equal(t, int64(0), rate)
			assert.Equal(t, nil, err)
		case uint, uint8, uint16, uint32, uint64:
			assert.Equal(t, uint64(0), rate)
			assert.Equal(t, nil, err)
		default:
			assert.Fail(t, "unsupported type")
		}
	}
}

func TestAggregator_GetSum_InvalidKey(t *testing.T) {
	a := NewAggregator()
	sum, err := a.GetSum("invalid")
	assert.Equal(t, nil, sum)
	assert.Equal(t, syscall.ENOENT, err)
}

func TestAggregator_GetSum_ValidKey(t *testing.T) {
	a := NewAggregator()
	for i, typ := range unityTypes {
		assert.Nil(t, a.Insert(i, typ))
		switch typ.(type) {
		case float32, float64:
			sum, err := a.GetSum(i)
			assert.Equal(t, float64(1), sum)
			assert.Nil(t, err)
		case int, int8, int16, int32, int64:
			sum, err := a.GetSum(i)
			assert.Equal(t, int64(1), sum)
			assert.Nil(t, err)
		case uint, uint8, uint16, uint32, uint64:
			sum, err := a.GetSum(i)
			assert.Equal(t, uint64(1), sum)
			assert.Nil(t, err)
		default:
			assert.Fail(t, "unsupported type")
		}
	}

	for i, typ := range unityTypes {
		assert.Nil(t, a.Insert(i, typ))
		switch typ.(type) {
		case float32, float64:
			sum, err := a.GetSum(i)
			assert.Equal(t, float64(2), sum)
			assert.Nil(t, err)
		case int, int8, int16, int32, int64:
			sum, err := a.GetSum(i)
			assert.Equal(t, int64(2), sum)
			assert.Nil(t, err)
		case uint, uint8, uint16, uint32, uint64:
			sum, err := a.GetSum(i)
			assert.Equal(t, uint64(2), sum)
			assert.Nil(t, err)
		default:
			assert.Fail(t, "unsupported type")
		}
	}
}

func TestAggregator_Insert_InconsistentType(t *testing.T) {
	a := NewAggregator()
	for i, typ := range unityTypes {
		assert.Nil(t, a.Insert(i, typ))
	}

	for i, typ := range unityTypes {
		j := (i + 1) % len(unityTypes)
		assert.Equal(t, syscall.EINVAL, a.Insert(j, typ))
	}
}

func TestAggregator_Insert_InvalidType(t *testing.T) {
	a := NewAggregator()
	assert.Equal(t, syscall.ENOTSUP, a.Insert("key", time.Duration(0)))
}

func TestAggregator_Insert_ValidType(t *testing.T) {
	a := NewAggregator()
	for i, typ := range unityTypes {
		assert.Nil(t, a.Insert(i, typ))
	}

	for i, typ := range unityTypes {
		assert.Nil(t, a.Insert(i, typ))
	}
}

func TestAggregator_NewAggregator(t *testing.T) {
	a := NewAggregator()

	assert.NotNil(t, a)
	assert.NotNil(t, a.db)
	assert.Equal(t, 0, len(a.db))
	assert.NotNil(t, a.mu)
}
