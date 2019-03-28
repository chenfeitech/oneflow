package lua_helper

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
    // "encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"config"
	"utils/helper"
	"model"
	"middleware/remote_utils"

	"code.google.com/p/uuid"
	"github.com/yuin/gopher-lua"
	// "layeh.com/gopher-luar"
	log "github.com/cihub/seelog"
	"golang.org/x/crypto/ssh"
)

type RemoteExecRec struct {
	Host string
	Uuid string
}

func init() {
	md5_sum := md5.New()
	bin, err := os.Open(local_task_agent_path)
	if err != nil {
		log.Error("Calc local agent ", local_task_agent_path, " failed:", err)
		os.Exit(1)
	} else {
		io.Copy(md5_sum, bin)
		task_agent_md5 = hex.EncodeToString(md5_sum.Sum(nil))
		log.Info("Local agent md5:", task_agent_md5)
	}
}

var (
	g_mutex     = &sync.Mutex{}
	host_mutexs = make(map[string]*sync.Mutex, 0)
)

var (
	task_agent_filename = "flow_agent"

	remote_task_agent_dir  = *config.ServerRoot + "/agent/"
	remote_task_agent_path = remote_task_agent_dir + task_agent_filename
)

var (
	local_task_agent_path = func() string {
		if runtime.GOOS == "darwin" {
			return *config.ServerRoot + "/bin/"
		} else {
			return *config.ServerRoot + "/bin/"
		}
	}() + task_agent_filename
	task_agent_md5 string

	not_download_pattern = []string{
		"/bin/*",
		"/usr/local/bin/*",
		"/sbin/*",
	}
)

func (l *iState) Remote_init(L *lua.LState) {
	L.SetGlobal("hello_world", L.NewFunction(l.Lua_hello_world))
	L.SetGlobal("remote_exec_set_env", L.NewFunction(l.Lua_remote_exec_set_env))
	L.SetGlobal("remote_exec", L.NewFunction(l.Lua_remote_exec))
}

func (l *iState) Lua_hello_world(L *lua.LState) int {
	num := L.GetTop()
	name := L.ToString(1)
	fmt.Println("num: ", num, " hello world: ", name)
	L.Pop(num)
	return 0
}

func (l *iState) Lua_remote_exec_set_env(L *lua.LState) int {
	num := L.GetTop()
	if num != 2 {
		fmt.Println("request args num should be 2, but ", num)
		return 1
	}
	name := L.CheckString(1)
	value := L.CheckString(1)
	l.remote_exec_set_env(name, value)

	return 0;
}

func (l *iState) remote_exec_set_env(name string, value string) {
	if l.RemoteExecEnv == nil {
		l.RemoteExecEnv = make(map[string]string)
	}
	l.RemoteExecEnv[name] = fmt.Sprint(value)
}

func (l *iState) Lua_remote_exec(L *lua.LState) int {
	num := L.GetTop()
	if num < 2 {
		fmt.Println("request args num should be more then2, but ", num)
		return 1
	}
	host := L.CheckString(1)
	program := L.CheckString(2)
	args := ""
	if num > 2 {
		args = L.CheckString(3)
		for i:=4; i <= num; i++ {
			args = args + ", " + L.CheckString(i)
		}
	}
	guid, output, err := l.remote_exec(host, program, args)
	L.Push(lua.LString(guid))
	L.Push(lua.LString(output))
	// L.SetGlobal("guid", lua.LString(guid))
	// L.SetGlobal("output", lua.LString(output))
	// L.SetGlobal("err", lua.LString(err))
	fmt.Println("xxx guid: ", guid, " output: ", output, " err: ", err)

	return 2;
}

