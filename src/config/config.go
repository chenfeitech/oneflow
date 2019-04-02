package config

import (
	"flag"
	"utils/helper"
)

var (
	Online     = flag.Bool("online", Bool(true, false), "online flag")
	ServerRoot = flag.String("server_root`", "/mnt/oneflow/pkg/", "Root of flow server.")
	ServerHost = flag.String("server_host", helper.GetIPAddr(), "Host of flow server.")
	ServerPort = flag.String("server_port", "3001", "Port of flow server.")
)
