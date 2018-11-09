package lua_helper

import (
	"bufio"
	"bytes"
	// "encoding/json"
	"fmt"
	"io"
	// "math"
	"model"
	// "reflect"
	// "strings"

	"github.com/yuin/gopher-lua"
	// lua "github.com/aarzilli/golua/lua"
	log "github.com/cihub/seelog"
	// "github.com/stevedonovan/luar"
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
	l.Do_init(l.LState)
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

// func (l *iState) pprint(i int) {
// 	fmt.Fprintln(l.writer, i)
// }
//
// func (l *iState) print(L *lua.State) int {
// 	top := L.GetTop()
// 	for i := 1; i <= top; i++ {
// 		t := L.Type(i)
// 		if L.IsGoStruct(i) {
// 			//fmt.Println("IsGoStruct")
// 			l.writer.Write(([]byte)(fmt.Sprintf("%+v", L.ToGoStruct(i))))
// 			continue
// 		}
// 		if IsValueProxy(L, i) {
// 			//fmt.Println("IsValueProxy")
// 			value, _ := ValueOfProxy(L, i)
// 			kind := value.Kind()
// 			if kind == reflect.Ptr {
// 				kind = value.Elem().Kind()
// 			}
// 			if value.CanInterface() {
// 				json_bytes, _ := json.MarshalIndent(value.Interface(), "", "  ")
// 				fmt.Fprintln(l.writer, string(json_bytes))
// 			} else {
// 				log.Warn("Lua vm print value proxy failed, type:", kind)
// 			}
// 			continue
// 		}
// 		switch t {
// 		case lua.LUA_TSTRING:
// 			l.writer.Write(([]byte)(L.ToString(i)))
// 		case lua.LUA_TNUMBER:
// 			val := L.ToNumber(i)
// 			if val == math.Floor(val) {
// 				l.writer.Write(([]byte)(fmt.Sprint(int(val))))
// 			} else {
// 				l.writer.Write(([]byte)(fmt.Sprint(val)))
// 			}
// 		case lua.LUA_TBOOLEAN:
// 			if L.ToBoolean(i) {
// 				l.writer.Write(([]byte)("true"))
// 			} else {
// 				l.writer.Write(([]byte)("false"))
// 			}
// 		case lua.LUA_TNIL:
// 			l.writer.Write(([]byte)("nil"))
//
// 		case lua.LUA_TTABLE:
// 			L.GetField(lua.LUA_GLOBALSINDEX, "table_print")
// 			L.PushValue(i)
// 			L.Call(1, 0)
//
// 		default:
// 			log.Error("Print not support type:", L.Typename(i))
// 			l.writer.Write(([]byte)(L.ToString(i)))
// 		}
// 		l.writer.Write(([]byte)(" "))
// 	}
// 	l.writer.Write(([]byte)("\n"))
// 	return 0
// }
//
func (s *iState) GetRemoteTaskCount() int {
	return s.remote_task_count
}
//
// var (
// 	tslice    = typeof((*[]interface{})(nil))
// 	tmap      = typeof((*map[string]interface{})(nil))
// 	null      = Null(0)
// 	nullv     = valueOf(null)
// 	nullables = map[reflect.Kind]bool{
// 		reflect.Chan:      true,
// 		reflect.Func:      true,
// 		reflect.Interface: true,
// 		reflect.Map:       true,
// 		reflect.Ptr:       true,
// 		reflect.Slice:     true,
// 	}
// )
//
// func isNil(val reflect.Value) bool {
// 	kind := val.Type().Kind()
// 	return nullables[kind] && val.IsNil()
// }
//
// func typeof(v interface{}) reflect.Type {
// 	return reflect.TypeOf(v).Elem()
// }
//
// var valueOf = reflect.ValueOf
//
// type Null int
//
// func IsValueProxy(L *lua.State, idx int) bool {
// 	res := false
// 	if L.IsUserdata(idx) {
// 		L.GetMetaTable(idx)
// 		if !L.IsNil(-1) {
// 			L.GetField(-1, "luago.value")
// 			res = !L.IsNil(-1)
// 			L.Pop(1)
// 		}
// 		L.Pop(1)
// 	}
// 	return res
// }
//
// type valueProxy struct {
// 	value reflect.Value
// 	t     reflect.Type
// }
//
// func ValueOfProxy(L *lua.State, idx int) (reflect.Value, reflect.Type) {
// 	vp := (*valueProxy)(L.ToUserdata(idx))
// 	return vp.value, vp.t
// }
