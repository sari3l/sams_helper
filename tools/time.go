package tools

import (
	"fmt"
	"strconv"
	"time"
)

func UnixToTime(timestamp string) string {
	_time, _ := strconv.ParseInt(timestamp, 0, 64)
	tm := time.Unix(_time/1000, _time%1000)
	return fmt.Sprintf("%s", tm.Format("2006-01-02 03:04:05 PM"))
}
