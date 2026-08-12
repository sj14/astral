package main

import (
	"bytes"
	"strings"
	"testing"
)

func runCLI(t *testing.T, args ...string) (stdout, stderr string, code int) {
	t.Helper()
	var outBuf, errBuf bytes.Buffer
	code = run(args, &outBuf, &errBuf)
	return outBuf.String(), errBuf.String(), code
}

func TestCLIOutput(t *testing.T) {
	stdout, stderr, code := runCLI(t, "-lat", "51.58", "-long", "6.52", "-time", "2021-04-30T21:12:11+02:00")
	if code != 0 {
		t.Fatalf("run returned exit code %d, stderr:\n%s", code, stderr)
	}

	for _, want := range []string{
		"Latitude\t51.58",
		"Longitude\t6.52",
		"Elevation\t0",
		"Daylight\t14h48m10s",
		"Night-Time\t9h9m55s",
		"Moon Phase\tFull Moon (17.611222222222224)",
	} {
		if !strings.Contains(stdout, want) {
			t.Errorf("output missing %q\ngot:\n%s", want, stdout)
		}
	}

	// Events must be printed in chronological order.
	wantOrder := []string{
		"Midnight",
		"Dawn (Astronomical)",
		"Dawn (Nautical)",
		"Dawn (Civil)",
		"Golden Hour Start",
		"Sunrise",
		"Golden Hour End",
		"Noon",
		"Golden Hour Start",
		"Sunset",
		"Golden Hour End",
		"Dusk (Civil)",
		"Dusk (Nautical)",
		"Dusk (Astronomical)",
	}
	pos := 0
	for _, want := range wantOrder {
		idx := strings.Index(stdout[pos:], want)
		if idx == -1 {
			t.Fatalf("expected %q to appear after position %d, not found in:\n%s", want, pos, stdout)
		}
		pos += idx + len(want)
	}
}

func TestCLIOutputUnreachableEventsOmitted(t *testing.T) {
	// At this latitude around midsummer the sun never gets 18 degrees below
	// the horizon, so astronomical dawn/dusk don't occur. Regression test for
	// a bug where the resulting zero-value time.Time was still rendered as a
	// bogus "Jan  1 00:00" entry instead of being left out.
	stdout, stderr, code := runCLI(t, "-lat", "51.58", "-long", "6.52", "-time", "2026-06-21T12:00:00+02:00")
	if code != 0 {
		t.Fatalf("run returned exit code %d, stderr:\n%s", code, stderr)
	}

	if strings.Contains(stdout, "Jan  1") {
		t.Errorf("output contains a stale zero-value time entry:\n%s", stdout)
	}
	for _, unwanted := range []string{"Dawn (Astronomical)", "Dusk (Astronomical)"} {
		if strings.Contains(stdout, unwanted) {
			t.Errorf("expected %q to be omitted when unreachable, got:\n%s", unwanted, stdout)
		}
	}
	if !strings.Contains(stderr, "sun never reaches 18 degrees below the horizon") {
		t.Errorf("expected stderr to report the unreachable event, got:\n%s", stderr)
	}
}

func TestCLIOutputPolarDayOmitsSunriseSunset(t *testing.T) {
	// Above the Arctic Circle in midsummer the sun never sets, so sunrise,
	// sunset and everything derived from them must be omitted rather than
	// shown with a bogus zero-value time or a garbage duration.
	stdout, stderr, code := runCLI(t, "-lat", "78", "-long", "15", "-time", "2026-06-21T12:00:00+02:00")
	if code != 0 {
		t.Fatalf("run returned exit code %d, stderr:\n%s", code, stderr)
	}

	if !strings.Contains(stdout, "Daylight\tn/a") {
		t.Errorf("expected Daylight to be n/a, got:\n%s", stdout)
	}
	if !strings.Contains(stdout, "Night-Time\tn/a") {
		t.Errorf("expected Night-Time to be n/a, got:\n%s", stdout)
	}
	for _, unwanted := range []string{"Sunrise", "Sunset", "Dawn", "Dusk", "Golden Hour"} {
		if strings.Contains(stdout, unwanted) {
			t.Errorf("expected %q to be omitted during polar day, got:\n%s", unwanted, stdout)
		}
	}
	for _, want := range []string{"Midnight", "Noon"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("expected %q to still be present, got:\n%s", want, stdout)
		}
	}
}

func TestCLIVersion(t *testing.T) {
	stdout, stderr, code := runCLI(t, "-version")
	if code != 0 {
		t.Fatalf("run returned exit code %d, stderr:\n%s", code, stderr)
	}
	for _, want := range []string{"version:", "commit:", "date:"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("version output missing %q:\n%s", want, stdout)
		}
	}
}

func TestCLIInvalidTime(t *testing.T) {
	_, stderr, code := runCLI(t, "-time", "not-a-time")
	if code == 0 {
		t.Fatal("expected run to return a non-zero exit code for an invalid -time value")
	}
	if !strings.Contains(stderr, "failed parsing time") {
		t.Errorf("expected stderr to mention the parse failure, got:\n%s", stderr)
	}
}

func TestCLIUnknownFlag(t *testing.T) {
	_, stderr, code := runCLI(t, "-nope")
	if code == 0 {
		t.Fatal("expected run to return a non-zero exit code for an unknown flag")
	}
	if strings.Count(stderr, "flag provided but not defined") != 1 {
		t.Errorf("expected the flag error to be written exactly once, got:\n%s", stderr)
	}
}
