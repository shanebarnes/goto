package units

import (
    "fmt"
    "math"
    "strconv"
)

type prefix struct {
    sym string
    val float64
}

var binaryPrefix = [...]prefix {
    { sym: "",   val: math.Pow(1024, 0) },
    { sym: "Ki", val: math.Pow(1024, 1) },
    { sym: "Mi", val: math.Pow(1024, 2) },
    { sym: "Gi", val: math.Pow(1024, 3) },
    { sym: "Ti", val: math.Pow(1024, 4) },
    { sym: "Pi", val: math.Pow(1024, 5) },
    { sym: "Ei", val: math.Pow(1024, 6) },
    { sym: "Zi", val: math.Pow(1024, 7) },
    { sym: "Yi", val: math.Pow(1024, 8) },
}

var metricPrefixLt1 = [...]prefix {
    { sym: "",  val: math.Pow(1000,  0) },
    { sym: "m", val: math.Pow(1000, -1) },
    { sym: "u", val: math.Pow(1000, -2) },
    { sym: "n", val: math.Pow(1000, -3) },
    { sym: "p", val: math.Pow(1000, -4) },
    { sym: "f", val: math.Pow(1000, -5) },
    { sym: "a", val: math.Pow(1000, -6) },
    { sym: "z", val: math.Pow(1000, -7) },
    { sym: "y", val: math.Pow(1000, -8) },
}

var metricPrefixGe1 = [...]prefix {
    { sym: "",  val: math.Pow(1000,  0) },
    { sym: "k", val: math.Pow(1000,  1) },
    { sym: "M", val: math.Pow(1000,  2) },
    { sym: "G", val: math.Pow(1000,  3) },
    { sym: "T", val: math.Pow(1000,  4) },
    { sym: "P", val: math.Pow(1000,  5) },
    { sym: "E", val: math.Pow(1000,  6) },
    { sym: "Z", val: math.Pow(1000,  7) },
    { sym: "Y", val: math.Pow(1000,  8) },
}

const (
    SEC_IN_SECOND float64 = 1
    SEC_IN_MINUTE         = SEC_IN_SECOND *  60
    SEC_IN_HOUR           = SEC_IN_MINUTE *  60
    SEC_IN_DAY            = SEC_IN_HOUR   *  24
    SEC_IN_WEEK           = SEC_IN_DAY    *   7
    SEC_IN_MONTH          = SEC_IN_DAY    *  30
    SEC_IN_YEAR           = SEC_IN_DAY    * 365.25
)

func ToBinaryString(number float64, precision int, separator, quantity string) string {
    var i int
    var prefix prefix
    var symbol string
    var sfactor float64 = 1
    n := math.Abs(number)

    if number < 0 {
        sfactor = -1
    }

    for i, prefix = range binaryPrefix {
        if n < prefix.val {
            if i > 0 {
                i = i - 1
            }
            break
        }
    }

    n = n / binaryPrefix[i].val
    symbol = binaryPrefix[i].sym

    return strconv.FormatFloat(sfactor * n, 'f', precision, 64) + separator + symbol + quantity
}

func ToMetricString(number float64, precision int, separator, quantity string) string {
    var i int
    var prefix prefix
    var symbol string
    var sfactor float64 = 1
    n := math.Abs(number)

    if number < 0 {
        sfactor = -1
    }

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

    return strconv.FormatFloat(sfactor * n, 'f', precision, 64) + separator + symbol + quantity
}

func ToTimeString(durationSec float64) string {
    n := math.Abs(durationSec)

    years := math.Floor(n / SEC_IN_YEAR)
    n = n - (years * SEC_IN_YEAR)

    months := math.Floor(n / SEC_IN_MONTH)
    n = n - (months * SEC_IN_MONTH)

    days := math.Floor(n / SEC_IN_DAY)
    n = n - (days * SEC_IN_DAY)

    hours := math.Floor(n / SEC_IN_HOUR)
    n = n - (hours * SEC_IN_HOUR)

    minutes := math.Floor(n / SEC_IN_MINUTE)
    n = n - (minutes * SEC_IN_MINUTE)

    seconds := n

    return fmt.Sprintf("%03.0fy:%02.0fm:%02.0fd:%02.0fh:%02.0fm:%09fs",
                       years,
                       months,
                       days,
                       hours,
                       minutes,
                       seconds)
}
