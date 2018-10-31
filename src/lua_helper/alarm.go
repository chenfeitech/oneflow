package lua_helper

import (
	"utils/alarm"
)

func init() {
	LuaGlobal["RaiseAlarm"] = alarm.RaiseAlarm
}
