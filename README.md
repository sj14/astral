# Astral

![Action](https://github.com/sj14/astral/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/sj14/astral)](https://goreportcard.com/report/github.com/sj14/astral)
[![GoDoc](https://godoc.org/github.com/sj14/astral?status.png)](https://godoc.org/github.com/sj14/astral)

Calculations for the position of the sun and moon.

This is a Go port of the Python [astral](https://github.com/sffjunkie/astral) package.

The `astral` package provides the means to calculate the following times of the sun:

* dawn
* sunrise
* noon
* midnight
* sunset
* dusk
* daylight
* night
* twilight
* blue hour
* golden hour
* ~~rahukaalam~~ TODO

plus solar azimuth and elevation at a specific latitude/longitude.
It can also calculate the moon phase for a specific date.

## CLI

Besides the package for usage in you own programs, we also provide a tool for showing the data.

### Installation

```text
go get github.com/sj14/astral/cmd/astral
```

### Usage

```text
Usage of astral:
  -elev float
        elevation of the observer
  -format string
        time output format according to Go parsing (default "Jan _2 15:04:05")
  -lat float
        latitude of the observer
  -long float
        longitude of the observer
  -time string
        day/time used for the calculation (defaults to current time)
```

### Example

```text
$ astral -lat 51.58 -long 6.52
Sun Apr 18 11:12:45 CEST 2021
Latitude 51.58 Longitude 6.52 Elevation 0
Apr 18 05:55:07 Dawn
Apr 18 05:55:07 Blue Hour Start
Apr 18 06:09:04 Blue Hour End
Apr 18 06:09:04 Golden Hour Start
Apr 18 06:31:44 Sunrise
Apr 18 07:15:56 Golden Hour End
Apr 18 13:33:18 Noon
Apr 18 19:51:24 Golden Hour Start
Apr 18 20:35:45 Sunset
Apr 18 20:58:33 Golden Hour End
Apr 18 20:58:33 Blue Hour Start
Apr 18 21:12:35 Blue Hour End
Apr 18 21:12:35 Dusk
Apr 19 01:33:13 Midnight
Moon Phase: New Moon (5.47788888888889)
```
