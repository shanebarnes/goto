package units

import (
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
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
	{sym: "m", val: math.Pow(1000, -1)},
	{sym: "u", val: math.Pow(1000, -2)},
	{sym: "n", val: math.Pow(1000, -3)},
	{sym: "p", val: math.Pow(1000, -4)},
	{sym: "f", val: math.Pow(1000, -5)},
	{sym: "a", val: math.Pow(1000, -6)},
	{sym: "z", val: math.Pow(1000, -7)},
	{sym: "y", val: math.Pow(1000, -8)},
}

var metricPrefixGe1 = [...]prefixPair{
	{sym: "", val: math.Pow(1000, 0)},
	{sym: "k", val: math.Pow(1000, 1)},
	{sym: "M", val: math.Pow(1000, 2)},
	{sym: "G", val: math.Pow(1000, 3)},
	{sym: "T", val: math.Pow(1000, 4)},
	{sym: "P", val: math.Pow(1000, 5)},
	{sym: "E", val: math.Pow(1000, 6)},
	{sym: "Z", val: math.Pow(1000, 7)},
	{sym: "Y", val: math.Pow(1000, 8)},
}

var timePrefix = [...]prefixPair{
	{sym: "%.0f.", val: 86400},
	{sym: "%02.0f:", val: 3600},
	{sym: "%02.0f:", val: 60},
	{sym: "%09f", val: 1},
}

func getBinaryPrefixIndex(prefix string) int {
	rv := -1

	for i, units := range binaryPrefix {
		if prefix == units.sym {
			rv = i
			break
		}
	}

	return rv
}

func getMetricPrefixGe1Index(prefix string) int {
	rv := -1

	for i, units := range metricPrefixGe1 {
		if prefix == units.sym {
			rv = i
			break
		}
	}

	return rv
}

func getMetricPrefixLt1Index(prefix string) int {
	rv := -1

	for i, units := range metricPrefixLt1 {
		if prefix == units.sym {
			rv = i
			break
		}
	}

	return rv
}

func ToNumber(s string) (float64, error) {
	var f float64 = 0
	var p string

	i, err := fmt.Sscanf(s, "%f%s", &f, &p)

	if i == 2 {
		p = strings.TrimSpace(p)
		found := false

		switch len(p) {
		case 0:
			found = true
		case 1: // Metric prefix (e.g., k)
			for _, units := range metricPrefixGe1 {
				if units.sym == p {
					f = f * units.val
					found = true
					break
				}
			}

			if !found {
				for _, units := range metricPrefixLt1 {
					if units.sym == p {
						f = f * units.val
						found = true
						break
					}
				}
			}
		case 2: // Binary prefix (e.g., Ki)
			for _, units := range binaryPrefix {
				if units.sym == p {
					f = f * units.val
					found = true
					break
				}
			}
		default:
		}

		if found == false {
			err = errors.New("Invalid prefix: " + p)
		}
	} else if i == 1 && err == io.EOF {
		err = nil
	}

	return f, err
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
		var units prefixPair
		for i, units = range binaryPrefix {
			if f < units.val {
				if i > 0 {
					i = i - 1
				}
				break
			}
		}
	}

	f = f / binaryPrefix[i].val
	symbol := binaryPrefix[i].sym

	return strconv.FormatFloat(sfactor*f, 'f', precision, 64) + separator + symbol + quantity
}

func ToBinaryString(number float64, precision int, separator, quantity string) string {
	return ToBinaryStringWithPrefix(number, precision, separator, "-", quantity)
}

func ToMetricStringWithPrefix(number float64, precision int, separator, returnPrefix, quantity string) string {
	var i int
	var prefix prefixPair
	var symbol string
	var sfactor float64 = 1
	n := math.Abs(number)

	if number < 0 {
		sfactor = -1
	}

	if i = getMetricPrefixGe1Index(returnPrefix); i > -1 { // Return desired metric prefix
		n = n / metricPrefixGe1[i].val
		symbol = metricPrefixGe1[i].sym
	} else if i = getMetricPrefixLt1Index(returnPrefix); i > -1 { // Return desired metric prefix
		n = n / metricPrefixLt1[i].val
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

			n = n / metricPrefixLt1[i].val
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

			n = n / metricPrefixGe1[i].val
			symbol = metricPrefixGe1[i].sym
		}
	}

	return strconv.FormatFloat(sfactor*n, 'f', precision, 64) + separator + symbol + quantity
}

func ToMetricString(number float64, precision int, separator, quantity string) string {
	return ToMetricStringWithPrefix(number, precision, separator, "-", quantity)
}

func ToTimeString(durationSec float64) string {
	n := math.Abs(durationSec)
	ret := ""

	for _, prefix := range timePrefix {
		f := n
		if prefix.val > 1 {
			f = math.Floor(f / prefix.val)
			n = n - (f * prefix.val)
		}

		ret = ret + fmt.Sprintf(prefix.sym, f)
	}

	return ret
}
