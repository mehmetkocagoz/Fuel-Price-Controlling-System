package datafunctions

import (
	"fmt"
	"strings"
	"time"
)

func switchMonthToNumber(month string) string {
	switch month {
	case "January":
		return "01"
	case "February":
		return "02"
	case "March":
		return "03"
	case "April":
		return "04"
	case "May":
		return "05"
	case "June":
		return "06"
	case "July":
		return "07"
	case "August":
		return "08"
	case "September":
		return "09"
	case "October":
		return "10"
	case "November":
		return "11"
	case "December":
		return "12"
	}
	return "0"
}

func ConvertTimestamp(date string) int64 {
	//fmt.Println("converting..", date)
	// I know that our date will come like int string int format.
	// So first I'm going to convert it to int-int-int format.
	parsedDate := strings.Split(date, " ")
	month := switchMonthToNumber(parsedDate[1])
	date = parsedDate[2] + "-" + month + "-" + parsedDate[0]
	layout := "2006-01-02"
	t, err := time.Parse(layout, date)

	if err != nil {
		fmt.Println("time.Parse has failed: ", err)
	}
	// 10800000 is 3 hours in milliseconds.
	// I'm going to subtract it from the timestamp.
	// Because the data I get from tppd.com.tr is 3 hours ahead of the data I get from bloomberght.com
	return (t.Unix()*1000 - 10800000)
}
