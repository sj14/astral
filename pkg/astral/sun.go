package astral

import (
	"errors"
	"fmt"
	"math"
	"time"
)

// Using 32 arc minutes as sun's apparent diameter
const sunApperentRadius = 32.0 / (60.0 * 2.0)

func degrees(rad float64) float64 {
	return rad * (180 / math.Pi)
}

func radians(deg float64) float64 {
	return deg * (math.Pi / 180)
}

type SunDirection int

const (
	SunDirectionRising  SunDirection = 1
	SunDirectionSetting SunDirection = -1
)

const (
	DepressionCivil        float64 = 6.0
	DepressionNautical     float64 = 12.0
	DepressionAstronomical float64 = 18.0
)

type Observer struct {
	Latitude  float64
	Longitude float64
	Elevation float64
}

// Convert a floating point number of minutes to a time.Duration
func minutes_to_timedelta(minutes float64) time.Duration {
	nanoseconds := time.Duration(minutes * 60000000000)
	return nanoseconds
}

// Calculate the geometric mean longitude of the sun//
func geom_mean_long_sun(juliancentury float64) float64 {
	l0 := 280.46646 + juliancentury*(36000.76983+0.0003032*juliancentury)
	return math.Mod(l0, 360)
}

// Calculate the geometric mean anomaly of the sun//
func geom_mean_anomaly_sun(juliancentury float64) float64 {
	return 357.52911 + juliancentury*(35999.05029-0.0001537*juliancentury)
}

// Calculate the eccentricity of Earth's orbit//
func eccentric_location_earth_orbit(juliancentury float64) float64 {
	return 0.016708634 - juliancentury*(0.000042037+0.0000001267*juliancentury)
}

// Calculate the equation of the center of the sun//
func sun_eq_of_center(juliancentury float64) float64 {
	m := geom_mean_anomaly_sun(juliancentury)

	mrad := radians(m)
	sinm := math.Sin(mrad)
	sin2m := math.Sin(mrad + mrad)
	sin3m := math.Sin(mrad + mrad + mrad)

	c := sinm*(1.914602-juliancentury*(0.004817+0.000014*juliancentury)) + sin2m*(0.019993-0.000101*juliancentury) + sin3m*0.000289
	return c
}

// Calculate the sun's true longitude//
func sun_true_long(juliancentury float64) float64 {
	l0 := geom_mean_long_sun(juliancentury)
	c := sun_eq_of_center(juliancentury)

	return l0 + c
}

// Calculate the sun's true anomaly//
func sun_true_anomoly(juliancentury float64) float64 {
	m := geom_mean_anomaly_sun(juliancentury)
	c := sun_eq_of_center(juliancentury)

	return m + c
}

func sun_rad_vector(juliancentury float64) float64 {
	v := sun_true_anomoly(juliancentury)
	e := eccentric_location_earth_orbit(juliancentury)

	return (1.000001018 * (1 - e*e)) / (1 + e*math.Cos(radians(v)))
}

func sun_apparent_long(juliancentury float64) float64 {
	true_long := sun_true_long(juliancentury)

	omega := 125.04 - 1934.136*juliancentury
	return true_long - 0.00569 - 0.00478*math.Sin(radians(omega))
}

func mean_obliquity_of_ecliptic(juliancentury float64) float64 {
	seconds := 21.448 - juliancentury*(46.815+juliancentury*(0.00059-juliancentury*(0.001813)))
	return 23.0 + (26.0+(seconds/60.0))/60.0
}

func obliquity_correction(juliancentury float64) float64 {
	e0 := mean_obliquity_of_ecliptic(juliancentury)

	omega := 125.04 - 1934.136*juliancentury
	return e0 + 0.00256*math.Cos(radians(omega))
}

// Calculate the sun's right ascension
func sun_rt_ascension(juliancentury float64) float64 {
	oc := obliquity_correction(juliancentury)
	al := sun_apparent_long(juliancentury)

	tananum := math.Cos(radians(oc)) * math.Sin(radians(al))
	tanadenom := math.Cos(radians(al))
	return degrees(math.Atan2(tananum, tanadenom))
}

