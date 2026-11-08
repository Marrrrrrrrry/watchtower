package scheduler

import "testing"

// ValidateSpec guards against the robfig/cron v3.0.1 panic on bare TZ=/
// CRON_TZ= prefixes (see parser.go for the full rationale). These cases lock
// the guard's behavior so the upstream workaround cannot silently regress.
func TestValidateSpec(t *testing.T) {
	cases := []struct {
		name    string
		spec    string
		wantErr bool
	}{
		// the upstream panic shapes
		{"bare TZ prefix with location", "TZ=Europe/Rome", true},
		{"bare CRON_TZ prefix with location", "CRON_TZ=Asia/Hong_Kong", true},
		{"bare TZ= literal", "TZ=", true},
		{"prefix and whitespace only", "  TZ=  ", true},

		// valid shapes must keep passing
		{"TZ prefix with expression", "TZ=Europe/Rome 0 0 4 * * *", false},
		{"CRON_TZ prefix with expression", "CRON_TZ=Asia/Hong_Kong 0 0 6 * * *", false},
		{"plain 5-field", "0 0 6 * *", false},
		{"plain 6-field", "0 0 6 * * *", false},
		{"descriptor", "@daily", false},
		{"5-field quirk", "0 4 * * *", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateSpec(tc.spec)
			if tc.wantErr && err == nil {
				t.Fatalf("expected an error for spec %q, got nil", tc.spec)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error for spec %q: %v", tc.spec, err)
			}
		})
	}
}
