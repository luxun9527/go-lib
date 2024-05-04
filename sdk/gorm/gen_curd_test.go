package main

import (
	"context"
	"go-lib/sdk/gorm/gen/dao/query"
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
		log.Panicf("err:%v", err)
	}
	var u []*User
	if db.Find(&u).Error != nil {
		log.Panicf("err:%v", err)
	}

	//gen join
	dao.User.WithContext(ctx).Join(company, user.CompanyID.EqCol(company.ID)).Find()
	//gorm join
	db.Joins("LEFT JOIN Company c ON c.id =  u.company_id").Find(&u)

	// gen分页
	dao.User.WithContext(ctx).FindByPage(0, 10)
	//gorm 分页
	db.Offset(0).Limit(10).Find(&u)
}
