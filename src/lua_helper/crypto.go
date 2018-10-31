package lua_helper

import (
	"utils/helper"
)

func (l *iState) Lua_encrypt(text string) (string, error) {
	return helper.Encrypt(text)
}

func (l *iState) Lua_decrypt(cryptstring string) (string, error) {
	return helper.Decrypt(cryptstring)
}
