package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/logrusorgru/aurora/v3"
	"github.com/sj14/astral/pkg/astral"
)

var (
	// will be replaced during the build process
	version = "undefined"
	commit  = "undefined"
	date    = "undefined"
)

const dashes = "┈"

func main() {
	var (
		dateTimeFormat = TimeFormatAstral
		timeFormat     = "15:04"

		timeFlag      = flag.String("time", time.Now().Format(time.RFC3339), "day/time used for the calculation")
		latFlag       = flag.Float64("lat", 0, "latitude of the observer")
		longFlag      = flag.Float64("long", 0, "longitude of the observer")
		elevationFlag = flag.Float64("elev", 0, "elevation of the observer")
		versionFlag   = flag.Bool("version", false, fmt.Sprintf("print version information of this release (%v)", version))
	)
	flag.StringVar(&dateTimeFormat, "dtfmt", dateTimeFormat, "date/time format to use for output")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("version: %v\n", version)
		fmt.Printf("commit: %v\n", commit)
		fmt.Printf("date: %v\n", date)
		os.Exit(0)
	}

	dateTimeFormat, err := FormatName(dateTimeFormat)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error '-dtfmt': %s", err)
		return
	}

	observer := astral.Observer{Latitude: *latFlag, Longitude: *longFlag, Elevation: *elevationFlag}

	t, err := time.Parse(time.RFC3339, *timeFlag)
	if err != nil {
		log.Fatalf("failed parsing time: %v\n", err)
	}

	dawnCivil, err := astral.Dawn(observer, t, astral.DepressionCivil)
	if err != nil {
		log.Println(err)
	}
	dawnAstronomical, err := astral.Dawn(observer, t, astral.DepressionAstronomical)
	if err != nil {
		log.Println(err)
	}
	dawnNautical, err := astral.Dawn(observer, t, astral.DepressionNautical)
	if err != nil {
		log.Println(err)
	}

	goldenRisingStart, goldenRisingEnd, err := astral.GoldenHour(observer, t, astral.SunDirectionRising)
	if err != nil {
		log.Println(err)
	}

	sunrise, err := astral.Sunrise(observer, t)
	if err != nil {
		log.Println(err)
	}

	sunriseNextDay, err := astral.Sunrise(observer, t.Add(24*time.Hour))
	if err != nil {
		log.Println(err)
	}

	noon := astral.Noon(observer, t)

	goldenSettingStart, goldenSettingEnd, err := astral.GoldenHour(observer, t, astral.SunDirectionSetting)
	if err != nil {
		log.Println(err)
	}

	sunset, err := astral.Sunset(observer, t)
	if err != nil {
		log.Println(err)
	}

	duskCivil, err := astral.Dusk(observer, t, astral.DepressionCivil)
	if err != nil {
		log.Println(err)
	}
	duskAstronomical, err := astral.Dusk(observer, t, astral.DepressionAstronomical)
	if err != nil {
		log.Println(err)
	}
	duskNautical, err := astral.Dusk(observer, t, astral.DepressionNautical)
	if err != nil {
		log.Println(err)
	}

	midnight := astral.Midnight(observer, t)

	moonPhase := astral.MoonPhase(t)
	moonDesc, err := astral.MoonPhaseDescription(moonPhase)
	if err != nil {
		log.Fatalf("failed parsing moon phase: %v", err)
	}

	dates := make(map[time.Time]colorDesc)
	dates[t] = colorDesc{special: true}
	dates[dawnAstronomical] = colorDesc{color: aurora.BgGray(8, " "), desc: "Dawn (Astronomical)"}
	dates[dawnNautical] = colorDesc{color: aurora.BgGray(15, " "), desc: "Dawn (Nautical)"}
	dates[dawnCivil] = colorDesc{color: aurora.BgIndex(111, " "), desc: "Dawn (Civil)         Twilight Start    Blue Hour Start"}
	dates[goldenRisingStart] = colorDesc{color: aurora.BgIndex(208, " "), desc: "Golden Hour Start                      Blue Hour End"}
	dates[sunrise] = colorDesc{color: aurora.BgIndex(214, " "), desc: "Sunrise              Twilight End"}
	dates[goldenRisingEnd] = colorDesc{color: aurora.BgIndex(226, " "), desc: "Golden Hour End"}
	dates[noon] = colorDesc{color: aurora.BgIndex(226, " "), desc: "Noon"}
	dates[goldenSettingStart] = colorDesc{color: aurora.BgIndex(214, " "), desc: "Golden Hour Start"}
	dates[sunset] = colorDesc{color: aurora.BgIndex(208, " "), desc: "Sunset               Twilight Start"}
	dates[goldenSettingEnd] = colorDesc{color: aurora.BgIndex(111, " "), desc: "Golden Hour End                        Blue Hour Start"}
	dates[duskCivil] = colorDesc{color: aurora.BgGray(18, " "), desc: "Dusk (Civil)         Twilight End      Blue Hour End "}
	dates[duskNautical] = colorDesc{color: aurora.BgGray(15, " "), desc: "Dusk (Nautical)"}
	dates[duskAstronomical] = colorDesc{color: aurora.BgGray(8, " "), desc: "Dusk (Astronomical)"}
	dates[midnight] = colorDesc{color: aurora.BgBlack(" "), desc: "Midnight"}

	var sortedTimes timeSlice

	for key := range dates {
		sortedTimes = append(sortedTimes, key)
	}
	sort.Sort(sortedTimes)

	fmt.Printf("Date/Time\t%v\n", t.Format(dateTimeFormat))
	fmt.Printf("Latitude\t%v\nLongitude\t%v\nElevation\t%v\n", *latFlag, *longFlag, *elevationFlag)
	fmt.Println()
	fmt.Printf("Daylight\t%v\n", sunset.Sub(sunrise).Truncate(1*time.Second))
	fmt.Printf("Night-Time\t%v\n", sunriseNextDay.Sub(sunset).Truncate(1*time.Second))
	fmt.Printf("Moon Phase\t%v (%.2f)\n", moonDesc, moonPhase)
	fmt.Println()

	longestDesc := 0
	for _, key := range sortedTimes {
		ld := len(dates[key].desc)
		if ld > longestDesc {
			longestDesc = ld
		}
	}

	lastColor := aurora.BgBlack(" ")
	for _, key := range sortedTimes {

		// calculate when the particular phase happend or will happen
		inHours := math.Abs(key.Sub(t).Truncate(1 * time.Hour).Hours())
		inMinutes := int(math.Abs((key.Sub(t).Truncate(1 * time.Minute).Minutes()))) % 60

		agoOrUntil := fmt.Sprintf("%02.0f:%02d", inHours, inMinutes)
		if key.Before(t) {
			agoOrUntil = fmt.Sprintf("-%s", agoOrUntil)
		} else {
			agoOrUntil = fmt.Sprintf("+%s", agoOrUntil)
		}

		// edge case for the given time
		if dates[key].special {
			prefixDashesCount := len(dateTimeFormat) - len(timeFormat) - 1
			if prefixDashesCount < 0 {
				prefixDashesCount = 0
			}

			prefixDashes := key.Format(strings.Repeat(dashes, prefixDashesCount))
			midDashes := strings.Repeat(dashes, len(agoOrUntil)+2)
			endDashes := strings.Repeat(dashes, longestDesc)
			t := key.Truncate(1 * time.Minute).Format(timeFormat)

			fmt.Printf("%v %v %v %v %v\n", prefixDashes, t, midDashes, lastColor, endDashes)
			continue
		}

		lastColor = dates[key].color
		fmt.Printf("%v (%v) %v %v\n", key.Format(dateTimeFormat), agoOrUntil, dates[key].color, dates[key].desc)
	}
}

