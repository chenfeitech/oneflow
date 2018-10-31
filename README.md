# Intor Oneaction
oneaction

# Install dir & work dir:

    you can configure the work dir for oneaction in config/config.go like this:
    ServerRoot = flag.String("server_root`", "/data/oneflow/", "Root of flow server.")

    for example the work dir is /data/oneflow/, so you should copy pkg dir which is in this project
to the work dir, or just make a link from pkg dir to the work dir. just like this:
    ln -s /home/helight/oneaction/pkg /data/oneflow

# Config
## workdir
    config/config.go
    ServerRoot = flag.String("server_root`", "/data/oneflow/", "Root of flow server.")
## database:
    config/database.go
    var (
        db_host = flag.String("db_host", func() string {
                              if runtime.GOOS == "darwin" {
                              return "127.0.0.1"
                              } else {
                              return "127.0.0.1"
                              }
                              }(), "mysql server host")
        db_port     = flag.Int("db_port", 3306, "mysql server port")
        db_username = flag.String("db_username", "admin", "mysql server username")
        db_password = flag.String("db_password", "mysql", "mysql server password")
        db_name     = flag.String("db_name", "aflow", "mysql server name")
    )

## server work port
    config/server.go
    ServerHost = flag.String("server_host", helper.GetIPAddr(), "Host of flow server.")
    ServerPort = flag.String("server_port", "3001", "Port of flow server.")
## nsq config & init

### config
    config/nsq.go
    just configure the nsqd address and nsqlookupd

### nsq install & init
#### main server
    ./nsqlookupd
    ./nsqd --lookupd-tcp-address=127.0.0.1:4160

#### web service for nsq
    ./nsqadmin --lookupd-http-address=127.0.0.1:4161
    http://localhost:4171/

#### create on topic
    curl -d "hello 1" "http://127.0.0.1:4151/pub?topic=status.task.flow.data"

# db init
## database create script
    /web/actionflow.sql

## init data
    insert into tbProducts (`PId`,`State`,`Name`,`DBHost`,`DBName`, `StarLevel`) 
        values ("222", 1, "223","127.0.0.1","223db", 1);
    insert into tbServer(`host`,`port`,`username`,`password`) values ("127.0.0.1", 22, "helight", "helight");
    insert into tbServer(`host`,`port`,`username`,`password`) values ("10.0.2.15", 22, "helight", "helight");

# lua demo
    put those lua script to the task and run it, first you should have the shell script in this lua script.
    
    '''
    print("123")
    puuid, output=gassert(remote_exec("127.0.0.1", "/data/actionflow/bin/script/test.sh", "222"))
    print("Process UUID:"..puuid)
    print("Remote exec output:"..output)
    '''
