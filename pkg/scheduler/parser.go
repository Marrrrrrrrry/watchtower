// Package scheduler holds the schedule-parsing contract for watchtower's
// --schedule flag and interval scheduling. Keeping the parser construction in
// one place documents the user-facing format and gives the characterization
// tests in this package a single anchor.
package scheduler

import (
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
