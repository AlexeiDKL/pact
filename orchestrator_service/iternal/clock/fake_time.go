package clock

import "time"

type fakeClock struct {
	now   time.Time
	after chan time.Time
}

func NewFakeClock(start time.Time) *fakeClock {
	return &fakeClock{
		now:   start,
		after: make(chan time.Time, 1),
	}
}

func (f *fakeClock) Now() time.Time {
	return f.now
}

func (f *fakeClock) After(d time.Duration) <-chan time.Time {
	f.now = f.now.Add(d)
	f.after <- f.now
	return f.after
}

func (f *fakeClock) NewTicker(d time.Duration) Ticker {
	ch := make(chan time.Time, 1)
	return &fakeTicker{ch: ch, clock: f, step: d}
}

type fakeTicker struct {
	ch    chan time.Time
	clock *fakeClock
	step  time.Duration
}

func (t *fakeTicker) C() <-chan time.Time {
	t.clock.now = t.clock.now.Add(t.step)
	t.ch <- t.clock.now
	return t.ch
}

func (t *fakeTicker) Stop() {}
