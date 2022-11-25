package tools

import "time"

const FormatTimeString = "2006-01-02 15:04:05"

func GetNowDatetime() string {
	return time.Now().Format(FormatTimeString)
}
