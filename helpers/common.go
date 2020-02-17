package helpers

import (
	"time"
)

func IsEmpty(data string) bool {
    if len(data) <= 0 {
        return true
    } else {
        return false
    }
}

func GetCurrentTime() string{
	dt := time.Now()
	return (dt.Format("15:04:05 02-01-2006"))
}