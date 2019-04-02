package server

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"config"
	"lua_helper"
	"middleware/remote_utils"
	"model"
	"utils/helper"
	"web_portal/form"

	log "github.com/cihub/seelog"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
)

var Router = mux.NewRouter()

// handlerFunc adapts a function to an http.Handler.
type handlerFunc func(http.ResponseWriter, *http.Request) error

func (f handlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := f(w, r)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func listHandler(w http.ResponseWriter, r *http.Request) error {
	flow_insts, err := model.FindFlow("")
	if err != nil {
		return err
	}
	return executeTemplate(w, "list", 200, map[string]interface{}{
		"Section":  "list_flow",
		"ListData": flow_insts,
	})
}

func addHandler(w http.ResponseWriter, r *http.Request) error {
	//model_object := &form.JobScheduleForm{}
	model_errors := make(map[string]string)
	return executeTemplate(w, "edit", 200, map[string]interface{}{
		"Section": "add_flow",
		"Errors":  model_errors,
		"Model":   form.Flow{model.Flow{Name: "Test", Description: ""}, nil},
	})
}

func editHandler(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]

	model_object := &form.Flow{}
	flow, err := model.GetFlowByKey(id)
	if err != nil {
		return err
	}
	if flow == nil {
		http.Error(w, "", http.StatusNotFound)
		return nil
	}
	model_object.Flow = *flow
	model_errors := make(map[string]string)
	tasks, err := model.FindTaskByFlowId(id)
	if err != nil {
		return err
	}

	model_object.Tasks = tasks
	return executeTemplate(w, "edit", 200, map[string]interface{}{
		"Section": "Add",
		"Errors":  model_errors,
		"Model":   model_object,
	})
}

func showHandler(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]

	model_object := &form.Flow{}
	flow, err := model.GetFlowByKey(id)
	if err != nil {
		return err
	}
	if flow == nil {
		http.Error(w, "", http.StatusNotFound)
		return nil
	}
	tasks, err := model.FindTaskByFlowId(id)
	if err != nil {
		return err
	}

	model_object.Flow = *flow
	model_object.Tasks = tasks

	return executeTemplate(w, "show", 200, map[string]interface{}{
		"Section": "Add",
		"Model":   model_object,
	})
}

func listFlowInstHandler(w http.ResponseWriter, r *http.Request) error {
	qday := r.URL.Query().Get("date")

	date, err := time.Parse("2006-01-02", qday)
	if err != nil {
		date = time.Now().Add(-24 * time.Hour).Truncate(24 * time.Hour)
	}
	flows, err := model.FindFlowInstByPage(0, 10, "fi.`running_day`=? ORDER BY last_update_time DESC", date.Format("2006-01-02"))
	if err != nil {
		return err
	}
	return executeTemplate(w, "list_flow_inst", 200, map[string]interface{}{
		"Section":  "list_flow_insts",
		"ListData": flows,
		"Date":     date,
	})
}

func showFlowInstHandler(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Parameter ID failed.", http.StatusNotFound)
		return nil
	}

	flow_inst, err := model.GetFlowInstById(id)
	if err != nil {
		return err
	}
	if flow_inst == nil {
		http.Error(w, "", http.StatusNotFound)
		return nil
	}
	task_insts, err := model.FindTaskInstByFlow(flow_inst.Id, flow_inst.FlowId)
	if err != nil {
		return err
	}

	return executeTemplate(w, "show_flow_inst", 200, map[string]interface{}{
		"Section":  "Show",
		"Model":    flow_inst,
		"ListData": task_insts,
	})
}

func rerunFlowInstHandler(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Parameter ID failed.", http.StatusNotFound)
		return nil
	}

	task_id := vars["task_id"]
	single := true
	if vars["single"] != "true" {
		single = false
	}

	flow_inst, err := model.GetFlowInstById(id)
	if err != nil {
		return log.Error("Find flow instance failed:", err)
	}
	if flow_inst == nil {
		return log.Error("Flow instance not exists.")
	}

	task_inst, err := model.GetTaskInstById(id, task_id)
	if err != nil {
		return log.Error("Find task instance failed:", err)
	}
	if task_inst == nil {
		return log.Error("Task instance not exists.")
	}

	task, err := model.GetTaskById(task_inst.FlowId, task_id)
	if err != nil {
		return log.Error("Find next state failed:", err)
	}
	if task == nil {
		return log.Error("Task not exists.")
	}

	if single {
		err = model.ResetTaskInstState(id, task.Id)
	} else {
		err = model.ResetTaskInstStateSinceOrderId(id, task.FlowId, task.OrderId)
	}
	if err != nil {
		return log.Error("Update task state failed:", err)
	}

	err = lua_helper.StartFlowTask(flow_inst.PId, flow_inst.Key, task, "", flow_inst.RunningDay, nil, task_inst)
	if err != nil {
		return log.Error("Start task failed:", err)
	}

	http.Redirect(w, r, r.Referer(), 302)
	return nil
}

func terminalHandler(w http.ResponseWriter, r *http.Request) error {
	body, _ := ioutil.ReadAll(r.Body)
	script := (string)(body)

	log.Debug("Run script:", script)

	L := lua_helper.GetState()
	defer lua_helper.RevokeState(L)

	err := L.DoString("in_terminal=true")
	if err != nil {
		w.Write(([]byte)(err.Error()))
	}
	err = L.DoString(script)
	if err != nil {
		w.Write(([]byte)(err.Error()))
	}

	output := L.GetOutput()
	log.Info("Script output:", output)
	w.Write(([]byte)(output))

	return nil
}

