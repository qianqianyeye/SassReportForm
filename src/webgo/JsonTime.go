package webgo

import (
	"time"
	"strconv"
	"strings"
)

type JsonDateTime time.Time
type JsonDate time.Time
type JsonTime time.Time

const (
	dateTimeFormart = "2006-01-02 15:04:05"
	dateFormart     = "2006-01-02"
	timeFormart     = "15:04:05"
)

func (p *JsonDateTime) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+dateTimeFormart+`"`, string(data), time.Local)
	*p = JsonDateTime(now)
	return
}
func (p *JsonDate) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+dateFormart+`"`, string(data), time.Local)
	*p = JsonDate(now)
	return
}
func (p *JsonTime) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+timeFormart+`"`, string(data), time.Local)
	*p = JsonTime(now)
	return
}
func (c JsonDateTime) MarshalJSON() ([]byte, error) {
	data := make([]byte, 0)
	data = append(data, '"')
	data = time.Time(c).AppendFormat(data, dateTimeFormart)
	data = append(data, '"')
	return data, nil
}
func (c JsonDate) MarshalJSON() ([]byte, error) {
	data := make([]byte, 0)
	data = append(data, '"')
	data = time.Time(c).AppendFormat(data, dateFormart)
	data = append(data, '"')
	return data, nil
}
func (c JsonTime) MarshalJSON() ([]byte, error) {
	data := make([]byte, 0)
	data = append(data, '"')
	data = time.Time(c).AppendFormat(data, timeFormart)
	data = append(data, '"')
	return data, nil
}
func (c JsonDateTime) String() string {
	return time.Time(c).Format(dateTimeFormart)
}

func (c JsonTime) String() string {
	return time.Time(c).Format(timeFormart)
}

//JSONDateTime 转为time
func (c JsonDateTime) CTime() time.Time {
	loc, _ := time.LoadLocation("Local")
	t:=time.Time(c).Format(dateTimeFormart)
	rt :=TimeZone(t)
	//fmt.Println(rt)
	result, _ := time.ParseInLocation("2006-01-02 15:04:05", rt, loc)
	return result
}

//比较是否是同一天
func GetDayIsEqual(startTime string, endTime string) bool {
	loc, _ := time.LoadLocation("Local")
	tstart, _ := time.ParseInLocation(dateTimeFormart, startTime, loc)
	tend, _ := time.ParseInLocation(dateTimeFormart, endTime, loc)
	ts := tstart.Format(dateFormart)
	te := tend.Format(dateFormart)
	if ts == te {
		return true
	} else {
		return false
	}
}

//将时间转换为字符串格式
func GetStringDateTime(times time.Time) string {
	ts := times.Format(dateTimeFormart)
	return ts
}

//获取日期差
func TimeSub(t1, t2 time.Time) int {
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, time.Local)
	t2 = time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, time.Local)
	return int(t1.Sub(t2).Hours() / 24)
}

//获取日期段 例如：2006-01-02 15:04:05 所属时间段为2006-01-02 15:00:00    2006-01-02 15:34:05所属时间段为2006-01-02 15:30:00
func TimeZone(startTime string) string {
	loc, _ := time.LoadLocation("Local") //重要：获取时区
	start, _ := time.ParseInLocation(dateTimeFormart, startTime, loc)
	startUnix := start.Unix()
	minute := start.Minute()
	sec := start.Second()
	if minute < 30 {
		subSec := minute*60 + sec
		startUnix = startUnix - int64(subSec)
	} else {
		subMinute := minute - 30
		subSec := subMinute*60 + sec
		startUnix = startUnix - int64(subSec)
	}
	parmTime := time.Unix(startUnix, 0).Format(dateTimeFormart)
	return parmTime
}

// StrToIntMonth 字符串月份转整数月份
func StrToIntMonth(month string) int {
	var data = map[string]int{
		"January":   0,
		"February":  1,
		"March":     2,
		"April":     3,
		"May":       4,
		"June":      5,
		"July":      6,
		"August":    7,
		"September": 8,
		"October":   9,
		"November":  10,
		"December":  11,
	}
	return data[month]
}

// GetTodayYMD 得到以sep为分隔符的年、月、日字符串(今天)
func GetTodayYMD(sep string) string {
	now := time.Now()
	year := now.Year()
	month := StrToIntMonth(now.Month().String())
	date := now.Day()

	var monthStr string
	var dateStr string
	if month < 9 {
		monthStr = "0" + strconv.Itoa(month+1)
	} else {
		monthStr = strconv.Itoa(month + 1)
	}

	if date < 10 {
		dateStr = "0" + strconv.Itoa(date)
	} else {
		dateStr = strconv.Itoa(date)
	}
	return strconv.Itoa(year) + sep + monthStr + sep + dateStr
}

// GetTodayYM 得到以sep为分隔符的年、月字符串(今天所属于的月份)
func GetTodayYM(sep string) string {
	now := time.Now()
	year := now.Year()
	month := StrToIntMonth(now.Month().String())

	var monthStr string
	if month < 9 {
		monthStr = "0" + strconv.Itoa(month+1)
	} else {
		monthStr = strconv.Itoa(month + 1)
	}
	return strconv.Itoa(year) + sep + monthStr
}

// GetYesterdayYMD 得到以sep为分隔符的年、月、日字符串(昨天)
func GetYesterdayYMD(sep string) string {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	todaySec := today.Unix()            //秒
	yesterdaySec := todaySec - 24*60*60 //秒
	yesterdayTime := time.Unix(yesterdaySec, 0)
	yesterdayYMD := yesterdayTime.Format("2006-01-02")
	return strings.Replace(yesterdayYMD, "-", sep, -1)
}

// GetTomorrowYMD 得到以sep为分隔符的年、月、日字符串(明天)
func GetTomorrowYMD(sep string) string {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	todaySec := today.Unix()           //秒
	tomorrowSec := todaySec + 24*60*60 //秒
	tomorrowTime := time.Unix(tomorrowSec, 0)
	tomorrowYMD := tomorrowTime.Format("2006-01-02")
	return strings.Replace(tomorrowYMD, "-", sep, -1)
}

// GetTodayTime 返回今天零点的time
func GetTodayTime() time.Time {
	now := time.Now()
	// now.Year(), now.Month(), now.Day() 是以本地时区为参照的年、月、日
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	return today
}

// GetYesterdayTime 返回昨天零点的time
func GetYesterdayTime() time.Time {
	now := time.Now()
	// now.Year(), now.Month(), now.Day() 是以本地时区为参照的年、月、日
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	yesterdaySec := today.Unix() - 24*60*60
	return time.Unix(yesterdaySec, 0)
}