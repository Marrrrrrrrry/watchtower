package scheduler

import (
	"testing"
	"time"
)

// These tests characterize the user-facing schedule format (the --schedule
// flag). They must keep passing as long as the format is unchanged; when the
// format is deliberately modernized, update this contract table together with
// the release notes.
func TestScheduleContract(t *testing.T) {
	ref := time.Date(2026, 1, 15, 10, 30, 45, 0, time.UTC) // a Thursday

	cases := []struct {
		name     string
		spec     string
		next     time.Time
		wantFail bool
	}{
		// 6-field specs (the documented format)
		{"daily at 4am", "0 0 4 * * *", time.Date(2026, 1, 16, 4, 0, 0, 0, time.UTC), false},
		{"every 5 minutes", "*/5 * * * * *", time.Date(2026, 1, 15, 10, 30, 50, 0, time.UTC), false},
		{"step hours", "0 30 */3 * * *", time.Date(2026, 1, 15, 12, 30, 0, 0, time.UTC), false},
		{"weekday filter", "0 0 4 * * 1-5", time.Date(2026, 1, 16, 4, 0, 0, 0, time.UTC), false}, // Friday
		{"month/day list", "0 30 4 1,15 * *", time.Date(2026, 2, 1, 4, 30, 0, 0, time.UTC), false},

		// 5-field specs keep the historical v1 mapping: fields are
		// second/minute/hour/day-of-month/month, day-of-week defaults to "*"
		{"5-field quirk: hourly at minute 4", "0 4 * * *", time.Date(2026, 1, 15, 11, 4, 0, 0, time.UTC), false},
		{"5-field quirk: every 10 seconds", "*/10 * * * *", time.Date(2026, 1, 15, 10, 30, 50, 0, time.UTC), false},

		// descriptors
		{"default interval string", "@every 86400s", time.Date(2026, 1, 16, 10, 30, 45, 0, time.UTC), false},
		{"every 90 seconds", "@every 90s", time.Date(2026, 1, 15, 10, 32, 15, 0, time.UTC), false},
		{"daily descriptor", "@daily", time.Date(2026, 1, 16, 0, 0, 0, 0, time.UTC), false},
		{"midnight descriptor", "@midnight", time.Date(2026, 1, 16, 0, 0, 0, 0, time.UTC), false},
		{"hourly descriptor", "@hourly", time.Date(2026, 1, 15, 11, 0, 0, 0, time.UTC), false},
		{"weekly descriptor", "@weekly", time.Date(2026, 1, 18, 0, 0, 0, 0, time.UTC), false}, // Sunday
		{"monthly descriptor", "@monthly", time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC), false},

		// additive v3 capability (v1 failed to start on these)
		{"timezone prefix", "TZ=Europe/Rome 0 0 4 * * *", time.Date(2026, 1, 16, 3, 0, 0, 0, time.UTC), false}, // 4am CET = 3am UTC

		// rejections
		{"too few fields", "0 4 * *", time.Time{}, true},
		{"too many fields", "0 0 4 * * * *", time.Time{}, true},
		{"garbage", "not a schedule", time.Time{}, true},
		{"empty", "", time.Time{}, true},
	}

	parser := V1CompatibleParser()
	for _, tc := range cases {
		schedule, err := parser.Parse(tc.spec)
		if tc.wantFail {
			if err == nil {
				t.Errorf("%s: spec %q unexpectedly parsed", tc.name, tc.spec)
			}
			continue
		}
		if err != nil {
			t.Errorf("%s: spec %q failed to parse: %v", tc.name, tc.spec, err)
			continue
		}
		if got := schedule.Next(ref); !got.Equal(tc.next) {
			t.Errorf("%s: spec %q next = %s, want %s", tc.name, tc.spec, got, tc.next)
		}
	}
}