// Calculate the sun's declination
func sun_declination(juliancentury float64) float64 {
	e := obliquity_correction(juliancentury)
	lambd := sun_apparent_long(juliancentury)

	sint := math.Sin(radians(e)) * math.Sin(radians(lambd))
	return degrees(math.Asin(sint))
}

func var_y(juliancentury float64) float64 {
	epsilon := obliquity_correction(juliancentury)
	y := math.Tan(radians(epsilon) / 2.0)
	return y * y
}

func eq_of_time(juliancentury float64) float64 {
	l0 := geom_mean_long_sun(juliancentury)
	e := eccentric_location_earth_orbit(juliancentury)
	m := geom_mean_anomaly_sun(juliancentury)

	y := var_y(juliancentury)

	sin2l0 := math.Sin(2.0 * radians(l0))
	sinm := math.Sin(radians(m))
	cos2l0 := math.Cos(2.0 * radians(l0))
	sin4l0 := math.Sin(4.0 * radians(l0))
	sin2m := math.Sin(2.0 * radians(m))

	Etime := y*sin2l0 - 2.0*e*sinm + 4.0*e*y*sinm*cos2l0 - 0.5*y*y*sin4l0 - 1.25*e*e*sin2m

	return degrees(Etime) * 4.0
}

// Calculate the hour angle of the sun
//
// See https://en.wikipedia.org/wiki/Hour_angle#Solar_hour_angle
// Args:
//
//	latitude: The latitude of the obersver
//	declination: The declination of the sun
//	zenith: The zenith angle of the sun
//	direction: The direction of traversal of the sun
//
// Raises:
//
//	ValueError
func hour_angle(latitude float64, declination float64, zenith float64, direction SunDirection) (float64, error) {
	latitude_rad := radians(latitude)
	declination_rad := radians(declination)
	zenith_rad := radians(zenith)

	h := (math.Cos(zenith_rad) - math.Sin(latitude_rad)*math.Sin(declination_rad)) / (math.Cos(latitude_rad) * math.Cos(declination_rad))

	hourAngle := math.Acos(h)
	if math.IsNaN(hourAngle) {
		return 0, errors.New("not able to determine hour angle")
	}
	if direction == SunDirectionSetting {
		hourAngle = -hourAngle
	}
	return hourAngle, nil
}

// Calculate the extra degrees of depression that you can see round the earth
// due to the increase in elevation.
// Args:
//
//	elevation: Elevation above the earth in metres
//
// Returns:
//
//	A number of degrees to add to adjust for the elevation of the observer
func adjust_to_horizon(elevation float64) float64 {

	if elevation <= 0 {
		return 0
	}

	r := 6356900.0 // radius of the earth
	a1 := r
	h1 := r + elevation
	theta1 := math.Acos(a1 / h1)
	return degrees(theta1)
}

// Calculate the number of degrees to adjust for an obscuring feature
func adjust_to_obscuring_feature(elevation0, elevation1 float64) float64 {
	if elevation0 == 0.0 {
		return 0.0
	}

	sign := 1.0
	if elevation0 < 0.0 {
		sign = -1
	}

	return sign * degrees(math.Acos(math.Abs(elevation0)/math.Sqrt(math.Pow(elevation0, 2)+math.Pow(elevation1, 2))))
}

// Calculate the degrees of refraction of the sun due to the sun's elevation.
func refraction_at_zenith(zenith float64) float64 {

	elevation := 90 - zenith
	if elevation >= 85.0 {
		return 0
	}

	refractionCorrection := 0.0
	te := math.Tan(radians(elevation))
	if elevation > 5.0 {
		refractionCorrection = (58.1/te - 0.07/(te*te*te) + 0.000086/(te*te*te*te*te))
	} else if elevation > -0.575 {
		step1 := -12.79 + elevation*0.711
		step2 := 103.4 + elevation*step1
		step3 := -518.2 + elevation*step2
		refractionCorrection = 1735.0 + elevation*step3
	} else {
		refractionCorrection = -20.774 / te
	}
	refractionCorrection = refractionCorrection / 3600.0

	return refractionCorrection
}

