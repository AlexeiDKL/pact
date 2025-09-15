package scheduler

import (
	"sync"
	"testing"
	"time"

	"dkl.ru/pact/orchestrator_service/iternal/clock"
)

func TestScheduledWorker(t *testing.T) {
	start := time.Date(2025, 9, 9, 5, 0, 0, 0, time.UTC) // 05:00
	fc := clock.NewFakeClock(start)

	var wg sync.WaitGroup
	calls := 0
	howCall := 2
	wg.Add(howCall)

	var w *ScheduledWorker
	callback := func() {
		calls++
		wg.Done()
		if calls == howCall {
			w.Stop() // Останавливаем воркер после второго запуска
		}
	}
	w = NewScheduledWorker(
		24*time.Hour, // раз в сутки
		6*time.Hour,  // запуск в 06:00
		callback,
		fc,
	)

	w.Start()

	<-fc.After(time.Hour)      // Первый запуск
	<-fc.After(24 * time.Hour) // Второй запуск

	wg.Wait()

	if calls != 2 {
		t.Fatalf("ожидали 2 запуска, получили %d", calls)
	}
}
