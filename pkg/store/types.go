package store

// LogCommand sql日志存储内容定义
// 应该是经过sql解析后的内容
type LogCommand struct {
	// Env  环境, dev,test,prod
	Env string `json:"env"`
	// App 应用名称
	App      string `json:"app"`
	Database string `json:"database"`
	// SQL 原始sql
	SQL string `json:"sql"`
	// SqlFingerprint sql指纹,用来标识sql,排除sql中值的变化
	SqlFingerprint string `json:"sql_fingerprint"`
	// User db用户名
	User   string   `json:"user"`
	Tables []string `json:"tables"`
	Fields []string `json:"fields"`
	// PerformanceScore 性能分数
	PerformanceScore float32 `json:"performance_score"`
	// Duration sql执行间隔(秒),精确到毫秒
	Duration int64 `json:"duration"`
	// 接收到sql的时间戳表示
	Occtime int64 `json:"occtime"`
}