// Calculate the time in the UTC timezone when the sun transits the specificed zenith
// Args:
//
//	observer: An observer viewing the sun at a specific, latitude, longitude and elevation
//	date: The date to calculate for
//	zenith: The zenith angle for which to calculate the transit time
//	direction: The direction that the sun is traversing
//
// Raises:
//
//	ValueError if the zenith is not transitted by the sun
//
// Returns:
//
//	the time when the sun transits the specificed zenith
func time_of_transit(observer Observer, date time.Time, zenith float64, direction SunDirection) (time.Time, error) {
	latitude := observer.Latitude
	if observer.Latitude > 89.8 {
		latitude = 89.8
	} else if observer.Latitude < -89.8 {
		latitude = -89.8
	}

	// TODO:
	// if isinstance(observer.elevation, float) && observer.elevation > 0.0 {
	// 	adjustment_for_elevation = adjust_to_horizon(observer.elevation)
	adjustment_for_elevation := 0.0
	if observer.Elevation > 0.0 {
		adjustment_for_elevation = adjust_to_horizon(observer.Elevation)
	}
	// } else if isinstance(observer.elevation, tuple) {
	// 	adjustment_for_elevation = adjust_to_obscuring_feature(observer.elevation)
	// }

	adjustment_for_refraction := refraction_at_zenith(zenith + adjustment_for_elevation)

	jd := julianday(date)
	jc := jday_to_jcentury(jd)
	solarDec := sun_declination(jc)

	hourangle, err := hour_angle(latitude, solarDec, zenith+adjustment_for_elevation-adjustment_for_refraction, direction)
	if err != nil {
		return time.Time{}, err
	}

	delta := -observer.Longitude - degrees(hourangle)
	timeDiff := 4.0 * delta
	timeUTC := 720.0 + timeDiff - eq_of_time(jc)

	jc = jday_to_jcentury(jcentury_to_jday(jc) + timeUTC/1440.0)
	solarDec = sun_declination(jc)
	hourangle, err = hour_angle(latitude, solarDec, zenith+adjustment_for_elevation+adjustment_for_refraction, direction)
	if err != nil {
		return time.Time{}, err
	}

	delta = -observer.Longitude - degrees(hourangle)
	timeDiff = 4.0 * delta
	timeUTC = 720 + timeDiff - eq_of_time(jc)

	td := minutes_to_timedelta(timeUTC)
	dt := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC).Add(td).In(date.Location())
	return dt, nil
}

// Calculates the time when the sun is at the specified elevation on the specified date.
// Note:
//
//	This method uses positive elevations for those above the horizon.
//	Elevations greater than 90 degrees are converted to a setting sun
//	i.e. an elevation of 110 will calculate a setting sun at 70 degrees.
//
// Args:
//
//	elevation: Elevation of the sun in degrees above the horizon to calculate for.
//	observer:  Observer to calculate for
//	date:      Date to calculate for.
//	direction: Determines whether the calculated time is for the sun rising or setting.
//	           Use ``SunDirectionRising`` or ``SunDirectionSetting``. Default is rising.
//
// Returns:
//
//	Date and time at which the sun is at the specified elevation.
func TimeAtElevation(observer Observer, elevation float64, date time.Time, direction SunDirection) (time.Time, error) {
	if elevation > 90.0 {
		elevation = 180.0 - elevation
		direction = SunDirectionSetting
	}

	zenith := 90 - elevation
	t, err := time_of_transit(observer, date, zenith, direction)
	if err != nil {
		return time.Time{}, fmt.Errorf("sun never reaches an elevation of %v degrees at this location", elevation)
	}
	return t, nil
}

