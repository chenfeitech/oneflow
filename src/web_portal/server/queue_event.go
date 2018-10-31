package server

import (
	"bytes"
	"config"
	"web_portal/jsonrpc2"
	"middleware/jsonrpc_client"
	"net/http"
	"strconv"
	"strings"

	json "web_portal/jsonrpc2/json2"
	nsq "github.com/EncoreJiang/go-nsq"
	"github.com/cihub/seelog"
)

func startEventLoop(concurrent int, service *FlowService) error {
	jsonRPCService := jsonrpc2.NewServer()
	jsonCodec := json.NewCodec()
	jsonRPCService.RegisterCodec(jsonCodec, "application/json")
	jsonRPCService.RegisterCodec(jsonCodec, "application/json; charset=UTF-8")
	jsonRPCService.RegisterService(service, "")

	nsqConfig := nsq.NewConfig()
	nsqConfig.MaxInFlight = 10
	q, err := nsq.NewConsumer(config.NSQFlowTaskStatusTopic, "executor", nsqConfig)
	if err != nil {
		return err
	}
	q.AddConcurrentHandlers(&eventExecutor{jsonRPCService}, concurrent)
	if len(config.NSQLookupdAddress) != 0 {
		err = q.ConnectToNSQLookupds(strings.Split(config.NSQLookupdAddress, ";"))
	} else {
		err = q.ConnectToNSQD(config.NSQDAddress)
	}
	if err != nil {
		return err
	}
	return nil
}

type eventExecutor struct {
	jsonRPCService *jsonrpc2.Server
}

func (te *eventExecutor) HandleMessage(message *nsq.Message) error {
	request := &http.Request{}
	request.Header = http.Header{}
	request.Header.Add("X-Timestamp", strconv.FormatInt(message.Timestamp, 10))
	resp := te.jsonRPCService.ServeBytes(request, message.Body)
	var reply interface{}
	err := jsonrpc_client.DecodeClientResponse(bytes.NewReader(resp), &reply)
	if err != nil {
		seelog.Errorf("Handle event %s failed: %v", message.Body, err)
		return err
	}
	return nil
}

func (te *eventExecutor) LogFailedMessage(message *nsq.Message) {
	seelog.Errorf("Handle event failed: %s", message.Body)
}
