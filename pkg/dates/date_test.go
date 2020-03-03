package dates

import (
	"testing"
	"time"
)

func TestLatestBusinessEndOfDay(t *testing.T) {
	easternTime, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name string
		date time.Time
		want time.Time
	}{
		{
			name: "MondayBeforeClose",
			date: time.Date(2020, time.March, 2, 0, 0, 0, 0, easternTime),
			want: time.Date(2020, time.February, 28, 17, 0, 0, 0, easternTime),
		},
		{
			name: "MondayAfterClose",
			date: time.Date(2020, time.March, 2, 17, 0, 0, 0, easternTime),
			want: time.Date(2020, time.March, 2, 17, 0, 0, 0, easternTime),
		},
		{
			name: "TuesdayBeforeClose",
			date: time.Date(2020, time.March, 3, 0, 0, 0, 0, easternTime),
			want: time.Date(2020, time.March, 2, 17, 0, 0, 0, easternTime),
		},
		{
			name: "TuesdayBeforeClose",
			date: time.Date(2020, time.March, 3, 17, 0, 0, 0, easternTime),
			want: time.Date(2020, time.March, 3, 17, 0, 0, 0, easternTime),
		},
		{
			name: "WednesdayBeforeClose",
			date: time.Date(2020, time.March, 4, 0, 0, 0, 0, easternTime),
			want: time.Date(2020, time.March, 3, 17, 0, 0, 0, easternTime),
		},
		{
			name: "WednesdayAfterClose",
			date: time.Date(2020, time.March, 4, 17, 0, 0, 0, easternTime),
			want: time.Date(2020, time.March, 4, 17, 0, 0, 0, easternTime),
		},
		{
			name: "ThursdayBeforeClose",
			date: time.Date(2020, time.March, 5, 0, 0, 0, 0, easternTime),
			want: time.Date(2020, time.March, 4, 17, 0, 0, 0, easternTime),
		},
		{
			name: "ThursdayAfterClose",
			date: time.Date(2020, time.March, 5, 17, 0, 0, 0, easternTime),
			want: time.Date(2020, time.March, 5, 17, 0, 0, 0, easternTime),
		},
		{
			name: "FridayBeforeClose",
			date: time.Date(2020, time.March, 6, 0, 0, 0, 0, easternTime),
			want: time.Date(2020, time.March, 5, 17, 0, 0, 0, easternTime),
		},
		{
			name: "FridayAfterClose",
			date: time.Date(2020, time.March, 6, 17, 0, 0, 0, easternTime),
			want: time.Date(2020, time.March, 6, 17, 0, 0, 0, easternTime),
		},
		{
			name: "SaturdayBeforeClose",
			date: time.Date(2020, time.March, 7, 0, 0, 0, 0, easternTime),
			want: time.Date(2020, time.March, 6, 17, 0, 0, 0, easternTime),
		},
		{
			name: "SaturdayAfterClose",
			date: time.Date(2020, time.March, 7, 17, 0, 0, 0, easternTime),
			want: time.Date(2020, time.March, 6, 17, 0, 0, 0, easternTime),
		},
		{
			name: "SundayBeforeClose",
			date: time.Date(2020, time.March, 8, 0, 0, 0, 0, easternTime),
			want: time.Date(2020, time.March, 6, 17, 0, 0, 0, easternTime),
		},
		{
			name: "SundayAfterClose",
			date: time.Date(2020, time.March, 8, 17, 0, 0, 0, easternTime),
			want: time.Date(2020, time.March, 6, 17, 0, 0, 0, easternTime),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := LatestBusinessEndOfDay(tt.date)
			if actual != tt.want {
				t.Errorf("want %q; got %q", tt.want, actual)
			}
		})
	}
}
