package parse

type JoinItem struct {
	Table string
	Field string
}

type SelectStmt struct {
	App          string
	FromTables   []string
	SelectFields []string
	Joins        []JoinItem //a.id = b.id
}

type UpdateStmt struct {
	App   string
	Table string
}

type InsertStmt struct {
	App   string
	Table string
}

type DeleteStmt struct {
	App   string
	Table string
}
