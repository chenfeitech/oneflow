package config

import (
	"flag"
)

var (
	Online = flag.Bool("online", Bool(true, false), "online flag")
	ServerRoot = flag.String("server_root`", "/data/oneflow/", "Root of flow server.")
)
