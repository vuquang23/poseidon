package time

import "time"

func RoundDown(t int64, d time.Duration) int64 {
	return time.Unix(t, 0).Truncate(time.Minute).Unix()
}
