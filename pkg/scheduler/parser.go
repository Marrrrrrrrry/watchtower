// Package scheduler holds the schedule-parsing contract for watchtower's
// --schedule flag and interval scheduling. Keeping the parser construction in
// one place documents the user-facing format and gives the characterization
// tests in this package a single anchor.
package scheduler

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/robfig/cron/v3"
)

// logrusLogger adapts logrus to cron.Logger so scheduler panics (recovered by
// cron.Recover) surface through watchtower's logging pipeline and honor
// --log-level/--log-format, instead of a raw stdlib line on stdout.
type logrusLogger struct{}

func (logrusLogger) Info(msg string, keysAndValues ...interface{}) {
	log.Debugf("cron: %s %v", msg, keysAndValues)
}

func (logrusLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	log.WithError(err).Errorf("cron: %s %v", msg, keysAndValues)
}

// V1CompatibleParser returns a cron parser that replicates, field for field,
// the default parser of robfig/cron v1.2.0 (the parser watchtower has shipped
// with historically). It accepts 5 or 6 fields:
//
//	6 fields: second, minute, hour, day-of-month, month, day-of-week
//	5 fields: second, minute, hour, day-of-month, month (day-of-week
//	          defaults to "*"; NOT the standard 5-field interpretation)
//
// Descriptors (@every, @daily, ...) behave as in v1. TZ=/CRON_TZ= prefixes
// are additionally accepted (v3 handles them unconditionally in Parse); this
// is purely additive — v1 failed to start on such specs, so no existing
// deployment can rely on the old rejection.
func V1CompatibleParser() cron.Parser {
	return cron.NewParser(
		cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor,
	)
}

// New returns a cron scheduler configured for full parity with robfig/cron
// v1.2.0: the v1-compatible parser, plus a job wrapper that recovers from
// panics and logs them through logrus. The panic recovery matters — v1
// recovered panicking jobs by default, while v3 crashes the process unless
// cron.Recover is installed explicitly.
func New() *cron.Cron {
	return cron.New(
		cron.WithParser(V1CompatibleParser()),
		cron.WithChain(cron.Recover(logrusLogger{})),
	)
}

// ValidateSpec pre-flights a schedule specification before it is handed to
// AddFunc.
//
// Why this guard exists: robfig/cron v3.0.1 panics with a slice-bounds error
// (parser.go:99) on a spec that is a bare "TZ=..." or "CRON_TZ=..." prefix
// with no trailing expression — e.g. "TZ=Europe/Rome" — instead of returning
// a parse error like v1.2.0 did for the same input. The panic escapes AddFunc
// before any cron.Recover chain is active, so it crashes the process rather
// than failing with a usable message. Reported upstream (robfig/cron); this
// check converts the crash into a normal error and can be removed once a
// fixed cron release ships.
func ValidateSpec(spec string) error {
	trimmed := strings.TrimSpace(spec)
	for _, prefix := range []string{"CRON_TZ=", "TZ="} {
		// The panic shape is the prefix with no space after it: cron's TZ
		// branch slices the spec at the first space, which does not exist.
		if strings.HasPrefix(trimmed, prefix) && !strings.Contains(trimmed, " ") {
			return fmt.Errorf("invalid schedule %q: %s prefix without a schedule expression after it", spec, strings.TrimSuffix(prefix, "="))
		}
	}
	return nil
}
