package main

import (
	"fmt"
	"lua_helper"
	"runtime/debug"
)

func main() {
	debug.SetPanicOnFault(true)
	L := lua_helper.GetState()
	err := L.DoString("print(remote_exec('3', '34', '34', '34',3,5,6))")
	if err != nil {
		debug.PrintStack()
	}

	fmt.Println("Lua output:", L.GetOutput())
}
