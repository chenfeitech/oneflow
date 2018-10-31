package lua_helper

import (
	"utils/conv"
	"github.com/stevedonovan/luar"
)

func init() {
	LuaPackages["conv"] = luar.Map{
		"String":  conv.String,
		"Int64":   conv.Int64,
		"Uint64":  conv.Uint64,
		"Int":     conv.Int,
		"Uint":    conv.Uint,
		"Float64": conv.Float64,
		"Time":    conv.Time,
	}
}
