package main

import (
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/rawsql"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "E:\\demoproject\\go-lib\\sdk\\gorm\\gen\\gen_sql_file\\query",
	})

	// https://github.com/go-gorm/rawsql/blob/master/tests/gen_test.go
	gormdb, _ := gorm.Open(rawsql.New(rawsql.Config{
		//SQL:      rawsql,                      // 建表sql
		FilePath: []string{
			//"./sql/user.sql", // 建表sql文件
			"E:\\demoproject\\go-lib\\sdk\\gorm\\gen\\gen_sql_file\\test.sql", // 建表sql目录
		},
	}))
	g.UseDB(gormdb) // 重新引用你的 gorm db
	// 按照约定为结构 `model.User` 生成基本类型安全的DAO API
	g.ApplyBasic(
		// 基于 `user` 表生成 `User` 结构
		g.GenerateModel("account"),
	)

	g.ApplyBasic(
		// 从当前数据库生成所有表结构
		g.GenerateAllTable()...,
	)

	// 生成代码
	g.Execute()
}
