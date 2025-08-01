package core

import (
	"fmt"
	"time"
)

func TimeToStringTimestamp(t time.Time) (string, error) {
	if t.IsZero() {
		return "", fmt.Errorf("time is zero")
	}
	return fmt.Sprintf("%d", t.Unix()), nil
}
