package scheduler

import (
	"time"

	"dkl.ru/pact/orchestrator_service/iternal/clock"
)

type ScheduledWorker struct {
	Interval  time.Duration
	StartTime time.Duration
	Task      func()
	StopChan  chan struct{}
	Clock     clock.Clock
}

func NewScheduledWorker(interval, startTime time.Duration, task func(), clk clock.Clock) *ScheduledWorker {
	return &ScheduledWorker{
		Interval:  interval,
		StartTime: startTime,
		Task:      task,
		StopChan:  make(chan struct{}),
		Clock:     clk,
	}
}

func (w *ScheduledWorker) Start() {
	go func() {
		for {
			now := w.Clock.Now()
			startToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).
				Add(w.StartTime)

			if now.After(startToday) {
				startToday = startToday.Add(24 * time.Hour)
			}

			sleepDuration := time.Until(startToday)
			select {
			case <-w.Clock.After(sleepDuration):
				w.Task()
			case <-w.StopChan:
				return
			}

			ticker := w.Clock.NewTicker(w.Interval)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C():
					w.Task()
				case <-w.StopChan:
					return
				}
			}
		}
	}()
}

func (w *ScheduledWorker) Stop() {
	close(w.StopChan)
}
