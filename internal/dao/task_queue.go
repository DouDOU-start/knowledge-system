// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package dao

import (
	"knowledge-system-api/internal/dao/internal"

	"github.com/gogf/gf/v2/database/gdb"
)

// taskQueueDao is the manager for logic model data accessing and custom defined data operations functions management.
type taskQueueDao struct {
	*internal.TaskQueueDao
}

var (
	// TaskQueue is globally public accessible object for table task_queue operations.
	TaskQueue = taskQueueDao{
		internal.NewTaskQueueDao(),
	}
)

// FillCustomDao adds custom DAO functions for the DAO.
func (d *taskQueueDao) FillCustomDao(modelHandler ...gdb.ModelHandler) {
	d.TaskQueueDao = internal.NewTaskQueueDao(modelHandler...)
}
