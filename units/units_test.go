package units

import (
    "testing"
)

func TestUnitsToBinaryString(t *testing.T) {
    var act, exp string

    act = ToBinaryString(0, 0, "", "B")
    exp = "0B"
    if act != exp {
        t.Errorf("Actual: %s, Expected: %s\n", act, exp)
    }

    act = ToBinaryString(1023, 0, "", "B")
    exp = "1023B"
    if act != exp {
        t.Errorf("Actual: %s, Expected: %s\n", act, exp)
    }

    act = ToBinaryString(1024, -1, "", "B")
    exp = "1KiB"
    if act != exp {
        t.Errorf("Actual: %s, Expected: %s\n", act, exp)
    }

    act = ToBinaryString(3.5*1024*1024*1024*1024, -1, "", "B")
    exp = "3.5TiB"
    if act != exp {
        t.Errorf("Actual: %s, Expected: %s\n", act, exp)
    }
}

func TestUnitsToMetricString(t *testing.T) {
    var act, exp string

    act = ToMetricString(0.000001, 2, " ", "m")
    exp = "1.00 um"
    if act != exp {
        t.Errorf("Actual: %s, Expected: %s\n", act, exp)
    }

    act = ToMetricString(0.025, 3, " ", "s")
    exp = "25.000 ms"
    if act != exp {
        t.Errorf("Actual: %s, Expected: %s\n", act, exp)
    }

    act = ToMetricString(0, 0, " ", "m")
    exp = "0 m"
    if act != exp {
        t.Errorf("Actual: %s, Expected: %s\n", act, exp)
    }

    act = ToMetricString(1000, 0, "", "g")
    exp = "1kg"
    if act != exp {
        t.Errorf("Actual: %s, Expected: %s\n", act, exp)
    }

    act = ToMetricString(500123000, -1, "-", "W")
    exp = "500.123-MW"
    if act != exp {
        t.Errorf("Actual: %s, Expected: %s\n", act, exp)
    }

    act = ToMetricString(1048576, 5, "", "B")
    exp = "1.04858MB"
    if act != exp {
        t.Errorf("Actual: %s, Expected: %s\n", act, exp)
    }

    act = ToMetricString(-9020, -1, " ", "N")
    exp = "-9.02 kN"
    if act != exp {
        t.Errorf("Actual: %s, Expected: %s\n", act, exp)
    }
}

func TestUnitsToTimeString(t *testing.T) {
    var act, exp string

    act = ToTimeString(59)
    exp = "0.00:00:59.000000"
    if act != exp {
        t.Errorf("Actual: %s, Expected: %s\n", act, exp)
    }

    act = ToTimeString(60)
    exp = "0.00:01:00.000000"
    if act != exp {
        t.Errorf("Actual: %s, Expected: %s\n", act, exp)
    }

    act = ToTimeString(75)
    exp = "0.00:01:15.000000"
    if act != exp {
        t.Errorf("Actual: %s, Expected: %s\n", act, exp)
    }

    act = ToTimeString(86399)
    exp = "0.23:59:59.000000"
    if act != exp {
        t.Errorf("Actual: %s, Expected: %s\n", act, exp)
    }

    act = ToTimeString(86401.123456789)
    exp = "1.00:00:01.123457"
    if act != exp {
        t.Errorf("Actual: %s, Expected: %s\n", act, exp)
    }

    act = ToTimeString(867600.050)
    exp = "10.01:00:00.050000"
    if act != exp {
        t.Errorf("Actual: %s, Expected: %s\n", act, exp)
    }

    act = ToTimeString(2680200.000789)
    exp = "31.00:30:00.000789"
    if act != exp {
        t.Errorf("Actual: %s, Expected: %s\n", act, exp)
    }

    act = ToTimeString(31557600)
    exp = "365.06:00:00.000000"
    if act != exp {
        t.Errorf("Actual: %s, Expected: %s\n", act, exp)
    }
}
