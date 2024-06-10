package main

import (
	"context"
	"go-lib/sdk/gorm/gen/dao/model"
	"go-lib/sdk/gorm/gen/dao/query"
	"gorm.io/gorm/clause"
	"log"
	"testing"
)

func TestGenSelect(t *testing.T) {
	dao := query.Use(db)
	//card := dao.Card
	user := dao.User
	company := dao.Company
	ctx := context.Background()

	_, err := user.WithContext(ctx).Find()
	if err != nil {
		log.Printf("err:%v", err)
	}
	var u []*User
	if db.Where("id=1").Find(&u).Error != nil {
		log.Printf("err:%v", err)
	}

	//gen join
	dao.User.WithContext(ctx).Join(company, user.CompanyID.EqCol(company.ID)).Find()
	//gorm join
	db.Joins("LEFT JOIN Company c ON c.id =  u.company_id").Find(&u)

	// gen分页
	dao.User.WithContext(ctx).FindByPage(0, 10)
	//gorm 分页
	db.Offset(0).Limit(10).Find(&u)

	// gen子查询
	subQuery := company.WithContext(ctx).
		Select(company.ID)
	//查登录信息
	_, err = user.WithContext(ctx).
		Where(user.Columns(user.CompanyID).In(subQuery)).
		Find()
	dao.User.WithContext(ctx).Clauses(clause.OnConflict{
		//Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"role": "user"}),
	}).Create(&model.User{})
}
