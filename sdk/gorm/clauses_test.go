package main

import (
	"go-lib/sdk/gorm/gen/dao/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"testing"
)

type ClauseUser struct {
	model.User
	C CustomClause
}

func TestClauses(t *testing.T) {
	user := &User{}
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&user)
	//INSERT INTO `user` (`username`,`age`,`fav`,`company_id`,`created_at`,`updated_at`,`deleted_at`) VALUES ('',0,'',0,1714919066,1714919066,0) ON DUPLICATE KEY UPDATE `id`=`id`
	db.Clauses(clause.Eq{
		Column: "id",
		Value:  1,
	}).Find(&user)
	//SELECT * FROM `user` WHERE `id` = 1 AND `user`.`deleted_at` = 0 AND `user`.`id` = 19
	db.Find(&ClauseUser{})
	//SELECT * FROM `user` WHERE `user`.`c` IS NULL
}

type CustomClause int32

func (CustomClause) QueryClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{CustomClauseQuery{Field: f}}
}

type CustomClauseQuery struct {
	Field *schema.Field
}

func (sd CustomClauseQuery) Name() string {
	return ""
}

func (sd CustomClauseQuery) Build(builder clause.Builder) {

}

func (sd CustomClauseQuery) MergeClause(c *clause.Clause) {

}

func (sd CustomClauseQuery) ModifyStatement(stmt *gorm.Statement) {
	stmt.AddClause(clause.Where{Exprs: []clause.Expression{
		clause.Eq{Column: clause.Column{Table: clause.CurrentTable, Name: sd.Field.DBName}, Value: nil},
	}})
}
