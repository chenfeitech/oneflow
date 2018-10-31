# oneaction
oneaction

# install dir:

link your work dir to /data/oneflow

ln -s /home/helight/oneaction/pkg /data/oneflow

# config database

config/database.go

# nsq init

./nsqlookupd 
./nsqd --lookupd-tcp-address=127.0.0.1:4160

# web service for nsq 
./nsqadmin --lookupd-http-address=127.0.0.1:4161

http://localhost:4171/


## create on topic
curl -d "hello 1" "http://127.0.0.1:4151/pub?topic=status.task.flow.data"


# db demo


insert into tbProducts (`PId`,`State`,`Name`,`DBHost`,`DBName`, `StarLevel`) values("222", 1, "223","127.0.0.1","223db", 1);

insert into tbServer(`host`,`port`,`username`,`password`) values("127.0.0.1", 22, "helight", "helight");

insert into tbServer(`host`,`port`,`username`,`password`) values("10.0.2.15", 22, "helight", "helight");
