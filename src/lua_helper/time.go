package lua_helper

import (
	"github.com/stevedonovan/luar"
	"time"
)

func init() {
	LuaPackages["time"] = luar.Map{
		"Now":             time.Now,
		"Date":            time.Date,
		"Parse":           time.Parse,
		"ParseInLocation": time.ParseInLocation,
		"Unix":            time.Unix,

		"ParseDuration": time.ParseDuration,
		"Since":         time.Since,
		"Duration":      Duration,

		"FixedZone":    time.FixedZone,
		"LoadLocation": time.LoadLocation,
		"Local":        time.Local,

		"Nanosecond":  (int64)(time.Nanosecond),
		"Microsecond": (int64)(time.Microsecond),
		"Millisecond": (int64)(time.Millisecond),
		"Second":      (int64)(time.Second),
		"Minute":      (int64)(time.Minute),
		"Hour":        (int64)(time.Hour),

		"January":   (int)(time.January),
		"February":  (int)(time.February),
		"March":     (int)(time.March),
		"April":     (int)(time.April),
		"May":       (int)(time.May),
		"June":      (int)(time.June),
		"July":      (int)(time.July),
		"August":    (int)(time.August),
		"September": (int)(time.September),
		"October":   (int)(time.October),
		"November":  (int)(time.November),
		"December":  (int)(time.December),

		"Sunday":    int(time.Sunday),
		"Monday":    int(time.Monday),
		"Tuesday":   int(time.Tuesday),
		"Wednesday": int(time.Wednesday),
		"Thursday":  int(time.Thursday),
		"Friday":    int(time.Friday),
		"Saturday":  int(time.Saturday),
	}
}

func Duration(t int64) time.Duration {
	return time.Duration(t)
}
