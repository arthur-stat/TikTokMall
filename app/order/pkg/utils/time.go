package utils

import "time"

// FormatTime 格式化时间
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// ParseTime 解析时间字符串
func ParseTime(s string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", s)
}

// UnixToTime 将Unix时间戳转换为Time
func UnixToTime(unix int32) time.Time {
	return time.Unix(int64(unix), 0)
}

// TimeToUnix 将Time转换为Unix时间戳
func TimeToUnix(t time.Time) int32 {
	return int32(t.Unix())
}
