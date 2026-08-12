package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/logrusorgru/aurora/v4"
	"github.com/sj14/astral/pkg/astral"
)

var (
	// will be replaced during the build process
	version = "undefined"
	commit  = "undefined"
	date    = "undefined"
)

const (
	dashes = "┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run executes the CLI, writing normal output to stdout and diagnostics
// (non-fatal event errors, flag usage, fatal errors) to stderr, and
// returning the process exit code. It never calls os.Exit itself so it
// can be exercised directly from tests.
func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("astral", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var (
		timeFlag      = fs.String("time", time.Now().Format(time.RFC3339), "Day/time used for the calculation.")
		latFlag       = fs.Float64("lat", 0, "Latitude of the observer. Northern latitudes should be positive.")
		longFlag      = fs.Float64("long", 0, "Longitude of the observer. Eastern longitudes should be positive.")
		elevationFlag = fs.Float64("elev", 0, "Elevation and/or distance to nearest obscuring feature in metres above/below the location")
		versionFlag   = fs.Bool("version", false, fmt.Sprintf("Print version information of this release (%v).", version))
	)
	if err := fs.Parse(args); err != nil {
		// fs already wrote the error and usage to stderr.
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}

	if *versionFlag {
		fmt.Fprintf(stdout, "version: %v\n", version)
		fmt.Fprintf(stdout, "commit: %v\n", commit)
		fmt.Fprintf(stdout, "date: %v\n", date)
		return 0
	}

	observer := astral.Observer{Latitude: *latFlag, Longitude: *longFlag, Elevation: *elevationFlag}

	t, err := time.Parse(time.RFC3339, *timeFlag)
	if err != nil {
		fmt.Fprintf(stderr, "failed parsing time: %v\n", err)
		return 1
	}

	var events []event
	// addEvent records an event for display, unless err is non-nil (in which
	// case the event never occurred and its zero-value time must not be shown).
	addEvent := func(at time.Time, err error, color aurora.Value, desc string) {
		if err != nil {
			fmt.Fprintln(stderr, err)
			return
		}
		events = append(events, event{t: at, color: color, desc: desc})
	}

	addEvent(t, nil, aurora.Value{}, dashes)

	dawnCivil, err := astral.Dawn(observer, t, astral.DepressionCivil)
	addEvent(dawnCivil, err, aurora.BgIndex(111, " "), "Dawn (Civil)         Twilight Start    Blue Hour Start")

	dawnAstronomical, err := astral.Dawn(observer, t, astral.DepressionAstronomical)
	addEvent(dawnAstronomical, err, aurora.BgGray(8, " "), "Dawn (Astronomical)")

	dawnNautical, err := astral.Dawn(observer, t, astral.DepressionNautical)
	addEvent(dawnNautical, err, aurora.BgGray(15, " "), "Dawn (Nautical)")

	goldenRisingStart, goldenRisingEnd, err := astral.GoldenHour(observer, t, astral.SunDirectionRising)
	addEvent(goldenRisingStart, err, aurora.BgIndex(208, " "), "Golden Hour Start                      Blue Hour End")
	addEvent(goldenRisingEnd, err, aurora.BgIndex(226, " "), "Golden Hour End")

	sunrise, sunriseErr := astral.Sunrise(observer, t)
	addEvent(sunrise, sunriseErr, aurora.BgIndex(214, " "), "Sunrise              Twilight End")

	sunriseNextDay, sunriseNextDayErr := astral.Sunrise(observer, t.Add(24*time.Hour))
	if sunriseNextDayErr != nil {
		fmt.Fprintln(stderr, sunriseNextDayErr)
	}

	noon := astral.Noon(observer, t)
	addEvent(noon, nil, aurora.BgIndex(226, " "), "Noon")

	goldenSettingStart, goldenSettingEnd, err := astral.GoldenHour(observer, t, astral.SunDirectionSetting)
	addEvent(goldenSettingStart, err, aurora.BgIndex(214, " "), "Golden Hour Start")
	addEvent(goldenSettingEnd, err, aurora.BgIndex(111, " "), "Golden Hour End                        Blue Hour Start")

	sunset, sunsetErr := astral.Sunset(observer, t)
	addEvent(sunset, sunsetErr, aurora.BgIndex(208, " "), "Sunset               Twilight Start")

	duskCivil, err := astral.Dusk(observer, t, astral.DepressionCivil)
	addEvent(duskCivil, err, aurora.BgGray(18, " "), "Dusk (Civil)         Twilight End      Blue Hour End ")

	duskNautical, err := astral.Dusk(observer, t, astral.DepressionNautical)
	addEvent(duskNautical, err, aurora.BgGray(15, " "), "Dusk (Nautical)")

	duskAstronomical, err := astral.Dusk(observer, t, astral.DepressionAstronomical)
	addEvent(duskAstronomical, err, aurora.BgGray(8, " "), "Dusk (Astronomical)")

	midnight := astral.Midnight(observer, t)
	addEvent(midnight, nil, aurora.BgBlack(" "), "Midnight")

	moonPhase := astral.MoonPhase(t)
	moonDesc, err := astral.MoonPhaseDescription(moonPhase)
	if err != nil {
		fmt.Fprintf(stderr, "failed parsing moon phase: %v\n", err)
		return 1
	}

	sort.Slice(events, func(i, j int) bool { return events[i].t.Before(events[j].t) })

	fmt.Fprintf(stdout, "Date/Time\t%v\n", t.Format(time.UnixDate))
	fmt.Fprintf(stdout, "Latitude\t%v\nLongitude\t%v\nElevation\t%v\n", *latFlag, *longFlag, *elevationFlag)
	fmt.Fprintln(stdout)
	if sunriseErr != nil || sunsetErr != nil {
		fmt.Fprintln(stdout, "Daylight\tn/a")
	} else {
		fmt.Fprintf(stdout, "Daylight\t%v\n", sunset.Sub(sunrise).Truncate(1*time.Second))
	}
	if sunsetErr != nil || sunriseNextDayErr != nil {
		fmt.Fprintln(stdout, "Night-Time\tn/a")
	} else {
		fmt.Fprintf(stdout, "Night-Time\t%v\n", sunriseNextDay.Sub(sunset).Truncate(1*time.Second))
	}
	fmt.Fprintf(stdout, "Moon Phase\t%v (%v)\n", moonDesc, moonPhase)
	fmt.Fprintln(stdout)
	printEvents(stdout, events, t)

	return 0
}

const (
	dateTimeFormat = "Jan _2 15:04"
	timeFormat     = "15:04"
)

func printEvents(w io.Writer, events []event, t time.Time) {
	lastColor := aurora.BgBlack(" ")
	for _, ev := range events {
		// calculate when the particular phase happend or will happen
		inHours := math.Abs(ev.t.Sub(t).Truncate(1 * time.Hour).Hours())
		inMinutes := int(math.Abs((ev.t.Sub(t).Truncate(1 * time.Minute).Minutes()))) % 60

		agoOrUntil := fmt.Sprintf("%02.0f:%02d", inHours, inMinutes)
		if ev.t.Before(t) {
			agoOrUntil = fmt.Sprintf("-%s", agoOrUntil)
		} else {
			agoOrUntil = fmt.Sprintf("+%s", agoOrUntil)
		}

		// edge case for the given time
		if ev.desc == dashes {
			prefixDashesCount := len(dateTimeFormat) - len(timeFormat) - 1
			if prefixDashesCount < 0 {
				prefixDashesCount = 0
			}

			prefixDashes := ev.t.Format(strings.Repeat("┈", prefixDashesCount))
			midDashes := strings.Repeat("┈", len(agoOrUntil)+2)
			tStr := ev.t.Truncate(1 * time.Minute).Format(timeFormat)

			fmt.Fprintf(w, "%v %v %v %v %v\n", prefixDashes, tStr, midDashes, lastColor, ev.desc)
			continue
		}

		lastColor = ev.color
		fmt.Fprintf(w, "%v (%v) %v %v\n", ev.t.Format(dateTimeFormat), agoOrUntil, ev.color, ev.desc)
	}
}

type event struct {
	t     time.Time
	color aurora.Value
	desc  string
}
