package units

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitsToBinaryString(t *testing.T) {
	var str string

	str = ToBinaryString(0, 0, "", "B")
	assert.Equal(t, "0B", str)

	str = ToBinaryString(1023, 0, "", "B")
	assert.Equal(t, "1023B", str)

	str = ToBinaryString(1024, -1, "", "B")
	assert.Equal(t, "1KiB", str)

	str = ToBinaryString(3.5*1024*1024*1024*1024, -1, "", "B")
	assert.Equal(t, "3.5TiB", str)

	str = ToBinaryString(1048575, 3, "", "B")
	assert.Equal(t, "1023.999KiB", str)

	// FIXME: Should return 1.00MiB
	str = ToBinaryString(1048575, 2, "", "B")
	assert.Equal(t, "1024.00KiB", str)
}

func TestUnitsToBinaryStringWithPrefix(t *testing.T) {
	var str string
	val := 123456789.

	str = ToBinaryStringWithPrefix(val, 3, " ", "wrong", "B")
	assert.Equal(t, "117.738 MiB", str)

	str = ToBinaryStringWithPrefix(val, 3, " ", "Ki", "B")
	assert.Equal(t, "120563.271 KiB", str)

	str = ToBinaryStringWithPrefix(val, 3, " ", "Mi", "B")
	assert.Equal(t, "117.738 MiB", str)

	str = ToBinaryStringWithPrefix(val, 3, " ", "Gi", "B")
	assert.Equal(t, "0.115 GiB", str)
}

func TestUnitsToMetricString(t *testing.T) {
	var str string

	str = ToMetricString(0.000001, 2, " ", "m")
	assert.Equal(t, "1.00 um", str)

	str = ToMetricString(0.025, 3, " ", "s")
	assert.Equal(t, "25.000 ms", str)

	str = ToMetricString(0, 0, " ", "m")
	assert.Equal(t, "0 m", str)

	str = ToMetricString(1000, 0, "", "g")
	assert.Equal(t, "1kg", str)

	str = ToMetricString(500123000, -1, "-", "W")
	assert.Equal(t, "500.123-MW", str)

	str = ToMetricString(1048576, 5, "", "B")
	assert.Equal(t, "1.04858MB", str)

	str = ToMetricString(-9020, -1, " ", "N")
	assert.Equal(t, "-9.02 kN", str)

	str = ToMetricString(999999, -1, "", "B")
	assert.Equal(t, "999.999kB", str)

	// FIXME: Should return 1.00MB
	str = ToMetricString(999999, 2, "", "B")
	assert.Equal(t, "1000.00kB", str)
}

func TestUnitsToMetricStringWithPrefix(t *testing.T) {
	var str string
	val := 123456789.

	str = ToMetricStringWithPrefix(val, 3, " ", "wrong", "W")
	assert.Equal(t, "123.457 MW", str)

	str = ToMetricStringWithPrefix(val, 3, " ", "m", "W")
	assert.Equal(t, "123456789000.000 mW", str)

	str = ToMetricStringWithPrefix(val, 3, " ", "u", "W")
	assert.Equal(t, "123456789000000.000 uW", str)

	str = ToMetricStringWithPrefix(val, 3, " ", "n", "W")
	assert.Equal(t, "123456789000000000.000 nW", str)

	str = ToMetricStringWithPrefix(val, 3, " ", "k", "W")
	assert.Equal(t, "123456.789 kW", str)

	str = ToMetricStringWithPrefix(val, 3, " ", "M", "W")
	assert.Equal(t, "123.457 MW", str)

	str = ToMetricStringWithPrefix(val, 3, " ", "G", "W")
	assert.Equal(t, "0.123 GW", str)

	str = ToMetricStringWithPrefix(999999, 2, " ", "M", "B")
	assert.Equal(t, "1.00 MB", str)

	str = ToMetricStringWithPrefix(994999, 2, " ", "M", "B")
	assert.Equal(t, "0.99 MB", str)
}

func TestUnitsToNumber(t *testing.T) {
	var f float64
	var err error

	f, err = ToNumber("123")
	assert.Nil(t, err)
	assert.Equal(t, 123., f)

	f, err = ToNumber("1.048576M")
	assert.Nil(t, err)
	assert.Equal(t, 1048576., f)

	f, err = ToNumber("1.048576 M")
	assert.Nil(t, err)
	assert.Equal(t, 1048576., f)

	f, err = ToNumber("1.048576m")
	assert.Nil(t, err)
	assert.Equal(t, 0.001048576, f)

	f, err = ToNumber("1.0A")
	assert.NotNil(t, err)
	assert.Equal(t, 1., f)
}

func TestUnitsToTimeString(t *testing.T) {
	var str string

	str = ToTimeString(59)
	assert.Equal(t, "0.00:00:59.000000", str)

	str = ToTimeString(60)
	assert.Equal(t, "0.00:01:00.000000", str)

	str = ToTimeString(75)
	assert.Equal(t, "0.00:01:15.000000", str)

	str = ToTimeString(86399)
	assert.Equal(t, "0.23:59:59.000000", str)

	str = ToTimeString(86401.123456789)
	assert.Equal(t, "1.00:00:01.123457", str)

	str = ToTimeString(867600.050)
	assert.Equal(t, "10.01:00:00.050000", str)

	str = ToTimeString(2680200.000789)
	assert.Equal(t, "31.00:30:00.000789", str)

	str = ToTimeString(31557600)
	assert.Equal(t, "365.06:00:00.000000", str)
}
