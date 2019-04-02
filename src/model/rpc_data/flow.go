package rpc_data


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

type StateDataArgs struct {
	PId     string
	Key     string
	FlowId  string
	TaskId  string
	State   int
	Creator string
	Date    string

	ExtraData map[string]string
}

type StateReply struct {
	Message string
}

type RerunTaskArgs struct {
	FlowInstId int
	TaskId     string
	SingleTask bool
	Creator    string
}

type RerunTaskReply struct {
}

type GetRemoteLogArgs struct {
	Host string
	Uuid string
	Date string
}

type GetRemoteLogReply struct {
	Cmdline string
	Output  string
	Error   string
}

type RunScriptArgs struct {
	Script string
}

type RunScriptReply struct {
	Output string
}

type StartFlowArgs struct {
	Id            string
	TaskId        string `json:"task_id"`
	PId           string `json:"pid"`
	Key           string
	Date          string
	Creator       string
	StartupScript string `json:"startup_script"`
}

type StartFlowReply struct {
	FlowInstId int `json:"flow_inst_id"`
}

type KillTaskInstanceArgs struct {
	FlowInstId int
	TaskId     string
}

type SetTaskInstanceSuccessArgs struct {
	FlowInstId int
	TaskId     string
}

type SetFlowInstanceSuccessArgs struct {
	FlowInstId int
}

type TaskAlarmArgs struct {
	FlowInstId int
	TaskId     string
	Content    string
}
