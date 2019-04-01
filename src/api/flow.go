package api

import (
	// "fmt"
	"io/ioutil"
	"errors"
	"encoding/json"

	"model"
	"utils/conv"
	// "model/jobmanager"

	log "github.com/cihub/seelog"
	"github.com/go-sql-driver/mysql"
)

func init() {
	SetRouterRegister(func(router *RouterGroup) {
		flowRouteGroup := router.Group("/oneflow/")
		flowRouteGroup.StdPOST("AddFlow", AddFlow)
		flowRouteGroup.StdPOST("UpdateFlow", UpdateFlow)
		flowRouteGroup.StdGET("GetFlows", GetFlows)
	})
}

type RpcQuestArgs struct {
	Id	int
	Method	string
	Params []FlowDataArgs
}

type FlowDataArgs struct {
	Id            string
	Name          string
	Description   string
	Creator       string
	StartupScript string
	Tasks         []TaskData
	DeleteTaskIds []string
}

type TaskData struct {
	Id          string
	Name        string
	Script      string
	MaxRetries  interface{} `json:"max_retries"`
	Description string
}
/*
{"id":2,"method":"FlowService.UpdateFlow",
"params":[{"id":"NEW_FLOW","name":"新流程","creator":"helight","description":"","startupScript":"",
	"tasks":[{"id":"TASK_1","name":"任务1","description":"","max_retries":"0","script":""},
		{"id":"TASK_2","name":"任务2","description":"","max_retries":"0","script":""}],
		"deleteTaskIds":[]}]}
api	
*/
func UpdateFlow(c *Context) (code int, message string, data interface{}) {
	//1.解析参数
	newFlow, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return 50001, err.Error(), nil
	}

	c.Info("post: --> ", string(newFlow), " <--data")
	var paramQuest RpcQuestArgs
	err = json.Unmarshal([]byte(newFlow), &paramQuest)
	if err != nil { 
		return -1, err.Error(), "Do Unmarshal err"
	} 
	// fmt.Printf ( "%+v" , paramQuest) 
	args := paramQuest.Params[0]
	flow := &model.Flow{Id: args.Id, Name: args.Name, Description: args.Description, StartupScript: args.StartupScript}

	tasks := make([]*model.Task, 0, len(args.Tasks))
	for i, taskdata := range args.Tasks {
		task := &model.Task{}
		task.Id = taskdata.Id
		task.Name = taskdata.Name
		task.Description = taskdata.Description
		task.OrderId = i
		task.FlowId = flow.Id
		task.MaxRetries, _ = conv.Int(taskdata.MaxRetries)
		script := taskdata.Script
		task.Script = &script
		tasks = append(tasks, task)
	}

	err = model.UpdateFlow(flow, tasks, args.DeleteTaskIds)
	if err != nil {
		log.Error(err)
		return -1, err.Error(), "err"
	}
	return 0, "ok", nil
}
/*
{"id":2,"method":"FlowService.AddFlow",
"params":[{"id":"NEW_FLOW3","name":"新流程","creator":"helight","description":null,"startupScript":"",
"tasks":[{"id":"TASK_1","name":"任务1","description":"","max_retries":"0","script":""},
{"id":"TASK_2","name":"任务2","description":"","max_retries":"0","script":""},
{"id":"TASK_3","name":"任务3","description":"","max_retries":"0","script":""}],
"deleteTaskIds":[]}]}
*/
// func (h *FlowService) AddFlow(r *http.Request, args *FlowDataArgs, reply *AddFlowReply) error {
func AddFlow(c *Context) (code int, message string, data interface{}) {
	//1.解析参数
	newFlow, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return 50001, err.Error(), nil
	}
	
	c.Info("post: --> ", string(newFlow), " <--data")
	var paramQuest RpcQuestArgs
	err = json.Unmarshal([]byte(newFlow), &paramQuest)
	if err != nil { 
        return -1, err.Error(), "Do Unmarshal err"
	} 
	// fmt.Printf ( "%+v" , paramQuest) 
	args := paramQuest.Params[0]
	// return -1, string(newFlow), args
	flow := &model.Flow{Id: args.Id, Name: args.Name, Description: args.Description, StartupScript: args.StartupScript, Creator: args.Creator}
	
	tasks := make([]*model.Task, 0, len(args.Tasks))
	for i, taskdata := range args.Tasks {
		task := &model.Task{}
		task.Id = taskdata.Id
		task.Name = taskdata.Name
		task.Description = taskdata.Description
		task.OrderId = i
		task.FlowId = flow.Id
		task.MaxRetries, _ = conv.Int(taskdata.MaxRetries)
		script := taskdata.Script
		task.Script = &script
		tasks = append(tasks, task)
	}
	err = model.AddFlow(flow, tasks)
	if err != nil {
		if dberr, match := err.(*mysql.MySQLError); match {
			if dberr.Number == 1062 {
				return -1, "Flow Id " + args.Id + " is exists.", errors.New("Flow Id " + args.Id + " is exists.")
			}
		}
		log.Error(err)
		return 50001, err.Error(), nil
	}
	return 0, "ok", nil
}

// begin_date=2019-03-29&end_date=2019-03-31&set_nu=all&process_type=all&process_state_type=-1&pid=all
func GetFlows(c *Context) (code int, message string, data interface{}) {

	begin_date := c.Query("begin_date")
	end_date := c.Query("end_date")
	set_nu := c.Query("set_nu")
	flow_id := c.Query("process_type")
	process_state_type := c.Query("process_state_type")
	pid := c.Query("pid")

	cond1 := " (fi.running_day between '"+begin_date+" 00:00:00' and '"+end_date+" 23:59:59')"
	cond2 := " and ('"+set_nu+"' = 'all') and ('"+flow_id+"' = 'all' or f.id = '"+flow_id+"')"
	cond3 := " and ('"+process_state_type+"' = 'all' or fi.state = "+process_state_type+" or ("+process_state_type+" = -1 and fi.`state` in (0,1,3)))"
	cond4 := " and ('"+pid+"' = 'all' or fi.pid = '"+pid+"') ORDER BY fi.running_day desc, fi.flow_id, last_update_time DESC "

	cond := cond1 + cond2 + cond3 +cond4
	flows, err := model.GetFlowInstAll(cond)
	if err != nil {
		return -1, "get db error", ""
	}
	return -1, "ok", flows
}