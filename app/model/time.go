package model

import (
	"fmt"
	"time"
)

type ModelTime time.Time

const (
	timeFormart = "2006-01-02 15:04:05"
	zone        = "Asia/Shanghai"
)

func (t *ModelTime) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+timeFormart+`"`, string(data), time.Local)
	*t = ModelTime(now)
	return
}

func (t ModelTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(t).Format(timeFormart))
	return []byte(stamp), nil
}

func (t ModelTime) Time() time.Time {
	return time.Time(t)
}

func (t ModelTime) Format() string {
	return time.Time(t).Format(timeFormart)
}

func (t ModelTime) String() string {
	return time.Time(t).Format(timeFormart)
}
