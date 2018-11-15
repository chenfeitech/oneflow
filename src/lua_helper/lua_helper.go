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
	L.LState = lua.NewState()
	L.RemoteExecRecords = make([]*RemoteExecRec, 0)
	L.RemoteExecUseRoot = true

	L.writer = bufio.NewWriter(&L.buffer)

	L.OpenLibs()
	L.iRegister()
	L.Register("gassert", assert)
	L.Register("print", L.print)
	// L.SetGlobal("print", L.NewFunction(l.print))

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

func (l *iState) print(L *lua.LState) int {
	top := L.GetTop()
	for i := 1; i <= top; i++ {
		lv := L.Get(i)
		l.writer.Write(([]byte)(lv.String()))
		l.writer.Write(([]byte)(" "))
	}
	l.writer.Write(([]byte)("\n"))
	return 0
}
