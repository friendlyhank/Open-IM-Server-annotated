package utils

import "time"

// Convert timestamp to time.Time type 秒级时间戳转时间
func UnixSecondToTime(second int64) time.Time {
	return time.Unix(second, 0)
}

// Get the current timestamp by Mill 毫秒时间戳
func GetCurrentTimestampByMill() int64 {
	return time.Now().UnixNano() / 1e6
}

// 毫秒时间戳转时间
func UnixMillSecondToTime(millSecond int64) time.Time {
	return time.Unix(0, millSecond*1e6)
}

// string转time
func TimeStringToTime(timeString string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", timeString)
	return t, err
}

func TimeToString(t time.Time) string {
	return t.Format("2006-01-02")
}
