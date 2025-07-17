package auxiliaryfunctions

import (
	"time"
)

func CreateFileVersion() int {
	return 1
	return int(time.Now().Unix())
}
