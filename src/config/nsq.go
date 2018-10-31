package config

import (
	"flag"
	"runtime"
)

var (
	NSQDAddress            string
	NSQDServiceName        string
	NSQLookupdAddress      string
	NSQFlowTaskStatusTopic = "status.task.flow.idata"
)

func init() {

	if runtime.GOOS == "darwin" {
		flag.StringVar(&NSQDAddress, "nsqd_address", "127.0.0.1:4150", "NSQ Address")
	} else {
		flag.StringVar(&NSQDAddress, "nsqd_address", "127.0.0.1:4150;127.0.0.1:4150", "NSQD Address")
		// flag.StringVar(&NSQDServiceName, "nsqd_service_name", "nsqd", "NSQD Service Name")
		flag.StringVar(&NSQLookupdAddress, "nsq_lookupd_address", "127.0.0.1:4161;127.0.0.1:4161", "NSQ Lookupd Address")
	}
}
