package dao

import (
	"github.com/gogf/gf/v2/frame/g"
)

// TaskQueue 是任务队列的DAO操作对象
var TaskQueue = g.DB().Model("task_queue")
