package config

import (
	"database/sql"
	"flag"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	log "github.com/cihub/seelog"
	_ "github.com/go-sql-driver/mysql"
)

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

var (
	db           *sql.DB
	db_mutex sync.Mutex
	dbs      = map[string]*sql.DB{}
)

func GetDBConnect() *sql.DB {
	if db != nil {
		return db
	}

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local&charset=utf8&timeout=15s", *db_username, *db_password, *db_host, *db_port, *db_name)
	var err error
	db, err = sql.Open("mysql", connStr)
	if err != nil {
		log.Error("Connect to database failed:", connStr, err)
		panic(err)
	}
	db.SetMaxIdleConns(10)
	return db
}

func GetFlowDBConnect(host string, dbname string) (*sql.DB, error) {
	key := host + "_" + dbname
	db_mutex.Lock()
	if db, ok := dbs[key]; ok {
		db_mutex.Unlock()
		return db, nil
	}
	db_mutex.Unlock()

	username := "admin"
	password := "mysql"

	is_domain := strings.Contains(host, ":")
	if is_domain {
		password = "mysql"
	} else {
		host = host + ":3306"
	}

	connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&loc=Local&charset=utf8&timeout=15s", username, password, host, dbname)

	log.Error("Connect to database :", connStr)
	var err error
	var db *sql.DB
	for i := 0; i < 5; i++ {
		db, err = sql.Open("mysql", connStr)
		if err != nil {
			log.Error("Connect to database failed:", connStr, err)
			time.Sleep(3 * time.Second)
			continue
		}
		db.SetMaxIdleConns(10)
		db_mutex.Lock()
		dbs[key] = db
		db_mutex.Unlock()
		return db, nil
	}

	return db, err
}
