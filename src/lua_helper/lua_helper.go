package lua_helper

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"lua_helper/packages"
	"math"
	"model"
	"reflect"
	"strings"

	lua "github.com/aarzilli/golua/lua"
	log "github.com/cihub/seelog"
	"github.com/stevedonovan/luar"
)

var (
	LuaGlobal   = luar.Map{}
	LuaPackages = make(map[string]luar.Map)
)

type iState struct {
	*lua.State
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
	funcs := luar.Map{}
	for name, fun := range LuaGlobal {
		funcs[name] = fun
	}

	lval := reflect.ValueOf(l)
	ltype := reflect.TypeOf(l)
	for i := 0; i < ltype.NumMethod(); i++ {
		method := ltype.Method(i)
		if strings.HasPrefix(method.Name, "Lua_") {
			fun_name := strings.TrimPrefix(method.Name, "Lua_")
			log.Info("Register ", fun_name)
			funcs[fun_name] = lval.MethodByName(method.Name).Interface()
		}
	}
	luar.Register(l.State, "", funcs)
}

func (l *iState) GetOutput() string {
	if buf_writer, ok := l.writer.(*bufio.Writer); ok {
		buf_writer.Flush()
	}
	return l.buffer.String()
}

func (s *iState) Lua_get_package(name string) string {
	return packages.Packages[name]
}

// TODO
func GetStateByPId(pid string) *iState {
	return GetState()
}

func GetState() *iState {
	L := &iState{}
	L.State = luar.Init()
	L.RemoteExecRecords = make([]*RemoteExecRec, 0)
	L.RemoteExecUseRoot = true

	L.writer = bufio.NewWriter(&L.buffer)

	L.OpenLibs()
	L.iRegister()
	L.Register("print", L.print)
	L.Register("gassert", assert)
	L.Register("db_query_array", QueryArray)
	L.Register("db_query_dict", QueryDict)
	L.Register("flow_query_array", QueryFlowDBArray)
	L.Register("flow_query_dict", QueryFlowDBDict)

	luar.Register(L.State, "", luar.Map{
		"pprint": L.pprint,
	})

	for pkgname, pkgcontent := range LuaPackages {
		luar.Register(L.State, pkgname, pkgcontent)
	}

	fmt.Println(L.DoString(`
function table_print (lua_table, indent)
	if type(lua_table) ~= "table" then 
		return 
	end
	indent = indent or 0
	print("{")
	for k, v in pairs(lua_table) do
		if type(k) == "string" then
			k = string.format("%q", k)
		end
		local szSuffix = ""
		if type(v) == "table" then
			szSuffix = "{"
		end
		local szPrefix = string.rep("    ", indent)
		formatting = szPrefix.."["..k.."]".." = "..szSuffix
		if type(v) == "table" then
			print(formatting)
			table_print(v, indent + 1)
			print(szPrefix.."},")
		else
			if type(v) == "string" then
				print(formatting..string.format("%q", v)..",")
			else
				print(formatting, v, ",")
			end
		end
	end
	print("}")
end

__xml_decode = xml_decode
function xml_decode(xml)
	return luar.map2table(__xml_decode(xml))
end

function xml_node_build (tag, data)
	if type(tag) ~= "string" then
		return nil
	end
    if type(data) ~= "table" then
		return nil
	end
    local node = { XMLName={Local=tag}, Nodes={}}
  	local array_element_count = 0
	for k, v in pairs(data) do
    	if type(v) ~= "table" then
          if type(k) ~= "string" then
              k = tostring(k)
          end
          if type(v) ~= "string" then
              v = tostring(v)
          end
      	  table.insert(node.Nodes,  { XMLName={Local=k}, Data=v})
        else
          if type(k) == "number" then
              table.insert(node.Nodes, xml_node_build(tag, v))
        	  array_element_count = array_element_count + 1
          else
            if type(k) ~= "string" then
                k = tostring(k)
            end
            local sub_nodes = xml_node_build(k, v)
            if sub_nodes.XMLName then
            	table.insert(node.Nodes, sub_nodes)
            else
          		for i, sub_node in ipairs(sub_nodes) do
            		table.insert(node.Nodes, sub_node)
            	end
            end
          end
        end
    end
    if array_element_count == #data and array_element_count>0 then
        return node.Nodes
    else
    	return node
    end
end

__xml_encode = xml_encode
function xml_encode(tag, root)
	return __xml_encode(xml_node_build(tag, root))
end

package.preload['template'] = function ()
    return loadstring(get_package('template'))()
end

	`))
	return L
}

func RevokeState(L *iState) {
	L.Close()
}

// func pcall(L *lua.State) (ret int) {
// 	fmt.Println("---------------------pcall   Top:", L.GetTop())
// 	defer func() {
// 		if err2 := recover(); err2 != nil {
// 			fmt.Println("---------------------pcall   recover:", err2)
// 			luar.GoToLua(L, typeof(err2), valueOf(err2), false)
// 			ret = 1
// 			return
// 		}
// 	}()

// 	L.MustCall(1, 0)
// 	fmt.Println("---------------------pcall   Top:", L.GetTop())
// 	pos := L.GetTop() - ret + 1
// 	L.PushNil()
// 	L.Insert(pos)

