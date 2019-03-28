export GOPATH=$(shell pwd)/../:$(shell pwd)/../go-package/
export CC=gcc -no-pie

#@echo "gopath:"$(GOPATH)
#GOBIN=/usr/local/GOBIN/bin/go
GOBIN=go

FServer=flow_server
FAgent=task_agent/flow_agent

# find . -type d -exec golint {} \;

all: server agent

#server:export GOROOT=/usr/local/GOBIN/
server: export GOPATH=$(shell pwd)/../:$(shell pwd)/../go-package/
server:
	$(GOBIN) build -o $(FServer) flow_server.go
	cp $(FServer) ../pkg/bin/

#agent:export GOROOT=/usr/local/GOBIN/
agent:
	$(GOBIN) build -o app/client app/client.go
	$(GOBIN) build -o app/luar_demo app/luar_demo.go
	$(GOBIN) build -o app/microtest app/microtest.go
	$(GOBIN) build -o $(FAgent) task_agent/main.go
	cp $(FAgent) ../pkg/bin/

#test:export GOROOT=/usr/local/GOBIN/
test:
#	$(GOBIN) test -v lua_helper/lua_test.go lua_helper/lua_helper.go lua_helper/remote.go
	$(GOBIN) test -v scheduler/parser.go  scheduler/parser_test.go scheduler/spec.go


clean:
	go clean
	rm $(FServer) $(FAgent) app/client app/luar_demo app/microtest -f

