package helper

import (
	"fmt"
	"time"
)

func GenerateTimeUUID() string {
	now := time.Now()
	timeComponent := fmt.Sprintf("%04d%02d%02d%02d%02d%02d%09d",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(),
		now.Nanosecond())
	if len(timeComponent) > 32 {
		timeComponent = timeComponent[:32]
	}
	for len(timeComponent) < 32 {
		timeComponent += "0"
	}
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		timeComponent[0:8],
		timeComponent[8:12],
		timeComponent[12:16],
		timeComponent[16:20],
		timeComponent[20:32])
}
