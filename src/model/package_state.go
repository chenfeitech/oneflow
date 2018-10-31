package model

import (
	"config"
	"time"
)

type PackageState struct {
	Sequence            int        `sql:"Sequence"`
	PackageId           *int       `sql:"PackageId"`
	JobId               *int       `sql:"JobId"`
	State               *int       `sql:"State"`
	PartialNum          *int       `sql:"PartialNum"`
	EmptyPartialNum     int        `sql:"EmptyPartialNum"`
	Count               *int       `sql:"Count"`
	ReducedCount        *int       `sql:"Reduced_count"`
	Reduced             *int       `sql:"Reduced"`
	FilePath            *string    `sql:"FilePath"`
	Size                *int       `sql:"Size"`
	SizeAfterCompressed *int       `sql:"SizeAfterCompressed"`
	Compressed          *int       `sql:"Compressed"`
	CreateTime          time.Time  `sql:"CreateTime"`
	UpdateTime          *time.Time `sql:"UpdateTime"`
}

func FindPackageState(condition string, args ...interface{}) ([]*PackageState, error) {
	sqlStr := "SELECT `Sequence`, `PackageId`, `JobId`, `State`, `PartialNum`, `EmptyPartialNum`, `Count`, `Reduced_count`, `Reduced`, `FilePath`, `Size`, `SizeAfterCompressed`, `Compressed`, `CreateTime`, `UpdateTime` FROM `tbPackageState`"
	if len(condition) > 0 {
		sqlStr = sqlStr + " WHERE " + condition
	}
	results := make([]*PackageState, 0)

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
			model := PackageState{}
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
				new(interface{}),
				new(interface{}),
				new(interface{}),
				new(interface{}),
				new(interface{}),
				new(interface{}),
			}
			rows.Scan(values...)
			model.Sequence = (int)((*(values[0].(*interface{}))).(int64))
			if *(values[1].(*interface{})) == nil {
				model.PackageId = nil
			} else {
				t_PackageId := (int)((*(values[1].(*interface{}))).(int64))
				model.PackageId = &t_PackageId
			}
			if *(values[2].(*interface{})) == nil {
				model.JobId = nil
			} else {
				t_JobId := (int)((*(values[2].(*interface{}))).(int64))
				model.JobId = &t_JobId
			}
			if *(values[3].(*interface{})) == nil {
				model.State = nil
			} else {
				t_State := (int)((*(values[3].(*interface{}))).(int64))
				model.State = &t_State
			}
			if *(values[4].(*interface{})) == nil {
				model.PartialNum = nil
			} else {
				t_PartialNum := (int)((*(values[4].(*interface{}))).(int64))
				model.PartialNum = &t_PartialNum
			}
			model.EmptyPartialNum = (int)((*(values[5].(*interface{}))).(int64))
			if *(values[6].(*interface{})) == nil {
				model.Count = nil
			} else {
				t_Count := (int)((*(values[6].(*interface{}))).(int64))
				model.Count = &t_Count
			}
			if *(values[7].(*interface{})) == nil {
				model.ReducedCount = nil
			} else {
				t_ReducedCount := (int)((*(values[7].(*interface{}))).(int64))
				model.ReducedCount = &t_ReducedCount
			}
			if *(values[8].(*interface{})) == nil {
				model.Reduced = nil
			} else {
				t_Reduced := (int)((*(values[8].(*interface{}))).(int64))
				model.Reduced = &t_Reduced
			}
			if *(values[9].(*interface{})) == nil {
				model.FilePath = nil
			} else {
				t_FilePath := (string)((*(values[9].(*interface{}))).([]uint8))
				model.FilePath = &t_FilePath
			}
			if *(values[10].(*interface{})) == nil {
				model.Size = nil
			} else {
				t_Size := (int)((*(values[10].(*interface{}))).(int64))
				model.Size = &t_Size
			}
			if *(values[11].(*interface{})) == nil {
				model.SizeAfterCompressed = nil
			} else {
				t_SizeAfterCompressed := (int)((*(values[11].(*interface{}))).(int64))
				model.SizeAfterCompressed = &t_SizeAfterCompressed
			}
			if *(values[12].(*interface{})) == nil {
				model.Compressed = nil
			} else {
				t_Compressed := (int)((*(values[12].(*interface{}))).(int64))
				model.Compressed = &t_Compressed
			}
			model.CreateTime = (*(values[13].(*interface{}))).(time.Time)
			if *(values[14].(*interface{})) == nil {
				model.UpdateTime = nil
			} else {
				t_UpdateTime := (*(values[14].(*interface{}))).(time.Time)
				model.UpdateTime = &t_UpdateTime
			}

			results = append(results, &model)
		}
	}
	return results, nil
}

func GetPackageState(condition string, args ...interface{}) (*PackageState, error) {
	results, err := FindPackageState(condition, args...)

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

func GetPackageStateByPackageId(package_id int) (*PackageState, error) {
	return GetPackageState("PackageId=?", package_id)
}
