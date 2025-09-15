package clock

import "time"

type realClock struct{}

func NewRealClock() Clock {
	return &realClock{}
}

func (realClock) Now() time.Time {
	return time.Now()
}

func (realClock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

func (realClock) NewTicker(d time.Duration) Ticker {
	return &realTicker{time.NewTicker(d)}
}

type realTicker struct {
	*time.Ticker
}

func (t *realTicker) C() <-chan time.Time {
	return t.Ticker.C
}