func (l *iState) remote_exec(host string, program string, args ...interface{}) (guid string, output string, err error) {
	if len(host) == 0 {
		panic(errors.New("Must give remote host to exec program!"))
	}
	guid = uuid.New()

	log.Info(append([]interface{}{"remote_exec " + guid + " " + host + " " + program}, args...)...)

	cmd := remote_task_agent_path + " run " + "-uuid " + guid + " -p \"" + program + "\" "
	// cmd := []string{remote_task_agent_path, "run", "-uuid", guid, "-p", program}
	fargs := make([]string, len(args), len(args))
	for i := 0; i < len(args); i++ {
		switch val := args[i].(type) {
		case reflect.Value:
			arg := strings.Replace(fmt.Sprint(val.Interface()), "\"", "\\\"", -1)
			fargs[i] = fmt.Sprintf("-a \"%v\"", arg)
		default:
			arg := strings.Replace(fmt.Sprint(args[i]), "\"", "\\\"", -1)
			fargs[i] = fmt.Sprintf("-a \"%v\"", arg)
		}
		// cmd = append(cmd, "-a", conv.String(arg))
	}
	cmd = cmd + strings.Join(fargs, " ")

	must_download := true
	for _, pattern := range not_download_pattern {
		if match, _ := path.Match(pattern, program); match {
			must_download = false
			break
		}
	}
	env := ExecEnv{}

	if must_download {
		env["REQUIRE_FILES"] = program
	}
	env["SERVER_HOST"] = *config.ServerHost
	env["SERVER_PORT"] = *config.ServerPort
	env["SERVER_PREFIX"] = "data_flow"

	if l.Task != nil {
		env["FLOW_ID"] = l.Task.FlowId
		env["TASK_ID"] = l.Task.Id
	}
	if l.FlowInstance != nil {
		env["FLOW_INST_ID"] = strconv.Itoa(l.FlowInstance.Id)
		env["PID"] = l.FlowInstance.PId
		env["KEY"] = l.FlowInstance.Key
		env["DATE"] = l.FlowInstance.RunningDay.Format("2006-01-02")
	}

	for k, v := range l.RemoteExecEnv {
		env[k] = v
	}

	cmdline := env.Encode() + " " + cmd

	output, err = RemoteExec(host, cmdline)
	// output, err = RemoteExec(l.RemoteExecUseRoot, host, env, cmd...)
	if err != nil {
		log.Error(err)
		panic(err)
	}

	l.remote_task_count++
	l.RemoteExecRecords = append(l.RemoteExecRecords, &RemoteExecRec{Host: host, Uuid: guid})
	return guid, output, nil
}

