package utils

import (
	"time"
)

// 得到今天的前一天(日期) 比如今天是20160901 得到的日期为20160831
func GetYesDate() string {
	nTime := time.Now()
	yesTime := nTime.AddDate(0, 0, -1)
	return yesTime.Format("20060102")
}
