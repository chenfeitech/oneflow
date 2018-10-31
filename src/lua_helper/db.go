package lua_helper

import (
	"reflect"
	"time"

	"model"
	"config"
	"database/sql"

	lua "github.com/aarzilli/golua/lua"
	log "github.com/cihub/seelog"
	"github.com/stevedonovan/luar"
)

func init() {
	LuaGlobal["GetProductsByKey"] = model.GetProductsByKey
}

func (l *iState) Lua_db_insert(query string, args ...interface{}) (id int64) {
	result, err := config.GetDBConnect().Exec(query, args...)
	if err != nil {
		l.RaiseError(err.Error())
		return
	}
	id, _ = result.LastInsertId()
	return id
}

func (l *iState) Lua_db_update(query string, args ...interface{}) (rowsAffected int64) {
	result, err := config.GetDBConnect().Exec(query, args...)
	if err != nil {
		l.RaiseError(err.Error())
		return
	}
	rowsAffected, _ = result.RowsAffected()
	return rowsAffected
}

func (l *iState) Lua_db_exec(query string, args ...interface{}) (rowsAffected int64) {
	result, err := config.GetDBConnect().Exec(query, args...)
	if err != nil {
		l.RaiseError(err.Error())
		return
	}
	rowsAffected, _ = result.RowsAffected()
	return rowsAffected
}

func (l *iState) Lua_products_insert(pid string, query string, args ...interface{}) (id int64) {
	products, err := model.GetProductsByKey(pid)
	if err != nil {
		l.RaiseError(err.Error())
		return 0
	}
	if products == nil {
		l.RaiseError("No such pid database.")
		return 0
	}
	conn, err := config.GetFlowDBConnect(products.DBHost, products.DBName)
	if err != nil {
		l.RaiseError(err.Error())
		return 0
	}
	result, err := conn.Exec(query, args...)
	if err != nil {
		l.RaiseError(err.Error())
		return
	}
	id, _ = result.LastInsertId()
	return id
}

func (l *iState) Lua_products_update(pid string, query string, args ...interface{}) (rowsAffected int64) {
	products, err := model.GetProductsByKey(pid)
	if err != nil {
		l.RaiseError(err.Error())
		return 0
	}
	if products == nil {
		l.RaiseError("No such pid database.")
		return 0
	}
	conn, err := config.GetFlowDBConnect(products.DBHost, products.DBName)
	if err != nil {
		l.RaiseError(err.Error())
		return 0
	}
	result, err := conn.Exec(query, args...)
	if err != nil {
		l.RaiseError(err.Error())
		return
	}
	rowsAffected, _ = result.RowsAffected()
	return rowsAffected
}

func (l *iState) Lua_products_exec(pid string, query string, args ...interface{}) (rowsAffected int64) {
	products, err := model.GetProductsByKey(pid)
	if err != nil {
		l.RaiseError(err.Error())
		return 0
	}
	if products == nil {
		l.RaiseError("No such pid database.")
		return 0
	}
	conn, err := config.GetFlowDBConnect(products.DBHost, products.DBName)
	if err != nil {
		l.RaiseError(err.Error())
		return 0
	}
	result, err := conn.Exec(query, args...)
	if err != nil {
		l.RaiseError(err.Error())
		return
	}
	rowsAffected, _ = result.RowsAffected()
	return rowsAffected
}

func QueryArray(L *lua.State) int {
	conn := config.GetDBConnect()
	return Query(L, conn, 0, false)
}

func QueryDict(L *lua.State) int {
	conn := config.GetDBConnect()
	return Query(L, conn, 0, true)
}

func QueryFlowDBArray(L *lua.State) int {
	pid := L.CheckString(1)
	products, err := model.GetProductsByKey(pid)
	if err != nil {
		L.RaiseError(err.Error())
		return 0
	}
	if products == nil {
		L.RaiseError("No such pid database.")
		return 0
	}
	conn, err := config.GetFlowDBConnect(products.DBHost, products.DBName)
	if err != nil {
		L.RaiseError(err.Error())
		return 0
	}
	return Query(L, conn, 1, false)
}

func QueryFlowDBDict(L *lua.State) int {
	pid := L.CheckString(1)
	products, err := model.GetProductsByKey(pid)
	if err != nil {
		L.RaiseError(err.Error())
		return 0
	}
	if products == nil {
		L.RaiseError("No such pid database.")
		return 0
	}
	conn, err := config.GetFlowDBConnect(products.DBHost, products.DBName)
	if err != nil {
		L.RaiseError(err.Error())
		return 0
	}
	return Query(L, conn, 1, true)
}

func Query(L *lua.State, conn *sql.DB, skip_args int, dict bool) int {
	sql := L.CheckString(1 + skip_args)

	top := L.GetTop()
	args := make([]interface{}, 0, top-1-skip_args)

	for i := 2 + skip_args; i <= top; i++ {
		t := L.Type(i)
		switch t {
		case lua.LUA_TSTRING:
			args = append(args, L.ToString(i))
		case lua.LUA_TNUMBER:
			args = append(args, L.ToNumber(i))
		case lua.LUA_TBOOLEAN:
			if L.ToBoolean(i) {
				args = append(args, 1)
			} else {
				args = append(args, 0)
			}
		case lua.LUA_TNIL:
			args = append(args, nil)
		case lua.LUA_TUSERDATA:
			if IsValueProxy(L, i) {
				value, _ := ValueOfProxy(L, i)
				kind := value.Kind()
				if kind == reflect.Ptr {
					kind = value.Elem().Kind()
				}
				if kind == reflect.Interface || kind == reflect.Struct {
					args = append(args, value.Interface())
				} else if kind == reflect.Map {
					args = append(args, value.Interface())
				} else {
					log.Error("Not supported type:" + L.Typename(int(t)))
					L.RaiseError("Not supported type:" + L.Typename(int(t)))
				}
			}
		default:
			log.Error("Not supported type:" + L.Typename(int(t)))
			L.RaiseError("Not supported type:" + L.Typename(int(t)))
		}
	}

	stmt, err := conn.Prepare(sql)
	if err != nil {
		L.RaiseError(err.Error())
		return 0
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		L.RaiseError(err.Error())
		return 0
	}
	defer rows.Close()

	L.NewTable()
	var row_ind int64 = 1

	for rows.Next() {
		L.PushInteger(row_ind)
		L.NewTable()
		columns, _ := rows.Columns()
		values := make([]interface{}, len(columns))
		for i := range values {
			values[i] = new(interface{})
		}
		err := rows.Scan(values...)
		if err != nil {
			L.RaiseError(err.Error())
			return 0
		}
		for i, col := range columns {
			if dict {
				L.PushString(col)
			} else {
				L.PushInteger(int64(i + 1))
			}
			switch val := (*(values[i].(*interface{}))).(type) {
			case int64:
				L.PushInteger(val)
			case float64:
				L.PushNumber(val)
			case nil:
				L.PushNil()
			case []uint8:
				L.PushString(string(val))
			case time.Time:
				log.Debug(val)
				luar.GoToLua(L, reflect.TypeOf(val), reflect.ValueOf(val), false)
			default:
				luar.GoToLua(L, reflect.TypeOf(val), reflect.ValueOf(val), false)
			}
			L.SetTable(-3)
		}
		L.SetTable(-3)
		row_ind++
	}
	return 1
}
