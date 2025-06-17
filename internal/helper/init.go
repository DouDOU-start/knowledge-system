package helper

// InitFunctions 初始化函数
type InitFunctions struct {
	// InitServices 初始化服务层
	InitServices func()
}

var (
	// 全局初始化函数
	globalInitFunctions InitFunctions
)

// SetInitFunctions 设置初始化函数
func SetInitFunctions(fns InitFunctions) {
	globalInitFunctions = fns
}

// InitAll 执行所有初始化
func InitAll() {
	if globalInitFunctions.InitServices != nil {
		globalInitFunctions.InitServices()
	}
}
