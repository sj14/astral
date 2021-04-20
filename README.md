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
Date/Time       Wed Apr 21 17:30:26 CEST 2021
Latitude        51.58
Longitude       6.52
Elevation       0

Daylight        14h15m19s
Night-Time      9h42m37s
Moon Phase      First Quarter (8.122333333333334)

Apr 21 04:09:44 Dawn (Astronomical)
Apr 21 05:02:44 Dawn (Nautical)
Apr 21 05:48:21 Dawn (Civil)         Twilight Start    Blue Hour Start
Apr 21 06:02:31 Golden Hour Start                      Blue Hour End
Apr 21 06:25:28 Sunrise              Twilight End
Apr 21 07:10:01 Golden Hour End
Apr 21 13:32:40 Noon
┈┈┈┈┈┈ 17:30:26 ┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈┈
Apr 21 19:56:03 Golden Hour Start
Apr 21 20:40:47 Sunset               Twilight Start
Apr 21 21:03:52 Golden Hour End                        Blue Hour Start
Apr 21 21:18:07 Dusk (Civil)         Twilight End      Blue Hour End 
Apr 21 22:04:06 Dusk (Nautical)
Apr 21 22:57:48 Dusk (Astronomical)
Apr 22 01:32:35 Midnight
```