type colorDesc struct {
	color   aurora.Value
	desc    string
	special bool
}

type timeSlice []time.Time

func (p timeSlice) Len() int {
	return len(p)
}

func (p timeSlice) Less(i, j int) bool {
	return p[i].Before(p[j])
}

func (p timeSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

const (
	// astrals default time format
	TimeFormatAstral = "Jan _2 15:04"
	// TimeFormatGo handles Go's default time.Now() format
	// (e.g. 2019-01-26 09:43:57.377055 +0100 CET m=+0.644739467)
	TimeFormatGo = "2006-01-02 15:04:05.999999999 -0700 MST"
	// TimeFormatSimple handles "2019-01-25 21:51:38"
	TimeFormatSimple = "2006-01-02 15:04:05.999999999"
	// TimeFormatHTTP instead of importing main with http.TimeFormat
	// which would increase the binary size significantly.
	TimeFormatHTTP = "Mon, 02 Jan 2006 15:04:05 GMT"
)

func FormatName(format string) (string, error) {

	if format == TimeFormatAstral {
		return TimeFormatAstral, nil
	}

	switch strings.ToLower(format) {
	case "":
		return TimeFormatGo, nil
	case "unix":
		return time.UnixDate, nil
	case "ruby":
		return time.RubyDate, nil
	case "ansic":
		return time.ANSIC, nil
	case "rfc822":
		return time.RFC822, nil
	case "rfc822z":
		return time.RFC822Z, nil
	case "rfc850":
		return time.RFC850, nil
	case "rfc1123":
		return time.RFC1123, nil
	case "rfc1123z":
		return time.RFC1123Z, nil
	case "rfc3339":
		return time.RFC3339, nil
	case "rfc3339nano":
		return time.RFC3339Nano, nil
	case "stamp":
		return time.Stamp, nil
	case "stampmilli":
		return time.StampMilli, nil
	case "stampmicro":
		return time.StampMicro, nil
	case "stampnano":
		return time.StampNano, nil
	case "http":
		return TimeFormatHTTP, nil
	default:
		return TimeFormatGo, fmt.Errorf("failed to parse format %q", format)
	}
}
