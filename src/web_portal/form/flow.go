package form

import (
	"model"
)

type Flow struct {
	model.Flow
	Tasks []*model.Task
}
