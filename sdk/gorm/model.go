package main

// User1Columns get sql column name.获取数据库列名
var User1Columns = struct {
	ID        string
	Name      string
	Age       string
	CreatedAt string
	UpdatedAt string
}{
	ID:        "id",
	Name:      "name",
	Age:       "age",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}
