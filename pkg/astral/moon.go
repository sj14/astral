package astral

import (
	"fmt"
	"math"
	"time"
)

func properAngle(value float64) float64 {
	if value > 0.0 {
		value /= 360.0
		return (value - math.Floor(value)) * 360.0
	}

	tmp := math.Ceil(math.Abs(value / 360.0))
	return value + tmp*360.0
}

func phaseAsfloat(date time.Time) float64 {
	jd := julianday(date)
	DT := math.Pow((jd-2382148), 2) / (41048480 * 86400)
	T := (jd + DT - 2451545.0) / 36525
	T2 := math.Pow(T, 2)
	T3 := math.Pow(T, 3)
	D := 297.85 + (445267.1115 * T) - (0.0016300 * T2) + (T3 / 545868)
	D = radians(properAngle(D))
	M := 357.53 + (35999.0503 * T)
	M = radians(properAngle(M))
	M1 := 134.96 + (477198.8676 * T) + (0.0089970 * T2) + (T3 / 69699)
	M1 = radians(properAngle(M1))
	elong := degrees(D) + 6.29*math.Sin(M1)
	elong -= 2.10 * math.Sin(M)
	elong += 1.27 * math.Sin(2*D-M1)
	elong += 0.66 * math.Sin(2*D)
	elong = properAngle(elong)
	elong = math.Round(elong)
	moon := ((elong + 6.43) / 360) * 28
	return moon
}

// Calculates the phase of the moon on the specified date.
// Args:
//
//	date: The date to calculate the phase for. Dates are always in the UTC timezone.
//	      If not specified then today's date is used.
//
// Returns:
//
//	A number designating the phase.
//	============  ==============
//	0 .. 6.99     New moon
//	7 .. 13.99    First quarter
//	14 .. 20.99   Full moon
//	21 .. 27.99   Last quarter
//	============  ==============
func MoonPhase(date time.Time) float64 {
	moon := phaseAsfloat(date)
	if moon >= 28.0 {
		moon -= 28.0
	}
	return moon
}

// MoonPhaseDescription returns the description of the given moon phase.
func MoonPhaseDescription(x float64) (string, error) {
	if x < 0 || x >= 28 {
		return "", fmt.Errorf("%v is out of the expected range (0-27.99)", x)
	}
	if x < 7 {
		return "New Moon", nil
	}
	if x >= 7 && x < 14 {
		return "First Quarter", nil
	}
	if x >= 14 && x < 21 {
		return "Full Moon", nil
	}
	if x >= 21 && x < 28 {
		return "Laster Quarter", nil
	}

	return "", fmt.Errorf("failed parsing %v", x)
}
