package measurment

import (
	"fmt"
	"time"
)

func Trace(s string) (string, time.Time) {
	fmt.Println("START:", s)
	return s, time.Now()
}

func Un(s string, startTime time.Time) {
	endTime := time.Now()
	fmt.Println("  END:", s, "ElapsedTime in seconds:", endTime.Sub(startTime))
}