// Calculate solar noon time when the sun is at its highest point.
// Args:
//
//	observer: An observer viewing the sun at a specific, latitude, longitude and elevation
//	date:     Date to calculate for. Default is today for the specified tzinfo.
//
// Returns:
//
//	Date and time at which noon occurs.
func Noon(observer Observer, date time.Time) time.Time {
	jc := jday_to_jcentury(julianday(date))
	eqtime := eq_of_time(jc)
	timeUTC := (720.0 - (4 * observer.Longitude) - eqtime) / 60.0

	hour := int(timeUTC)
	minute := int((timeUTC - float64(hour)) * 60)
	second := int((((timeUTC - float64(hour)) * 60.0) - float64(minute)) * 60)

	if second > 59 {
		second -= 60
		minute += 1
	} else if second < 0 {
		second += 60
		minute -= 1
	}
	if minute > 59 {
		minute -= 60
		hour += 1
	} else if minute < 0 {
		minute += 60
		hour -= 1
	}
	if hour > 23 {
		hour -= 24
		date = date.Add(24 * time.Hour)
	} else if hour < 0 {
		hour += 24
		date = date.Add(-24 * time.Hour)
	}
	noon := time.Date(date.Year(), date.Month(), date.Day(), hour, minute, second, 0, time.UTC).In(date.Location())
	return noon
}

// Calculate solar midnight time.
// Note:
//
//	This calculates the solar midnight that is closest
//	to 00:00:00 of the specified date i.e. it may return a time that is on
//	the previous day.
//
// Args:
//
//	observer: An observer viewing the sun at a specific, latitude, longitude and elevation
//	date:     Date to calculate for. Default is today for the specified tzinfo.
//
// Returns:
//
//	Date and time at which midnight occurs.
func Midnight(observer Observer, date time.Time) time.Time {
	date = time.Date(date.Year(), date.Month(), date.Day(), 12, 0, 0, 0, date.Location())
	jd := julianday(date)
	newt := jday_to_jcentury(jd + 0.5 + -observer.Longitude/360.0)

	eqtime := eq_of_time(newt)
	timeUTC := (-observer.Longitude * 4.0) - eqtime

	timeUTC = timeUTC / 60.0
	hour := int(timeUTC)
	minute := int((timeUTC - float64(hour)) * 60)
	second := int((((timeUTC - float64(hour)) * 60) - float64(minute)) * 60)

	if second > 59 {
		second -= 60
		minute += 1
	} else if second < 0 {
		second += 60
		minute -= 1
	}

	if minute > 59 {
		minute -= 60
		hour += 1
	} else if minute < 0 {
		minute += 60
		hour -= 1
	}

	if hour < 0 {
		hour += 24
		date = date.Add(-24 * time.Hour)
	}
	return time.Date(date.Year(), date.Month(), date.Day(), hour, minute, second, 0, time.UTC).In(date.Location())
}

func ZenithAndAzimuth(observer Observer, dateandtime time.Time, with_refraction bool) (float64, float64) {
	latitude := observer.Latitude

	if observer.Latitude > 89.8 {
		latitude = 89.8
	} else if observer.Latitude < -89.8 {
		latitude = -89.8
	}
	longitude := observer.Longitude

	utc_datetime := dateandtime.UTC()

	timenow := (utc_datetime.Hour() + (utc_datetime.Minute() / 60.0) + (utc_datetime.Second() / 3600.0))

	JD := julianday(dateandtime)
	t := jday_to_jcentury(JD + float64(timenow)/24.0)
	solarDec := sun_declination(t)
	eqtime := eq_of_time(t)

	solarTimeFix := eqtime - (4.0 * -longitude)
	trueSolarTime := float64(utc_datetime.Hour()*60+utc_datetime.Minute()+utc_datetime.Second()/60) + solarTimeFix
	//    in minutes as a float, fractional part is seconds

	for trueSolarTime > 1440 {
		trueSolarTime = trueSolarTime - 1440
	}

	hourangle := trueSolarTime/4.0 - 180.0
	//    Thanks to Louis Schwarzmayr for the next line:
	if hourangle < -180 {
		hourangle = hourangle + 360.0
	}

	harad := radians(hourangle)

	csz := math.Sin(radians(latitude))*math.Sin(radians(solarDec)) + math.Cos(radians(latitude))*math.Cos(radians(solarDec))*math.Cos(harad)

	if csz > 1.0 {
		csz = 1.0
	} else if csz < -1.0 {
		csz = -1.0
	}

	zenith := degrees(math.Acos(csz))

	azDenom := math.Cos(radians(latitude)) * math.Sin(radians(zenith))

	azimuth := 0.0
	if math.Abs(azDenom) > 0.001 {
		azRad := ((math.Sin(radians(latitude)) * math.Cos(radians(zenith))) - math.Sin(radians(solarDec))) / azDenom

		if math.Abs(azRad) > 1.0 {
			if azRad < 0 {
				azRad = -1.0
			} else {
				azRad = 1.0
			}
		}
		azimuth = 180.0 - degrees(math.Acos(azRad))

		if hourangle > 0.0 {
			azimuth = -azimuth
		}
	} else {
		if latitude > 0.0 {
			azimuth = 180.0
		} else {
			azimuth = 0.0
		}
	}

	if azimuth < 0.0 {
		azimuth = azimuth + 360.0
	}
	if with_refraction {
		zenith -= refraction_at_zenith(zenith)
	}
	return zenith, azimuth
}