// 	for i := 0; i < L.GetTop(); i++ {
// 		fmt.Println(L.Typename(int(L.Type(i + 1))))
// 	}
// 	return L.GetTop()
// }

func assert(L *lua.State) int {
	top := L.GetTop()
	if top == 0 {
		return 0
	}
	for i := 1; i <= top; i++ {
		L.PushValue(i)
	}
	return top
}

func (l *iState) pprint(i int) {
	fmt.Fprintln(l.writer, i)
}

func (l *iState) print(L *lua.State) int {
	top := L.GetTop()
	for i := 1; i <= top; i++ {
		t := L.Type(i)
		if L.IsGoStruct(i) {
			//fmt.Println("IsGoStruct")
			l.writer.Write(([]byte)(fmt.Sprintf("%+v", L.ToGoStruct(i))))
			continue
		}
		if IsValueProxy(L, i) {
			//fmt.Println("IsValueProxy")
			value, _ := ValueOfProxy(L, i)
			kind := value.Kind()
			if kind == reflect.Ptr {
				kind = value.Elem().Kind()
			}
			if value.CanInterface() {
				json_bytes, _ := json.MarshalIndent(value.Interface(), "", "  ")
				fmt.Fprintln(l.writer, string(json_bytes))
			} else {
				log.Warn("Lua vm print value proxy failed, type:", kind)
			}
			continue
		}
		switch t {
		case lua.LUA_TSTRING:
			l.writer.Write(([]byte)(L.ToString(i)))
		case lua.LUA_TNUMBER:
			val := L.ToNumber(i)
			if val == math.Floor(val) {
				l.writer.Write(([]byte)(fmt.Sprint(int(val))))
			} else {
				l.writer.Write(([]byte)(fmt.Sprint(val)))
			}
		case lua.LUA_TBOOLEAN:
			if L.ToBoolean(i) {
				l.writer.Write(([]byte)("true"))
			} else {
				l.writer.Write(([]byte)("false"))
			}
		case lua.LUA_TNIL:
			l.writer.Write(([]byte)("nil"))

		case lua.LUA_TTABLE:
			L.GetField(lua.LUA_GLOBALSINDEX, "table_print")
			L.PushValue(i)
			L.Call(1, 0)

		default:
			log.Error("Print not support type:", L.Typename(i))
			l.writer.Write(([]byte)(L.ToString(i)))
		}
		l.writer.Write(([]byte)(" "))
	}
	l.writer.Write(([]byte)("\n"))
	return 0
}

func (s *iState) GetRemoteTaskCount() int {
	return s.remote_task_count
}

func CopySliceToTable(L *lua.State, obj interface{}) int {
	vslice := valueOf(obj)
	if vslice.IsValid() && vslice.Type().Kind() == reflect.Slice {
		n := vslice.Len()
		L.CreateTable(n, 0)
		for i := 0; i < n; i++ {
			L.PushInteger(int64(i + 1))
			v := vslice.Index(i)

			if isNil(v) {
				v = nullv
			}
			sv := v
			if sv.Kind() == reflect.Ptr {
				sv = sv.Elem()
			}
			if sv.Kind() == reflect.Struct {
				if ev, ok := sv.Interface().(error); ok {
					L.PushString(ev.Error())
				} else if lv, ok := sv.Interface().(*luar.LuaObject); ok {
					lv.Push()
				} else {
					if (sv.Kind() == reflect.Ptr || sv.Kind() == reflect.Interface) && !sv.Elem().IsValid() {
						L.PushNil()
					} else {
						L.PushGoStruct(sv.Interface())
					}
				}
			} else {
				luar.GoToLua(L, nil, v, true)
			}
			L.SetTable(-3)
		}
		return 1
	} else {
		L.PushNil()
		L.PushString("not a slice!")
	}
	return 2
}

var (
	tslice    = typeof((*[]interface{})(nil))
	tmap      = typeof((*map[string]interface{})(nil))
	null      = Null(0)
	nullv     = valueOf(null)
	nullables = map[reflect.Kind]bool{
		reflect.Chan:      true,
		reflect.Func:      true,
		reflect.Interface: true,
		reflect.Map:       true,
		reflect.Ptr:       true,
		reflect.Slice:     true,
	}
)

func isNil(val reflect.Value) bool {
	kind := val.Type().Kind()
	return nullables[kind] && val.IsNil()
}

func typeof(v interface{}) reflect.Type {
	return reflect.TypeOf(v).Elem()
}

var valueOf = reflect.ValueOf

type Null int

func IsValueProxy(L *lua.State, idx int) bool {
	res := false
	if L.IsUserdata(idx) {
		L.GetMetaTable(idx)
		if !L.IsNil(-1) {
			L.GetField(-1, "luago.value")
			res = !L.IsNil(-1)
			L.Pop(1)
		}
		L.Pop(1)
	}
	return res
}

type valueProxy struct {
	value reflect.Value
	t     reflect.Type
}

func ValueOfProxy(L *lua.State, idx int) (reflect.Value, reflect.Type) {
	vp := (*valueProxy)(L.ToUserdata(idx))
	return vp.value, vp.t
}
