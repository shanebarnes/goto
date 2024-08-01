package units

import (
	"fmt"
	"math"
	"testing"

	"github.com/dustin/go-humanize"
	"github.com/stretchr/testify/assert"
)

func TestToBinaryString(t *testing.T) {
	tests := []struct {
		expectedStr string
		number      float64
		precision   int
		quantity    string
		separator   string
	}{
		{expectedStr: "0B", number: 0, precision: 0, quantity: "B", separator: ""},
		{expectedStr: "1023B", number: humanize.KiByte - 1, precision: 0, quantity: "B", separator: ""},
		{expectedStr: "1KiB", number: humanize.KiByte, precision: 0, quantity: "B", separator: ""},
		{expectedStr: "1023.999KiB", number: humanize.MiByte - 1, precision: 3, quantity: "B", separator: ""},
		{expectedStr: "1MiB", number: humanize.MiByte, precision: 0, quantity: "B", separator: ""},
		{expectedStr: "1023.999999MiB", number: humanize.GiByte - 1, precision: 6, quantity: "B", separator: ""},
		{expectedStr: "1GiB", number: humanize.GiByte, precision: 0, quantity: "B", separator: ""},
		{expectedStr: "1023.999999999GiB", number: humanize.TiByte - 1, precision: 9, quantity: "B", separator: ""},
		{expectedStr: "1TiB", number: humanize.TiByte, precision: 0, quantity: "B", separator: ""},
		{expectedStr: "1023.999999999999TiB", number: humanize.PiByte - 1, precision: 12, quantity: "B", separator: ""},
		{expectedStr: "1PiB", number: humanize.PiByte, precision: 0, quantity: "B", separator: ""},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("number=%f,precision=%d", test.number, test.precision), func(t *testing.T) {
			assert.Equal(t, test.expectedStr, ToBinaryString(test.number, test.precision, test.separator, test.quantity))
		})
	}
}

func TestToBinaryStringWithPrefix(t *testing.T) {
	tests := []struct {
		expectedStr string
		number      float64
		precision   int
		prefix      string
		quantity    string
		separator   string
	}{
		{expectedStr: "117.738 MiB", number: 123456789, precision: 3, prefix: "wrong", quantity: "B", separator: " "},
		{expectedStr: "120563.271 KiB", number: 123456789, precision: 3, prefix: "Ki", quantity: "B", separator: " "},
		{expectedStr: "117.738 MiB", number: 123456789, precision: 3, prefix: "Mi", quantity: "B", separator: " "},
		{expectedStr: "0.115 GiB", number: 123456789, precision: 3, prefix: "Gi", quantity: "B", separator: " "},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("number=%f,precision=%d", test.number, test.precision), func(t *testing.T) {
			assert.Equal(t, test.expectedStr, ToBinaryStringWithPrefix(test.number, test.precision, test.separator, test.prefix, test.quantity))
		})
	}
}

func TestToMetricString(t *testing.T) {
	tests := []struct {
		expectedStr string
		number      float64
		precision   int
		quantity    string
		separator   string
	}{
		{expectedStr: "0B", number: 0, precision: 0, quantity: "B", separator: ""},
		{expectedStr: "999B", number: humanize.KByte - 1, precision: 0, quantity: "B", separator: ""},
		{expectedStr: "1kB", number: humanize.KByte, precision: 0, quantity: "B", separator: ""},
		{expectedStr: "999.999kB", number: humanize.MByte - 1, precision: 3, quantity: "B", separator: ""},
		{expectedStr: "1MB", number: humanize.MByte, precision: 0, quantity: "B", separator: ""},
		{expectedStr: "999.999999MB", number: humanize.GByte - 1, precision: 6, quantity: "B", separator: ""},
		{expectedStr: "1GB", number: humanize.GByte, precision: 0, quantity: "B", separator: ""},
		{expectedStr: "999.999999999GB", number: humanize.TByte - 1, precision: 9, quantity: "B", separator: ""},
		{expectedStr: "1TB", number: humanize.TByte, precision: 0, quantity: "B", separator: ""},
		{expectedStr: "999.999999999999TB", number: humanize.PByte - 1, precision: 12, quantity: "B", separator: ""},
		{expectedStr: "1PB", number: humanize.PByte, precision: 0, quantity: "B", separator: ""},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("number=%f,precision=%d", test.number, test.precision), func(t *testing.T) {
			assert.Equal(t, test.expectedStr, ToMetricString(test.number, test.precision, test.separator, test.quantity))
		})
	}
}