// Calculate the zenith angle of the sun.
// Args:
//     observer:    Observer to calculate the solar zenith for
//     dateandtime: The date and time for which to calculate the angle.
//     with_refraction: If True adjust zenith to take refraction into account

// Returns:
//
//	The zenith angle in degrees.
func Zenith(observer Observer, dateandtime time.Time, with_refraction bool) float64 {
	zenith, _ := ZenithAndAzimuth(observer, dateandtime, with_refraction)
	return zenith
}

// Calculate the azimuth angle of the sun.
// Args:
//
//	observer:    Observer to calculate the solar azimuth for
//	dateandtime: The date and time for which to calculate the angle.
//
// Returns:
//
//	The azimuth angle in degrees clockwise from North.
func Azimuth(observer Observer, dateandtime time.Time) float64 {
	_, azimuth := ZenithAndAzimuth(observer, dateandtime, true)
	return azimuth
}

// Calculate the sun's angle of elevation.
// Args:
//
//	observer:    Observer to calculate the solar elevation for
//	dateandtime: The date and time for which to calculate the angle.
//	with_refraction: If True adjust elevation to take refraction into account
//
// Returns:
//
//	The elevation angle in degrees above the horizon.
func Elevation(observer Observer, dateandtime time.Time, with_refraction bool) float64 {
	return 90.0 - Zenith(observer, dateandtime, with_refraction)
}

// Calculate dawn time.
// Args:
//
//	observer:   Observer to calculate dawn for
//	date:       Date to calculate for.
//	depression: Number of degrees below the horizon to use to calculate dawn.
//	            Default is for Civil dawn i.e. 6.0
//	tzinfo:     Timezone to return times in. Default is UTC.
//
// Returns:
//
//	Date and time at which dawn occurs.
func Dawn(observer Observer, date time.Time, depression float64) (time.Time, error) {
	t, err := time_of_transit(observer, date, 90.0+depression, SunDirectionRising)
	if err != nil {
		return t, fmt.Errorf("sun never reaches %v degrees below the horizon, at this location", depression)
	}
	return t, nil
}

var (
	ErrAlwaysBelow = errors.New("sun is always below the horizon on this day, at this location")
	ErrAlwaysAbove = errors.New("sun is always above the horizon on this day, at this location")
)

// Calculate sunrise time.
// Args:
//
//	observer Observer to calculate sunrise for
//	date:     Date to calculate for.
//	tzinfo:   Timezone to return times in. Default is UTC.
//
// Returns:
//
//	Date and time at which sunrise occurs.
func Sunrise(observer Observer, date time.Time) (time.Time, error) {
	t, err := time_of_transit(observer, date, 90.0+sunApperentRadius, SunDirectionRising)

	if err != nil {
		z := Zenith(observer, Noon(observer, date), true)
		if z > 90.0 {
			return time.Time{}, ErrAlwaysBelow
		}
		return time.Time{}, ErrAlwaysAbove
	}

	return t, nil
}

