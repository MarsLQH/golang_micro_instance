package common

import (
	"time"
)

const TimeLayoutUtc8 = "2006-01-02 15:04:05 +0800"

func TimeTripTimezone(timeWithTimezone string) (string, error) {
	t, err := time.ParseInLocation(TimeLayoutUtc8, timeWithTimezone, time.Local)
	if err != nil {
		return "", err
	}
	return time.Unix(t.Unix(), 0).Format("2006-01-02 15:04:05"), nil
}
