package astral

import (
	"math"
	"time"
)

func julianday(date time.Time) float64 {
	date = date.UTC()
	// Calculate the Julian Day for the specified date//
	var (
		y = float64(date.Year())
		m = float64(date.Month())
		d = float64(date.Day())
	)

	if m <= 2 {
		y -= 1
		m += 12
	}

	a := math.Floor(y / 100)
	b := 2 - a + math.Floor(a/4)
	jd := math.Floor(365.25*(y+4716)) + math.Floor(30.6001*(m+1)) + d + b - 1524.5

	return jd
}

// Convert a Julian Day number to a Julian Century//
func jday_to_jcentury(julianday float64) float64 {
	return (julianday - 2451545.0) / 36525.0
}

// Convert a Julian Century number to a Julian Day//
func jcentury_to_jday(juliancentury float64) float64 {
	return (juliancentury * 36525.0) + 2451545.0
}