// Calculate sunset time.
// Args:
//
//	observer Observer to calculate sunset for
//	date:     Date to calculate for.
//	tzinfo:   Timezone to return times in. Default is UTC.
//
// Returns:
//
//	Date and time at which sunset occurs.
//
// Raises:
//
//	    ValueError: if the sun does not reach the horizon
//
//		if isinstance(tzinfo, str) {
//			tzinfo = pytz.timezone(tzinfo)
//		}
//		if date.IsZero() {
//			date := today(tzinfo)
//		}
func Sunset(observer Observer, date time.Time) (time.Time, error) {
	t, err := time_of_transit(observer, date, 90.0+sunApperentRadius, SunDirectionSetting)
	if err != nil {
		z := Zenith(observer, Noon(observer, date), true)
		if z > 90.0 {
			return time.Time{}, ErrAlwaysBelow
		}
		return time.Time{}, ErrAlwaysAbove
	}
	return t, nil

}

// Calculate dusk time.

// Args:
//     observer:   Observer to calculate dusk for
//     date:       Date to calculate for.
//     depression: Number of degrees below the horizon to use to calculate dusk.
//                 Default is for Civil dusk i.e. 6.0
//     tzinfo:     Timezone to return times in. Default is UTC.

// Returns:
//     Date and time at which dusk occurs.

// Raises:
//     ValueError: if dusk does not occur on the specified date
//

//	if isinstance(tzinfo, str) {
//		tzinfo = pytz.timezone(tzinfo)
//	}
//
//	if date.IsZero() {
//		date := today(tzinfo)
//	}
func Dusk(observer Observer, date time.Time, depression float64) (time.Time, error) {
	t, err := time_of_transit(observer, date, 90.0+depression, SunDirectionSetting)
	if err != nil {
		return t, fmt.Errorf("sun never reaches %v degrees below the horizon, at this location", depression)
	}
	return t, nil
}

// Calculate daylight start and end times.
// Args:
//
//	observer Observer to calculate daylight for
//	date:     Date to calculate for.
//	tzinfo:   Timezone to return times in. Default is UTC.
//
// Returns:
//
//	A tuple of the date and time at which daylight starts and ends.
//
// Raises:
//
//	ValueError: if the sun does not rise or does not set
func Daylight(observer Observer, date time.Time) (time.Time, time.Time, error) {
	start, err := Sunrise(observer, date)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	end, err := Sunset(observer, date)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return start, end, nil
}

// Calculate night start and end times.
// Night is calculated to be between astronomical dusk on the
// date specified and astronomical dawn of the next day.
// Args:
//
//	observer Observer to calculate night for
//	date:     Date to calculate for. Default is today's date for the
//	          specified tzinfo.
//	tzinfo:   Timezone to return times in. Default is UTC.
//
// Returns:
//
//	A tuple of the date and time at which night starts and ends.
//
// Raises:
//
//	ValueError: if dawn does not occur on the specified date or
//	            dusk on the following day
func Night(observer Observer, date time.Time) (time.Time, time.Time, error) {
	start, err := Dusk(observer, date, 6)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	tomorrow := date.Add(24 * time.Hour)
	end, err := Dawn(observer, tomorrow, 6)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return start, end, nil
}

// Returns the start and end times of Twilight
// when the sun is traversing in the specified direction.
// This method defines twilight as being between the time
// when the sun is at -6 degrees and sunrise/sunset.
// Args:
//
//	observer:  Observer to calculate twilight for
//	date:      Date for which to calculate the times.
//
//	direction: Determines whether the time is for the sun rising or setting.
//	              Use ``astral.SunDirectionRising`` or ``astral.SunDirectionSetting``.
//	tzinfo:    Timezone to return times in. Default is UTC.
//
// Returns:
//
//	A tuple of the date and time at which twilight starts and ends.
//
// Raises:
//
//	ValueError: if the sun does not rise or does not set
func Twilight(observer Observer, date time.Time, direction SunDirection) (time.Time, time.Time, error) {
	start, err := time_of_transit(observer, date, 90+6, direction)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	end, err := Sunset(observer, date)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	if direction == SunDirectionRising {
		end, err := Sunrise(observer, date)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		return start, end, nil
	}
	return end, start, nil
}

