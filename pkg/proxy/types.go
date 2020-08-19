package proxy

// SqlCommand sql命令
type SqlCommand struct {
	// App 应用名称
	App string `json:"app"`
	// Env  环境, dev,test,prod
	Env      string `json:"env"`
	Database string `json:"database"`
	SQL      string `json:"sql"`
	// User db用户名
	User string `json:"user"`

	// Duration sql执行间隔(秒),精确到毫秒
	Duration int64 `json:"duration"`

	// 接收到sql的时间戳表示
	Occtime int64 `json:"occtime"`
}