func RemoteKill(host string, day time.Time, uuid string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
		}
	}()

	if len(host) == 0 {
		panic(errors.New("Must give remote host to exec program!"))
	}

	cmd := remote_task_agent_path + " kill " + "-uuid " + uuid + " -date " + day.Format("20060102")
	output, err := RemoteExec(host, cmd)
	// cmd := []string{remote_task_agent_path, "kill", "-uuid", uuid, "-date", day.Format("20060102")}
	// output, err := RemoteExec(useRoot, host, nil, cmd...)
	log.Info("Remote kill output:", output)
	return err
}
/*
func RemoteExec1(useRoot bool, host string, env map[string]interface{}, cmdArgs ...string) (output string, err error) {
	g_mutex.Lock()
	host_mutex := host_mutexs[host]
	if host_mutex == nil {
		host_mutex = &sync.Mutex{}
		host_mutexs[host] = host_mutex
	}
	g_mutex.Unlock()

	host_mutex.Lock()
	defer host_mutex.Unlock()

	if len(host) == 0 {
		panic(errors.New("Must give remote host to exec program!"))
	}

	if env == nil {
		env = map[string]interface{}{}
	}

	envBytes, err := json.Marshal(env)
	if err != nil {
		return "", err
	}
	cmdArgsBytes, err := json.Marshal(cmdArgs)
	if err != nil {
		return "", err
	}

	scriptEnv := "var Env=" + string(envBytes) + ";\n"
	scriptCmdArgs := "var CmdArgs=" + string(cmdArgsBytes) + ";\n"
	script := scriptEnv + scriptCmdArgs + remoteExecScript
	log.Info("script:", script)
	var reply interface{}
	if useRoot {
		err = microservice.Request(host, "Script.Exec", script, &reply)
	} else {
		err = microservice.SafeRequest(host, "Script.Exec", script, &reply)
	}
	return conv.String(reply), err
}
*/
func RemoteExec(host string, cmdline string) (output string, err error) {
	g_mutex.Lock()
	host_mutex := host_mutexs[host]
	if host_mutex == nil {
		host_mutex = &sync.Mutex{}
		host_mutexs[host] = host_mutex
	}
	g_mutex.Unlock()

	host_mutex.Lock()
	defer host_mutex.Unlock()

	if len(host) == 0 {
		panic(errors.New("Must give remote host to exec program!"))
	}

	server, err := model.GetServerByHost(host)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	if server == nil {
		panic(errors.New("Host " + host + " not found in our database!"))
	}

	if server.Username == nil {
		server.Username = new(string)
		*server.Username = "root"
	}

	var password string
	if len(server.CryptoPassword) > 0 {
		password, err = helper.Decrypt(server.CryptoPassword)
		if err != nil {
			log.Error("Decrypt password ", server.CryptoPassword, " failed:", err)
		}
	}

	if len(password) == 0 {
		if server.Password == nil {
			password = "data2014"
		} else {
			password = *server.Password
		}
	}

	ssh_client, err := remote_utils.Connect(fmt.Sprintf("%s:%d", server.Host, server.Port), *server.Username, password)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	defer ssh_client.Close()

	checksum, err := get_remote_agent_checksum(ssh_client)
	if err != nil {
		panic(err)
	}

	if checksum != task_agent_md5 {
		log.Info("Remote agent md5:", checksum, " local agent md5:", task_agent_md5)
		err := upload_task_agent(ssh_client)
		if err != nil {
			log.Error(err)
			panic(err)
		}
	}

	session, err := ssh_client.NewSession()
	if err != nil {
		log.Error(err)
		panic(err)
	}

	retries := 0
	client := ssh_client
	for {
		if client == nil {
			client, err = remote_utils.Connect(fmt.Sprintf("%s:%d", server.Host, server.Port), *server.Username, password)
			if err == nil {
				defer client.Close()
			}
		}
		if client != nil {
			log.Info("Run ", cmdline, " on ", server.Host)
			output_bytes, err := session.CombinedOutput(cmdline)
			output = (string)(output_bytes)

			output_lines := strings.Split(output, "\n")
			for i := 0; i < len(output_lines); i++ {
				if len(strings.Trim(output_lines[i], " ")) != 0 {
					output_lines[i] = "[" + server.Host + "] " + output_lines[i]
				}
			}
			output = strings.Join(output_lines, "\n")

			if err != nil {
				log.Error(err)
				if _, ok := err.(*ssh.ExitError); ok {
					if strings.Contains(output, task_agent_filename+": Text file busy") {
						log.Info("Task agent file busy.")
					} else {
						panic(fmt.Errorf("%v\n Output: %s", err, output))
					}
				} else {
					client = nil
				}
			} else {
				break
			}
		}
		if retries < 5 {
			time.Sleep(2 * time.Second)
			retries++
		} else {
			panic(fmt.Errorf("%v\n Output: %s", err, output))
		}
	}
	return output, nil
}

func get_remote_agent_checksum(client *ssh.Client) (string, error) {
	log.Info("Get remote agent md5.")
	session, err := client.NewSession()
	if err != nil {
		log.Error(err)
		return "", err
	}

	output, err := session.Output(remote_task_agent_path + " " + "checksum")
	log.Info(remote_task_agent_path+" checksum ", (string)(output), err)
	if err != nil {
		log.Error(err)
		if _, ok := err.(*ssh.ExitError); ok {
			return "", nil
		} else {
			return "", err
		}
	}
	md5 := strings.Trim((string)(output), " \r\n\t")

	return md5, nil
}

func upload_task_agent(client *ssh.Client) error {
	log.Info("Upload remote agent.")

	bin, err := os.Open(local_task_agent_path)
	if err != nil {
		log.Error("Open local agent ", local_task_agent_path, " failed:", err)
		return err
	}
	err = remote_utils.Upload(client, bin, remote_task_agent_dir, task_agent_filename, 755)

	return err
}

type ExecEnv map[string]string

func (e *ExecEnv) Encode() string {
	var buf bytes.Buffer
	for k, v := range *e {
		buf.WriteString(k)
		buf.WriteString("=\"")
		buf.WriteString(strings.Replace(v, "\"", "\\\"", -1))
		buf.WriteString("\" ")
	}
	return buf.String()
}
