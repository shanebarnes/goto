package units

import (
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"unicode/utf8"
)

type prefixPair struct {
	sym string
	val float64
}

var binaryPrefix = [...]prefixPair{
	{sym: "", val: math.Pow(1024, 0)},
	{sym: "Ki", val: math.Pow(1024, 1)},
	{sym: "Mi", val: math.Pow(1024, 2)},
	{sym: "Gi", val: math.Pow(1024, 3)},
	{sym: "Ti", val: math.Pow(1024, 4)},
	{sym: "Pi", val: math.Pow(1024, 5)},
	{sym: "Ei", val: math.Pow(1024, 6)},
	{sym: "Zi", val: math.Pow(1024, 7)},
	{sym: "Yi", val: math.Pow(1024, 8)},
}

var metricPrefixLt1 = [...]prefixPair{
	{sym: "", val: math.Pow(1000, 0)},
	//{sym: "d", val: math.Pow(10, -1)},
	//{sym: "c", val: math.Pow(10, -2)},
	{sym: "m", val: math.Pow(10, -3)},
	{sym: "μ", val: math.Pow(10, -6)},
	{sym: "n", val: math.Pow(10, -9)},
	{sym: "p", val: math.Pow(10, -12)},
	{sym: "f", val: math.Pow(10, -15)},
	{sym: "a", val: math.Pow(10, -18)},
	{sym: "z", val: math.Pow(10, -21)},
	{sym: "y", val: math.Pow(10, -24)},
	{sym: "r", val: math.Pow(10, -27)},
	{sym: "q", val: math.Pow(10, -30)},
}

var metricPrefixGe1 = [...]prefixPair{
	{sym: "", val: math.Pow(10, 0)},
	//{sym: "da", val: math.Pow(10, 1)},
	//{sym: "h", val: math.Pow(10, 2)},
	{sym: "k", val: math.Pow(10, 3)},
	{sym: "M", val: math.Pow(10, 6)},
	{sym: "G", val: math.Pow(10, 9)},
	{sym: "T", val: math.Pow(10, 12)},
	{sym: "P", val: math.Pow(10, 15)},
	{sym: "E", val: math.Pow(10, 18)},
	{sym: "Z", val: math.Pow(10, 21)},
	{sym: "Y", val: math.Pow(10, 24)},
	{sym: "R", val: math.Pow(10, 27)},
	{sym: "Q", val: math.Pow(10, 30)},
}

var timePrefix = [...]prefixPair{
	{sym: "%.0f.", val: 86400},
	{sym: "%02.0f:", val: 3600},
	{sym: "%02.0f:", val: 60},
	{sym: "%012.09f", val: 1},
}

func getBinaryPrefixIndex(prefix string) int {
	for i, u := range binaryPrefix {
		if prefix == u.sym {
			return i
		}
	}

	return -1
}

func getMetricPrefixGe1Index(prefix string) int {
	for i, u := range metricPrefixGe1 {
		if prefix == u.sym {
			return i
		}
	}

	return -1
}

func getMetricPrefixLt1Index(prefix string) int {
	for i, u := range metricPrefixLt1 {
		if prefix == u.sym {
			return i
		}
	}

	return -1
}

func ToBinaryString(number float64, precision int, separator, quantity string) string {
	return ToBinaryStringWithPrefix(number, precision, separator, "-", quantity)
}

func ToBinaryStringWithPrefix(number float64, precision int, separator, returnPrefix, quantity string) string {
	var sfactor float64 = 1
	f := math.Abs(number)

	if number < 0 {
		sfactor = -1
	}

	// Convert to appropriate binary prefix that keeps unit value in the range [1, 1024)
	i := getBinaryPrefixIndex(returnPrefix)
	if i < 0 {
		var u prefixPair
		for i, u = range binaryPrefix {
			if f < u.val {
				if i > 0 {
					i = i - 1
				}
				break
			}
		}
	}

	f /= binaryPrefix[i].val
	symbol := binaryPrefix[i].sym

	return strconv.FormatFloat(sfactor*f, 'f', precision, 64) + separator + symbol + quantity
}

func ToMetricString(number float64, precision int, separator, quantity string) string {
	return ToMetricStringWithPrefix(number, precision, separator, "-", quantity)
}

func ToMetricStringWithPrefix(number float64, precision int, separator, returnPrefix, quantity string) string {
	var (
		i       int
		prefix  prefixPair
		symbol  string
		sfactor float64 = 1
	)
	n := math.Abs(number)

	if number < 0 {
		sfactor = -1
	}

	if i = getMetricPrefixGe1Index(returnPrefix); i > -1 { // Return desired metric prefix
		n /= metricPrefixGe1[i].val
		symbol = metricPrefixGe1[i].sym
	} else if i = getMetricPrefixLt1Index(returnPrefix); i > -1 { // Return desired metric prefix
		n /= metricPrefixLt1[i].val
		symbol = metricPrefixLt1[i].sym
	} else { // Convert to appropriate metric prefix that keeps unit value in the range [1, 1000)

		if n == 0 {
			symbol = ""
		} else if n < 1 {
			for i, prefix = range metricPrefixLt1 {
				if n >= prefix.val {
					break
				}
			}

			n /= metricPrefixLt1[i].val
			symbol = metricPrefixLt1[i].sym
		} else {
			for i, prefix = range metricPrefixGe1 {
				if n < prefix.val {
					if i > 0 {
						i = i - 1
					}
					break
				}
			}

			n /= metricPrefixGe1[i].val
			symbol = metricPrefixGe1[i].sym
		}
	}

	return strconv.FormatFloat(sfactor*n, 'f', precision, 64) + separator + symbol + quantity
}

func ToNumber(s string) (float64, error) {
	var (
		f      float64
		prefix string
	)

	n, err := fmt.Sscanf(s, "%f%s", &f, &prefix)
	if err != nil {
		switch {
		case errors.Is(err, io.EOF) && n == 1:
			return f, nil
		default:
			return 0, err
		}
	}

	if n == 2 {
		prefix = strings.TrimSpace(prefix)
		switch utf8.RuneCountInString(prefix) {
		case 0:
			return f, nil
		case 1: // Metric prefix (e.g., k)
			for _, u := range metricPrefixGe1 {
				if u.sym == prefix {
					return f * u.val, nil
				}
			}

			for _, u := range metricPrefixLt1 {
				if u.sym == prefix {
					return f * u.val, nil
				}
			}
		case 2: // Binary prefix (e.g., Ki)
			for _, u := range binaryPrefix {
				if u.sym == prefix {
					return f * u.val, nil
				}
			}
		}
		return 0, fmt.Errorf("invalid prefix: " + prefix)
	}

	return f, nil
}

func ToTimeString(durationInSeconds float64) string {
	durationInSeconds = math.Abs(durationInSeconds)

	var timeStr string

	for _, prefix := range timePrefix {
		f := durationInSeconds
		if prefix.val > 1 {
			f = math.Floor(f / prefix.val)
			durationInSeconds -= (f * prefix.val)
		}

		timeStr += fmt.Sprintf(prefix.sym, f)
	}

	return timeStr
}
