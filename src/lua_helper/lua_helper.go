package lua_helper

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"model"

	"github.com/yuin/gopher-lua"
	log "github.com/cihub/seelog"
)

// var (
// 	LuaGlobal   = luar.Map{}
// 	LuaPackages = make(map[string]luar.Map)
// )

type iState struct {
	*lua.LState
	writer            io.Writer
	buffer            bytes.Buffer
	remote_task_count int
	Products          *model.Products
	FlowInstance      *model.FlowInst
	Task              *model.Task
	TaskInstance      *model.TaskInst
	RemoteExecRecords []*RemoteExecRec
	RemoteExecEnv     map[string]string
	RemoteExecUseRoot bool
}

func (l *iState) iRegister() {
	log.Info("Do Register ")
	fmt.Println("Do Register ")
	l.Remote_init(l.LState)
//	funcs := luar.Map{}
//	for name, fun := range LuaGlobal {
//		funcs[name] = fun
//	}
//
//	lval := reflect.ValueOf(l)
//	ltype := reflect.TypeOf(l)
//	for i := 0; i < ltype.NumMethod(); i++ {
//		method := ltype.Method(i)
//		if strings.HasPrefix(method.Name, "Lua_") {
//			fun_name := strings.TrimPrefix(method.Name, "Lua_")
//			log.Info("Register ", fun_name)
//			funcs[fun_name] = lval.MethodByName(method.Name).Interface()
//		}
//	}
//	luar.Register(l.State, "", funcs)
}

func (l *iState) GetOutput() string {
	if buf_writer, ok := l.writer.(*bufio.Writer); ok {
		buf_writer.Flush()
	}
	return l.buffer.String()
}

// TODO
func GetStateByPId(pid string) *iState {
	return GetState()
}

func GetState() *iState {
	L := &iState{}
	// L.State = luar.Init()
	L.LState = lua.NewState()
	L.RemoteExecRecords = make([]*RemoteExecRec, 0)
	L.RemoteExecUseRoot = true

	L.writer = bufio.NewWriter(&L.buffer)

	L.OpenLibs()
	L.iRegister()
	L.Register("gassert", assert)
//	L.Register("print", L.print)
//	luar.Register(L.State, "", luar.Map{
//		"pprint": L.pprint,
//	})
//
//	for pkgname, pkgcontent := range LuaPackages {
//		luar.Register(L.State, pkgname, pkgcontent)
//	}
	return L
}

func RevokeState(L *iState) {
	L.Close()
}

func assert(L *lua.LState) int {
	top := L.GetTop()
	if top == 0 {
		return 0
	}
	for i := 1; i <= top; i++ {
		L.Push(lua.LNumber(i))
	}
	return top
}

func (s *iState) GetRemoteTaskCount() int {
	return s.remote_task_count
}

