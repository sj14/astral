package astral

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestJulianday(t *testing.T) {
	type args struct {
		date time.Time
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{name: "1", args: args{date: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)}, want: 2451544.5},
		{name: "2", args: args{date: time.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC)}, want: 2455927.5},
		{name: "3", args: args{date: time.Date(1987, 6, 19, 0, 0, 0, 0, time.UTC)}, want: 2446_965.5},
		{name: "4", args: args{date: time.Date(2013, 1, 1, 0, 0, 0, 0, time.UTC)}, want: 2456293.5},
		{name: "5", args: args{date: time.Date(1988, 6, 19, 0, 0, 0, 0, time.UTC)}, want: 2447_331.5},
		{name: "6", args: args{date: time.Date(2013, 6, 1, 0, 0, 0, 0, time.UTC)}, want: 2456444.5},
		{name: "7", args: args{date: time.Date(1867, 2, 1, 0, 0, 0, 0, time.UTC)}, want: 2402998.5},
		{name: "8", args: args{date: time.Date(3200, 11, 14, 0, 0, 0, 0, time.UTC)}, want: 2890153.5},
		{name: "9", args: args{date: time.Date(1957, 10, 4, 19, 26, 24, 0, time.UTC)}, want: 2.4361155e+06}, //2436116.31},
		{name: "10", args: args{date: time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC)}, want: 2.4515445e+06},   //2451545.0},
		{name: "11", args: args{date: time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC)}, want: 2451179.5},
		{name: "12", args: args{date: time.Date(1987, 1, 27, 0, 0, 0, 0, time.UTC)}, want: 2446_822.5},
		{name: "13", args: args{date: time.Date(1987, 6, 19, 12, 0, 0, 0, time.UTC)}, want: 2.4469655e+06}, // 2446_966.0},
		{name: "14", args: args{date: time.Date(1988, 1, 27, 0, 0, 0, 0, time.UTC)}, want: 2447_187.5},
		{name: "15", args: args{date: time.Date(1988, 6, 19, 12, 0, 0, 0, time.UTC)}, want: 2.4473315e+06}, //2447_332.0},
		{name: "16", args: args{date: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)}, want: 2415_020.5},
		{name: "17", args: args{date: time.Date(1600, 1, 1, 0, 0, 0, 0, time.UTC)}, want: 2305_447.5},
		{name: "18", args: args{date: time.Date(1600, 12, 31, 0, 0, 0, 0, time.UTC)}, want: 2305_812.5},
		{name: "19", args: args{date: time.Date(2012, 1, 1, 12, 0, 0, 0, time.UTC)}, want: 2.4559275e+06}, // 2455928.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := julianday(tt.args.date)
			require.Equal(t, tt.want, got, "wantStart: %v but got: %v", tt.want, got)
		})
	}
}

func TestJulianDayToCentury(t *testing.T) {
	type args struct {
		date float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{args: args{date: 2455927.5}, want: 0.119986311},
		{args: args{date: 2456293.5}, want: 0.130006845},
		{args: args{date: 2456444.5}, want: 0.134140999},
		{args: args{date: 2402998.5}, want: -1.329130732},
		{args: args{date: 2890153.5}, want: 12.00844627},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := jday_to_jcentury(tt.args.date)
			require.True(t, almostEqualf(tt.want, got, 0.000000001), "wantStart: %v but got: %v", tt.want, got)
		})
	}
}

func TestJulianCenturyToDay(t *testing.T) {
	type args struct {
		date float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{args: args{date: 0.119986311}, want: 2455927.5},
		{args: args{date: 0.130006845}, want: 2456293.5},
		{args: args{date: 0.134140999}, want: 2456444.5},
		{args: args{date: -1.32913073}, want: 2402998.5},
		{args: args{date: 12.00844627}, want: 2890153.5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := jcentury_to_jday(tt.args.date)
			// TODO: not sure if the accuracy is good enough
			require.True(t, almostEqualf(tt.want, got, 0.0001), "want: %v but got: %v", tt.want, got)
		})
	}
}
