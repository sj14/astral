package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/logrusorgru/aurora/v3"
	"github.com/sj14/astral"
)

func main() {
	var (
		dateTimeFormat = "Jan _2 15:04"
		timeFormat     = "15:04"

		timeFlag      = flag.String("time", time.Now().Format(time.RFC3339), "day/time used for the calculation")
		latFlag       = flag.Float64("lat", 0, "latitude of the observer")
		longFlag      = flag.Float64("long", 0, "longitude of the observer")
		elevationFlag = flag.Float64("elev", 0, "elevation of the observer")
	)
	flag.Parse()

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

	dashes := "┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈"

	dates := make(map[time.Time]colorDesc)
	dates[t] = colorDesc{desc: dashes}
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

	fmt.Printf("Date/Time\t%v\n", t.Format(time.UnixDate))
	fmt.Printf("Latitude\t%v\nLongitude\t%v\nElevation\t%v\n", *latFlag, *longFlag, *elevationFlag)
	fmt.Println()
	fmt.Printf("Daylight\t%v\n", sunset.Sub(sunrise).Truncate(1*time.Second))
	fmt.Printf("Night-Time\t%v\n", sunriseNextDay.Sub(sunset).Truncate(1*time.Second))
	fmt.Printf("Moon Phase\t%v (%v)\n", moonDesc, moonPhase)
	fmt.Println()

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
		if dates[key].desc == dashes {
			prefixDashesCount := len(dateTimeFormat) - len(timeFormat) - 1
			if prefixDashesCount < 0 {
				prefixDashesCount = 0
			}

			prefixDashes := key.Format(strings.Repeat("┈", prefixDashesCount))
			midDashes := strings.Repeat("┈", len(agoOrUntil)+2)
			t := key.Truncate(1 * time.Minute).Format(timeFormat)

			fmt.Printf("%v %v %v %v %v\n", prefixDashes, t, midDashes, lastColor, dates[key].desc)
			continue
		}

		lastColor = dates[key].color
		fmt.Printf("%v (%v) %v %v\n", key.Format(dateTimeFormat), agoOrUntil, dates[key].color, dates[key].desc)
	}
}

type colorDesc struct {
	color aurora.Value
	desc  string
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
