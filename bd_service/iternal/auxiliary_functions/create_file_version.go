package auxiliaryfunctions

import (
	"time"
)

func CreateFileVersion() int {
	return int(time.Now().Unix())
}
