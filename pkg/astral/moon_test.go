package astral

import (
	"testing"
	"time"
)

func TestMoon(t *testing.T) {
	type args struct {
		date time.Time
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{args: args{date: time.Date(2015, 12, 1, 0, 0, 0, 0, time.UTC)}, want: 19.477889},
		{args: args{date: time.Date(2015, 12, 2, 0, 0, 0, 0, time.UTC)}, want: 20.411222},
		{args: args{date: time.Date(2015, 12, 3, 0, 0, 0, 0, time.UTC)}, want: 21.266777},
		{args: args{date: time.Date(2014, 12, 1, 0, 0, 0, 0, time.UTC)}, want: 9.0556666},
		{args: args{date: time.Date(2014, 12, 2, 0, 0, 0, 0, time.UTC)}, want: 10.066777},
		{args: args{date: time.Date(2014, 1, 1, 0, 0, 0, 0, time.UTC)}, want: 0.033444},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MoonPhase(tt.args.date)
			almostEqualFloat(t, tt.want, got, 0.000001)
		})
	}
}
