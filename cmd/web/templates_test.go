package main

import (
	"testing"
	"time"

	"snippetbox.saran.net/internal/assert"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC single digit",
			tm:   time.Date(2022, 3, 17, 5, 30, 0, 0, time.UTC),
			want: "17 Mar 2022 at 05:30",
		},
		{
			name: "UTC double digit",
			tm:   time.Date(2022, 3, 17, 10, 30, 0, 0, time.UTC),
			want: "17 Mar 2022 at 10:30",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2022, 3, 17, 10, 30, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 Mar 2022 at 09:30",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)
			assert.Equal(t, hd, tt.want)
		})
	}
}
