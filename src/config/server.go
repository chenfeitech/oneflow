package config

import (
	"flag"
	"utils/helper"
)

var (
	ServerHost = flag.String("server_host", helper.GetIPAddr(), "Host of flow server.")
	ServerPort = flag.String("server_port", "3001", "Port of flow server.")
)