func logHandler(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	host := vars["host"]
	uuid := vars["uuid"]
	date := vars["date"]

	server, err := model.GetServerByHost(host)
	if err != nil {
		return log.Error("Find host "+host+" failed:", err)
	}
	if server == nil {
		return log.Error("Host " + host + " not found in our database!")
	}
	ssh_client, err := remote_utils.Connect(fmt.Sprintf("%s:%d", server.Host, server.Port), *server.Username, *server.Password)
	if err != nil {
		return log.Error("Connect to host "+host+" failed:", err)
	}
	defer ssh_client.Close()

	path := fmt.Sprintf("%s/proc/%s/%s", *config.ServerRoot, strings.Replace(date, "-", "", -1), uuid)

	log_filenames := []string{path + "/cmd", path + "/out.log", path + "/err.log"}

	log_datas := remote_utils.ReadFiles(ssh_client, log_filenames)

	log_strings := make([]string, len(log_datas))
	for i, log_data := range log_datas {
		switch val := log_data.(type) {
		case []byte:
			log_strings[i] = (string)(val)
		case error:
			log_strings[i] = "<p class='text-danger'>" + val.Error() + "</p>"
		}
	}

	return executeTemplate(w, "log", 200, map[string]interface{}{
		"Cmdline": log_strings[0],
		"Output":  log_strings[1],
		"Error":   log_strings[2],
	})
}

var (
	fileserver_path = flag.String("fileserver_path", *config.ServerRoot+"/fileserver", "Flow central file server.")
)

func fileInfoHandler(w http.ResponseWriter, r *http.Request) error {
	file_path := path.Clean(path.Join(*fileserver_path, r.URL.String()))

	if !strings.HasPrefix(file_path, *fileserver_path) {
		http.Error(w, "", http.StatusNotFound)
		return nil
	}
	file, err := os.Open(file_path)
	if err != nil {
		log.Error("Open file ", file_path, " failed:", err)
		http.Error(w, "", http.StatusNotFound)
		return nil
	}
	fstat, err := file.Stat()
	if err != nil {
		log.Error("Stat file ", file_path, " failed:", err)
		http.Error(w, "", http.StatusNotFound)
		return nil
	}
	if fstat.IsDir() {
		http.Error(w, "", http.StatusNotFound)
		return nil
	}

	w.Header().Add("MD5", helper.GetFileMd5(file))
	w.Header().Add("File", file_path)
	w.Header().Add("Mode", fmt.Sprintf("%o", (uint32)(fstat.Mode())))
	return nil
}

func fileGetHandler(w http.ResponseWriter, r *http.Request) error {
	file_path := path.Clean(path.Join(*fileserver_path, r.URL.String()))

	if !strings.HasPrefix(file_path, *fileserver_path) {
		http.Error(w, "", http.StatusNotFound)
		return nil
	}
	file, err := os.Open(file_path)
	if err != nil {
		log.Error("Open file ", file_path, " failed:", err)
		http.Error(w, "", http.StatusNotFound)
		return nil
	}
	fstat, err := file.Stat()
	if err != nil {
		log.Error("Stat file ", file_path, " failed:", err)
		http.Error(w, "", http.StatusNotFound)
		return nil
	}
	if fstat.IsDir() {
		http.Error(w, "", http.StatusNotFound)
		return nil
	}
	w.Header().Add("Mode", fmt.Sprintf("0%o", (uint32)(fstat.Mode())))
	w.Header().Add("Size", fmt.Sprint(fstat.Size()))
	io.Copy(w, file)
	return nil
}

func init() {
	r := Router
	r.StrictSlash(false)

	r.Handle("/data_flow/terminal", handlerFunc(terminalHandler))

	//r.Handle("/", handlerFunc(homeHandler))
	r.Handle("/data_flow/", handlerFunc(listFlowInstHandler))
	r.Handle("/data_flow/flows", handlerFunc(listHandler)).Name("list_flow")
	r.Handle("/data_flow/flow", handlerFunc(addHandler)).Name("add_flow")
	r.Handle("/data_flow/flow/{id}", handlerFunc(showHandler)).Name("show_flow")
	r.Handle("/data_flow/flow/{id}/edit", handlerFunc(editHandler)).Name("edit_flow")
	r.Handle("/data_flow/flow_insts", handlerFunc(listFlowInstHandler)).Name("list_flow_insts")
	r.Handle("/data_flow/flow_inst/{id:[0-9]+}", handlerFunc(showFlowInstHandler)).Name("show_flow_inst")
	r.Handle("/data_flow/flow_inst/rerun/{id:[0-9]+}/{task_id}/{single}", handlerFunc(rerunFlowInstHandler)).Name("rerun_flow_inst")
	r.Handle("/data_flow/flow_inst/log/{date}/{host}/{uuid}", handlerFunc(logHandler)).Name("log")

	r.PathPrefix("/data_flow/static/").Handler(http.StripPrefix("/data_flow/static/", http.FileServer(http.Dir("./static/")))).Name("static")

	jsonRPC := rpc.NewServer()
	jsonCodec := NewCodec()
	jsonRPC.RegisterCodec(jsonCodec, "application/json")
	jsonRPC.RegisterCodec(jsonCodec, "application/json; charset=UTF-8") // For firefox 11 and other browsers which append the charset=UTF-8
	service := new(FlowService)
	jsonRPC.RegisterService(service, "")

	r.Handle("/data_flow/api", jsonRPC).Name("api")

	r.Methods("HEAD").Handler(http.StripPrefix("/data_flow/file/", handlerFunc(fileInfoHandler)))
	r.Methods("GET").Handler(http.StripPrefix("/data_flow/file/", handlerFunc(fileGetHandler)))
	startEventLoop(10, service)
}
