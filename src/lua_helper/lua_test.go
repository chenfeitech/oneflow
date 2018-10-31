package lua_helper

import (
	"bufio"
	"bytes"
	"testing"
)

func Test_Remote(t *testing.T) {
	L := GetState()
	defer RevokeState(L)

	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	RedirectOutput(L, writer)

	err := L.DoString("print(remote_exec('10.32.74.83', 'ls', '/'))")

	writer.Flush()
	t.Log("-----------", buffer.String())
	if err != nil { //try a unit test on function
		t.Log(err)
		t.Error("remote_exec") // 如果不是如预期的那么就报错
	} else {
		t.Log("remote_exec") //记录一些你期望记录的信息
	}
}

func Test_Remote_Assert(t *testing.T) {
	L := GetState()
	defer RevokeState(L)

	err := L.DoString("print(gassert(remote_exec('10.32.74.83', 'ls')))")

	if err != nil { //try a unit test on function
		t.Log(err)
		t.Error("remote_exec") // 如果不是如预期的那么就报错
	} else {
		t.Log("remote_exec") //记录一些你期望记录的信息
	}
}

func Test_Remote_With_Error(t *testing.T) {
	L := GetState()
	defer RevokeState(L)

	err := L.DoString("print(remote_exec('', 'ls', '/'))")

	if err != nil { //try a unit test on function
		t.Log(err)
		t.Error("remote_exec") // 如果不是如预期的那么就报错
	} else {
		t.Log("remote_exec") //记录一些你期望记录的信息
	}
}

func Test_Remote_Assert_With_Error(t *testing.T) {
	L := GetState()
	defer RevokeState(L)

	err := L.DoString("print(gassert(remote_exec('', 'ls')))")

	if err == nil { //try a unit test on function
		t.Error("remote_exec") // 如果不是如预期的那么就报错
	} else {
		t.Log(err)
		t.Log("remote_exec") //记录一些你期望记录的信息
	}
}
