package main

import (
	"flag"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/sj14/astral"
)

func main() {
	var (
		timeFlag      = flag.String("time", time.Now().Format(time.RFC3339), "day/time used for the calculation")
		formatFlag    = flag.String("format", "Jan _2 15:04:05", "time output format according to Go parsing")
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

	dates := make(map[time.Time]string)
	dates[t] = dashes
	dates[dawnAstronomical] = "Dawn (Astronomical)"
	dates[dawnNautical] = "Dawn (Nautical)"
	dates[dawnCivil] = "Dawn (Civil)         Twilight Start    Blue Hour Start"
	dates[goldenRisingStart] = "Golden Hour Start                      Blue Hour End"
	dates[sunrise] = "Sunrise              Twilight End"
	dates[goldenRisingEnd] = "Golden Hour End"
	dates[noon] = "Noon"
	dates[goldenSettingStart] = "Golden Hour Start"
	dates[sunset] = "Sunset               Twilight Start"
	dates[goldenSettingEnd] = "Golden Hour End                        Blue Hour Start"
	dates[duskCivil] = "Dusk (Civil)         Twilight End      Blue Hour End "
	dates[duskAstronomical] = "Dusk (Astronomical)"
	dates[duskNautical] = "Dusk (Nautical)"
	dates[midnight] = "Midnight"

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

	for _, key := range sortedTimes {
		if dates[key] == dashes {
			hhMMss := "15:04:05"
			preDashes := len(*formatFlag) - len(hhMMss) - 1
			if preDashes < 0 {
				preDashes = 0
			}
			fmt.Printf("%v %v\n", key.Format(strings.Repeat("┈", preDashes)+" "+hhMMss), dates[key])
			continue
		}
		fmt.Printf("%v %v\n", key.Format(*formatFlag), dates[key])
	}
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
