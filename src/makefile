export GOPATH=$(shell pwd)/../:$(shell pwd)/../go-package/
export CC=gcc -no-pie

#go1.11=/usr/local/go1.11/bin/go
go1.11=go

FServer=flow_server
FAgent=task_agent/flow_agent

# find . -type d -exec golint {} \;

all: server agent

#server:export GOROOT=/usr/local/go1.11/
server:
	$(go1.11) build -o $(FServer) flow_server.go
	cp $(FServer) ../pkg/bin/

#agent:export GOROOT=/usr/local/go1.11/
agent:
	$(go1.11) build -o app/client app/client.go
	$(go1.11) build -o app/luar_demo app/luar_demo.go
	$(go1.11) build -o app/microtest app/microtest.go
	$(go1.11) build -o $(FAgent) task_agent/main.go
	cp $(FAgent) ../pkg/bin/

#test:export GOROOT=/usr/local/go1.11/
test:
#	$(go1.11) test -v lua_helper/lua_test.go lua_helper/lua_helper.go lua_helper/remote.go
	$(go1.11) test -v scheduler/parser.go  scheduler/parser_test.go scheduler/spec.go


clean:
	go clean
	rm $(FServer) $(FAgent) app/client app/luar_demo app/microtest -f

