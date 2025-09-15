package worker

import (
	"log"
	"time"
)

type ScheduledWorker struct {
	IntervalMs  int64
	StartTimeMs int64
	Task        func()
	StopChan    chan struct{}
}

func NewScheduledWorker(intervalMs, startTimeMs int64, task func()) *ScheduledWorker {
	return &ScheduledWorker{
		IntervalMs:  intervalMs,
		StartTimeMs: startTimeMs,
		Task:        task,
		StopChan:    make(chan struct{}),
	}
}

func (w *ScheduledWorker) Start() {
	go func() {
		for {
			now := time.Now()
			startToday := time.Date(
				now.Year(), now.Month(), now.Day(),
				0, 0, 0, 0, now.Location(),
			).Add(time.Duration(w.StartTimeMs) * time.Millisecond)

			if now.After(startToday) {
				startToday = startToday.Add(24 * time.Hour)
			}

			sleepDuration := time.Until(startToday)
			log.Printf("ScheduledWorker: waiting %v until first run", sleepDuration)

			select {
			case <-time.After(sleepDuration):
				// Первый запуск
				w.Task()
			case <-w.StopChan:
				log.Println("ScheduledWorker: stopped before first run")
				return
			}

			// Повторные запуски
			ticker := time.NewTicker(time.Duration(w.IntervalMs) * time.Millisecond)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					w.Task()
				case <-w.StopChan:
					log.Println("ScheduledWorker: stopped")
					return
				}
			}
		}
	}()
}

func (w *ScheduledWorker) Stop() {
	close(w.StopChan)
}
