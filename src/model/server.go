package model

import (
	"config"
	"time"
)

var _ = time.Now
var funcGetConnect = config.GetDBConnect

type Server struct {
	Host           string  `sql:host`
	Port           int     `sql:port`
	Username       *string `sql:username`
	Password       *string `sql:password`
	CryptoPassword string  `sql:crypto_password`
	Supervisors    *string `sql:supervisors`
	ContFailures   int     `sql:cont_failures`
	FailureUUID    *string `sql:failure_uuid`
	Tags           string  `sql:tags`
}

func AddServer(model *Server) (int64, error) {
	sqlStr := "INSERT INTO `tbServer` (`host`, `port`, `username`, `password`, `crypto_password`, `supervisors`, `cont_failures`, `failure_uuid`, `tags`) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := funcGetConnect().Exec(sqlStr, model.Host, model.Port, model.Username, model.Password, model.CryptoPassword, model.Supervisors, model.ContFailures, model.FailureUUID, model.Tags)
	if err != nil {
		return 0, err
	} else {
		return result.LastInsertId()
	}
}

func FindServer(condition string, args ...interface{}) ([]*Server, error) {
	sqlStr := "SELECT `host`, `port`, `username`, `password`, `crypto_password`, `supervisors`, `cont_failures`, `failure_uuid`, `tags` FROM `tbServer`"
	if len(condition) > 0 {
		sqlStr = sqlStr + " WHERE " + condition
	}
	sqlStr = sqlStr + " LIMIT 0,1000"
	results := make([]*Server, 0)

	stmt, err := funcGetConnect().Prepare(sqlStr)
	if err != nil {
		return results, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		return results, err
	} else {
		defer rows.Close()
		for rows.Next() {
			model := Server{}
			values := []interface{}{
				new(interface{}),
				new(interface{}),
				new(interface{}),
				new(interface{}),
				new(interface{}),
				new(interface{}),
				new(interface{}),
				new(interface{}),
				new(interface{}),
			}
			rows.Scan(values...)

			if *(values[0].(*interface{})) != nil {
				tmp := string((*(values[0].(*interface{}))).([]uint8))
				model.Host = tmp
			}

			if *(values[1].(*interface{})) != nil {
				tmp := int((*(values[1].(*interface{}))).(int64))
				model.Port = tmp
			}

			if *(values[2].(*interface{})) != nil {
				tmp := string((*(values[2].(*interface{}))).([]uint8))
				model.Username = &tmp
			}

			if *(values[3].(*interface{})) != nil {
				tmp := string((*(values[3].(*interface{}))).([]uint8))
				model.Password = &tmp
			}

			if *(values[4].(*interface{})) != nil {
				tmp := string((*(values[4].(*interface{}))).([]uint8))
				model.CryptoPassword = tmp
			}

			if *(values[5].(*interface{})) != nil {
				tmp := string((*(values[5].(*interface{}))).([]uint8))
				model.Supervisors = &tmp
			}

			if *(values[6].(*interface{})) != nil {
				tmp := int((*(values[6].(*interface{}))).(int64))
				model.ContFailures = tmp
			}

			if *(values[7].(*interface{})) != nil {
				tmp := string((*(values[7].(*interface{}))).([]uint8))
				model.FailureUUID = &tmp
			}

			if *(values[8].(*interface{})) != nil {
				tmp := string((*(values[8].(*interface{}))).([]uint8))
				model.Tags = tmp
			}
			results = append(results, &model)
		}
	}
	return results, nil
}

func GetServer(condition string, args ...interface{}) (*Server, error) {
	results, err := FindServer(condition, args...)

	if err != nil {
		return nil, err
	} else {
		if len(results) > 0 {
			return results[0], nil
		} else {
			return nil, nil
		}
	}
}

func GetServerByHost(host string) (*Server, error) {
	return GetServer("`host`=?", host)
}