func TestToMetricStringWithPrefix(t *testing.T) {
	tests := []struct {
		expectedStr string
		number      float64
		precision   int
		prefix      string
		quantity    string
		separator   string
	}{
		{expectedStr: "123.457 MB", number: 123456789, precision: 3, prefix: "wrong", quantity: "B", separator: " "},
		{expectedStr: "123456.789 kB", number: 123456789, precision: 3, prefix: "k", quantity: "B", separator: " "},
		{expectedStr: "123.457 MB", number: 123456789, precision: 3, prefix: "M", quantity: "B", separator: " "},
		{expectedStr: "0.123 GB", number: 123456789, precision: 3, prefix: "G", quantity: "B", separator: " "},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("number=%f,precision=%d", test.number, test.precision), func(t *testing.T) {
			assert.Equal(t, test.expectedStr, ToMetricStringWithPrefix(test.number, test.precision, test.separator, test.prefix, test.quantity))
		})
	}
}

func TestUnitsToNumber(t *testing.T) {
	tests := []struct {
		expectErr   bool
		expectedNum float64
		str         string
	}{
		{expectErr: true, expectedNum: 0, str: "123A"},
		{expectErr: false, expectedNum: 123 * math.Pow(10, -6), str: "123μ"},
		{expectErr: false, expectedNum: 123 * math.Pow(10, -6), str: "123 μ"},
		{expectErr: false, expectedNum: 123 * math.Pow(10, -3), str: "123m"},
		{expectErr: false, expectedNum: 123 * math.Pow(10, -3), str: "123 m"},
		{expectErr: false, expectedNum: 123, str: "123"},
		{expectErr: false, expectedNum: 1000, str: "1k"},
		{expectErr: false, expectedNum: 1000, str: "1 k"},
		{expectErr: false, expectedNum: 1024, str: "1Ki"},
		{expectErr: false, expectedNum: 1024, str: "1 Ki"},
		{expectErr: false, expectedNum: 1234, str: "1.234k"},
		{expectErr: false, expectedNum: 1234, str: "1.234 k"},
		{expectErr: false, expectedNum: 1263.616, str: "1.234Ki"},
		{expectErr: false, expectedNum: 1263.616, str: "1.234 Ki"},
	}

	for _, test := range tests {
		t.Run(test.str, func(t *testing.T) {
			f, err := ToNumber(test.str)
			if assert.Equal(t, test.expectErr, err != nil) && err == nil {
				assert.Equal(t, test.expectedNum, f)
			}
		})
	}
}

func TestToTimeString(t *testing.T) {
	tests := []struct {
		durationInSeconds float64
		expectedStr       string
	}{
		{durationInSeconds: 0, expectedStr: "0.00:00:00.000000000"},
		{durationInSeconds: 59.999999999, expectedStr: "0.00:00:59.999999999"},
		{durationInSeconds: 60, expectedStr: "0.00:01:00.000000000"},
		{durationInSeconds: 3599.999999999, expectedStr: "0.00:59:59.999999999"},
		{durationInSeconds: 3600, expectedStr: "0.01:00:00.000000000"},
		{durationInSeconds: 86399.999999999, expectedStr: "0.23:59:59.999999999"},
		{durationInSeconds: 86400, expectedStr: "1.00:00:00.000000000"},
		{durationInSeconds: 86400 * 365, expectedStr: "365.00:00:00.000000000"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%013.010fs", test.durationInSeconds), func(t *testing.T) {
			assert.Equal(t, test.expectedStr, ToTimeString(test.durationInSeconds))
		})
	}
}
