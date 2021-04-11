package main

import (
	"flag"
	"fmt"
	"log"
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

	fmt.Printf("%v\n", t.Format(time.UnixDate))
	fmt.Printf("Latitude %v Longitude %v Elevation %v\n", *latFlag, *longFlag, *elevationFlag)

	// sun := astral.Sun(observer, now, astral.DepressionCivil)
	// fmt.Printf("sun at %+v\n", sun)

	dawn, err := astral.Dawn(observer, t, astral.DepressionCivil)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("%v Dawn\n", dawn.Format(*formatFlag))

	blueStart, blueEnd, err := astral.BlueHour(observer, t, astral.SunDirectionRising)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("%v Blue Hour Start\n", blueStart.Format(*formatFlag))
	fmt.Printf("%v Blue Hour End\n", blueEnd.Format(*formatFlag))

	goldenRisingStart, goldenRisingEnd, err := astral.GoldenHour(observer, t, astral.SunDirectionRising)
	if err != nil {
		log.Println(err)
	}

	sunrise, err := astral.Sunrise(observer, t)
	if err != nil {
		log.Println(err)
	}

	fmt.Printf("%v Golden Hour Start\n", goldenRisingStart.Format(*formatFlag))
	fmt.Printf("%v Sunrise\n", sunrise.Format(*formatFlag))
	fmt.Printf("%v Golden Hour End\n", goldenRisingEnd.Format(*formatFlag))

	noon := astral.Noon(observer, t)
	fmt.Printf("%v Noon\n", noon.Format(*formatFlag))

	goldenSettingStart, goldenSettingEnd, err := astral.GoldenHour(observer, t, astral.SunDirectionSetting)
	if err != nil {
		log.Println(err)
	}

	sunset, err := astral.Sunset(observer, t)
	if err != nil {
		log.Println(err)
	}

	fmt.Printf("%v Golden Hour Start\n", goldenSettingStart.Format(*formatFlag))
	fmt.Printf("%v Sunset\n", sunset.Format(*formatFlag))
	fmt.Printf("%v Golden Hour End\n", goldenSettingEnd.Format(*formatFlag))

	blueSettingStart, blueSettingEnd, err := astral.BlueHour(observer, t, astral.SunDirectionSetting)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("%v Blue Hour Start\n", blueSettingStart.Format(*formatFlag))
	fmt.Printf("%v Blue Hour End\n", blueSettingEnd.Format(*formatFlag))

	dusk, err := astral.Dusk(observer, t, astral.DepressionCivil)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("%v Dusk\n", dusk.Format(*formatFlag))

	// daylightStart, daylightEnd, err := astral.Daylight(observer, now)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Printf("daylight start: %v daylight end: %v\n", daylightStart, daylightEnd)

	// nightStart, nightEnd, err := astral.Night(observer, now)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Printf("night start: %v night end: %v\n", nightStart, nightEnd)

	// rahukaalamStart, rahukaalamEnd := astral.Rahukaalam(observer, now, true)
	// fmt.Printf("rahukaalam start: %v rahukaalam end: %v\n", rahukaalamStart, rahukaalamEnd)

	// elevation := astral.Elevation(observer, t, true)
	// elevationTime, err := astral.TimeAtElevation(observer, elevation, t, astral.SunDirectionSetting)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Printf("%v Elevation\n", elevationTime.Format(*formatFlag))

	midnight := astral.Midnight(observer, t)
	fmt.Printf("%v Midnight\n", midnight.Format(*formatFlag))

	// twilightStart, twilightEnd, err := astral.Twilight(observer, now, astral.SunDirectionSetting) // same as sunset and dusk
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Printf("twilight start: %v twilight end: %v\n", twilightStart, twilightEnd)

	// zenith, azimuth := astral.ZenithAndAzimuth(observer, now, true)
	// fmt.Printf("zenith: %v azimuth: %v\n", zenith, azimuth)

	moonPhase := astral.MoonPhase(t)
	moonDesc, err := astral.MoonPhaseDescription(moonPhase)
	if err != nil {
		log.Fatalf("failed parsing moon phase: %v", err)
	}
	fmt.Printf("Moon Phase: %v (%v)\n", moonDesc, moonPhase)
}
