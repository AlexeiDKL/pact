package clock

import "time"

// Clock — абстракция над временем
type Clock interface {
	Now() time.Time
	After(d time.Duration) <-chan time.Time
	NewTicker(d time.Duration) Ticker
}

// Ticker — абстракция над time.Ticker
type Ticker interface {
	C() <-chan time.Time
	Stop()
}
