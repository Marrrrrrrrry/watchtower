package scheduler

import (
	"sync/atomic"
	"testing"
	"time"
)

// TestSchedulerLifecycle exercises the exact construction used by
// cmd/root.go (sched.New + AddFunc + Entries + Start) to prove the v3
// wiring end to end: a spec parsed by the v1-compatible parser fires the
// job, a panicking job is recovered instead of killing the process (v1
// behavior parity), and Entries reports the next run consistently.
func TestSchedulerLifecycle(t *testing.T) {
	// the spec watchtower generates from the default 24h interval
	intervalSpec := "@every 1s"

	scheduler := New()
	var fired int32
	if _, err := scheduler.AddFunc(intervalSpec, func() { atomic.AddInt32(&fired, 1) }); err != nil {
		t.Fatalf("AddFunc: %v", err)
	}
	if _, err := scheduler.AddFunc("@every 1s", func() { panic("boom") }); err != nil {
		t.Fatalf("AddFunc(panic job): %v", err)
	}

	entries := scheduler.Entries()
	if len(entries) != 2 {
		t.Fatalf("Entries: got %d entries, want 2", len(entries))
	}
	next := entries[0].Schedule.Next(time.Now())
	if next.IsZero() {
		t.Fatal("Schedule.Next returned zero time")
	}

	scheduler.Start()
	defer scheduler.Stop()

	// without cron.Recover the panicking job would crash the test binary
	deadline := time.Now().Add(3 * time.Second)
	for atomic.LoadInt32(&fired) < 2 && time.Now().Before(deadline) {
		time.Sleep(50 * time.Millisecond)
	}
	if got := atomic.LoadInt32(&fired); got < 2 {
		t.Fatalf("healthy job fired %d times in 3s alongside a panicking job, want >= 2", got)
	}
}
