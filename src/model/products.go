package model

import (
	"config"

    log "github.com/cihub/seelog"
)

type Products struct {
	PId		string `sql:"PId"`
	State		int    `sql:"state"`
	Name		string `sql:"Name"`
	DBHost		string `sql:"DBHost"`
	DBName		string `sql:"DBName"`
	StarLevel	int    `sql:"StarLevel"`
}

func FindProductsDB(condition string, args ...interface{}) ([]*Products, error) {
	sqlStr := "SELECT `PId`, `State`, `Name`, `DBHost`, `DBName`, `StarLevel` FROM `tbProducts`"
	if len(condition) > 0 {
		sqlStr = sqlStr + " WHERE " + condition
	}
	results := make([]*Products, 0)

    log.Debug("sql: ", sqlStr, " args: ", args)
	stmt, err := config.GetDBConnect().Prepare(sqlStr)
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
			model := Products{}
			values := []interface{}{
				new(interface{}),
				new(interface{}),
				new(interface{}),
				new(interface{}),
				new(interface{}),
				new(interface{}),
			}
			rows.Scan(values...)
			model.PId = (string)((*(values[0].(*interface{}))).([]uint8))

			if *(values[1].(*interface{})) == nil {
				model.State = 0
			} else {
				t_State := (int)((*(values[1].(*interface{}))).(int64))
				model.State = t_State
			}
			if *(values[2].(*interface{})) == nil {
				model.Name = ""
			} else {
				t_Name := (string)((*(values[2].(*interface{}))).([]uint8))
				model.Name = t_Name
			}
			if *(values[3].(*interface{})) == nil {
				model.DBHost = ""
			} else {
				t_DBHost := (string)((*(values[3].(*interface{}))).([]uint8))
				model.DBHost = t_DBHost
			}
			if *(values[4].(*interface{})) == nil {
				model.DBName = ""
			} else {
				t_DBName := (string)((*(values[4].(*interface{}))).([]uint8))
				model.DBName = t_DBName
			}
			model.StarLevel = (int)((*(values[5].(*interface{}))).(int64))

			results = append(results, &model)
		}
	}
	return results, nil
}

func GetProducts(condition string, args ...interface{}) (*Products, error) {
	results, err := FindProductsDB(condition, args...)

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

func GetProductsByKey(pid string) (*Products, error) {
	return GetProducts("`PId`=?", pid)
}
