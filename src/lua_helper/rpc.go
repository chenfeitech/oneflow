package lua_helper

import (
	"bytes"
	"encoding/json"

	"middleware/jsonrpc_client"

	log "github.com/cihub/seelog"
)

func (l *iState) Lua_jsonrpc(url string, method string, args map[string]interface{}) (result map[string]interface{}) {
	log.Info(args)
	reply := json.RawMessage{}

	err := jsonrpc_client.Request(url, method, args, &reply)
	if err != nil {
		panic(err)
	}

	var bs = []byte(reply)
	if bytes.Trim(bs, " \t\n")[0] != '{' {
		bs = []byte("{\"result\":" + string(bs) + "}")
	}
	err = json.Unmarshal(bs, &result)
	if err != nil {
		panic(err)
	}
	return
}
