# Intor Oneflow 介绍
oneflow，目标是web化的分布式流程系统，可以远程执行任务

# Install dir & work dir:

    you can configure the work dir for oneflow in config/config.go like this:
    ServerRoot = flag.String("server_root`", "/data/oneflow/", "Root of flow server.")
这个路径必须配置正确，否则无法启动程序

    for example the work dir is /data/oneflow/, so you should copy pkg dir which is in this project
to the work dir, or just make a link from pkg dir to the work dir. just like this:
``` sh
ln -s /home/helight/oneflow/pkg /data/oneflow           // 是工作目录，要和上面的配置一致。
ln -s /home/helight/oneflow/web /home/helight/www/      // 这是web页面
```

# Config
## workdir
    config/config.go
``` go
ServerRoot = flag.String("server_root`", "/data/oneflow/", "Root of flow server.")
```
## database:
数据库配置，非常简单。
    config/database.go
``` go
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
```
## server work port
服务端口配置
    config/server.go
``` go
ServerHost = flag.String("server_host", helper.GetIPAddr(), "Host of flow server.")
ServerPort = flag.String("server_port", "3001", "Port of flow server.")
```

# 编译 golang服务
目前项目是用dep来管理第三方库，所以要先下载第三方库，在进入src目录之后，使用`dep ensure`，执行之后再执行`make`即可。
```
cd src
dep ensure
make
```
## nsq config & init
使用nsq作为消息队列，所以需要配置nsq。
### config
    config/nsq.go
    just configure the nsqd address and nsqlookupd

### nsq install & init
#### main server
启动nsq主服务
``` sh
./nsqlookupd
./nsqd --lookupd-tcp-address=127.0.0.1:4160
```
#### web service for nsq
启动nsq的web管理端
``` sh
./nsqadmin --lookupd-http-address=127.0.0.1:4161 
```
http://localhost:4171/ 通过这个地址可以直接访问nsq的管理端
#### create on topic
创建一个topic
``` sh
curl -d "hello 1" "http://127.0.0.1:4151/pub?topic=status.task.flow.data"
```
# db init
数据库名称是：aflow
``` sql
create database aflow;
create user 'admin'@'localhost' IDENTIFIED BY 'mysql';
grant all privileges on aflow.* to admin@localhost identified by 'mysql';
```
## database create script
数据库初始化脚本
``` sql
mysql -uroot -Daflow < /web/oneflow.sql
```
## init data
初始化数据
``` sql
insert into tbProducts (`PId`,`State`,`Name`,`DBHost`,`DBName`, `StarLevel`) 
    values ("222", 1, "223","127.0.0.1","223db", 1);
insert into tbServer(`host`,`port`,`username`,`password`) values ("127.0.0.1", 22, "helight", "helight");
insert into tbServer(`host`,`port`,`username`,`password`) values ("10.0.2.15", 22, "helight", "helight");
insert into tbServer(`host`,`port`,`username`,`password`) values ("172.22.112.56", 22, "helight", "helight");
```
    ERROR:ssh: handshake failed: ssh: unable to authenticate, attempted methods [none password], no supported methods remain

# lua demo
整个系统使用lua语言来写执行脚本，下面是一个简单的例子
    put those lua script to the task and run it, first you should have the shell script in this lua script.    

``` lua 
print("123")
puuid, output=gassert(remote_exec("127.0.0.1", "/data/oneflow/bin/script/test.sh", "222"))
print("Process UUID:"..puuid)
print("Remote exec output:"..output)
```