// Returns the start and end times of the Golden Hour
// when the sun is traversing in the specified direction.
// This method uses the definition from PhotoPills i.e. the
// golden hour is when the sun is between 4 degrees below the horizon
// and 6 degrees above.
// Args:
//
//	observer:  Observer to calculate the golden hour for
//	date:      Date for which to calculate the times.
//
//	direction: Determines whether the time is for the sun rising or setting.
//	              Use ``SunDirectionRising`` or ``SunDirectionSetting``.
//	tzinfo:    Timezone to return times in. Default is UTC.
//
// Returns:
//
//	A tuple of the date and time at which the Golden Hour starts and ends.
//
// Raises:
//
//	ValueError: if the sun does not transit the elevations -4 & +6 degrees
func GoldenHour(observer Observer, date time.Time, direction SunDirection) (time.Time, time.Time, error) {
	start, err := time_of_transit(observer, date, 90+4, direction)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	end, err := time_of_transit(observer, date, 90-6, direction)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	if direction == SunDirectionRising {
		return start, end, nil
	}
	return end, start, nil
}

// Returns the start and end times of the Blue Hour
// when the sun is traversing in the specified direction.

// This method uses the definition from PhotoPills i.e. the
// blue hour is when the sun is between 6 and 4 degrees below the horizon.

// Args:
//     observer:  Observer to calculate the blue hour for
//     date:      Date for which to calculate the times.
//
//     direction: Determines whether the time is for the sun rising or setting.
//                   Use ``SunDirectionRising`` or ``SunDirectionSetting``.
//     tzinfo:    Timezone to return times in. Default is UTC.

// Returns:
//     A tuple of the date and time at which the Blue Hour starts and ends.

// Raises:
//
//	ValueError: if the sun does not transit the elevations -4 & -6 degrees
func BlueHour(observer Observer, date time.Time, direction SunDirection) (time.Time, time.Time, error) {
	start, err := time_of_transit(observer, date, 90+6, direction)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	end, err := time_of_transit(observer, date, 90+4, direction)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	if direction == SunDirectionRising {
		return start.In(date.Location()), end.In(date.Location()), nil
	}
	return end, start, nil

}

// TODO: not working correctly
// func Rahukaalam(observer Observer, date time.Time, daytime bool) (time.Time, time.Time, error) {
// 	// Calculate ruhakaalam times.

// 	// Args:
// 	//     observer Observer to calculate rahukaalam for
// 	//     date:     Date to calculate for.
// 	//     daytime:  If True calculate for the day time else calculate for the night time.
// 	//     tzinfo:   Timezone to return times in. Default is UTC.

// 	// Returns:
// 	//     Tuple containing the start and end times for Rahukaalam.

// 	// Raises:
// 	//     ValueError: if the sun does not rise or does not set
// 	//

// 	start, err := Sunrise(observer, date)
// 	if err != nil {
// 		return time.Time{}, time.Time{}, err
// 	}
// 	end, err := Sunset(observer, date)
// 	if err != nil {
// 		return time.Time{}, time.Time{}, err
// 	}

// 	if !daytime {
// 		start, err = Sunset(observer, date)
// 		if err != nil {
// 			return time.Time{}, time.Time{}, err
// 		}
// 		end, err = Sunrise(observer, date.Add(24*time.Hour))
// 		if err != nil {
// 			return time.Time{}, time.Time{}, err
// 		}
// 	}

// 	octant_duration := end.Sub(start).Seconds() / 8
// 	// octant_duration := datetime.timedelta((end - start).seconds / 8)

// 	// Mo,Sa,Fr,We,Th,Tu,Su
// 	octant_index := []int{1, 6, 4, 5, 3, 2, 7}

// 	weekday := date.Weekday()

// 	// convert to python weekday
// 	if weekday == time.Sunday {
// 		weekday = 6
// 	} else {
// 		weekday--
// 	}

// 	octant := octant_index[weekday]

// 	// TODO: is seconds correct?
// 	add := time.Duration(octant_duration*float64(octant)) * time.Second
// 	start = start.Add(add)
// 	end = start.Add(time.Duration(octant_duration))

// 	return start.In(date.Location()), end.In(date.Location()), nil
// }
