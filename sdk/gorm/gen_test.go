package main

import (
	"context"
	"go-lib/sdk/gorm/gen/dao/query"
	"log"
	"testing"
)

func TestGen(t *testing.T) {
	c := query.Use(db).Card
	log.Printf("%p", &c)
	query.Use(db).Card.WithContext(context.Background()).RawWhere("IF(a=?,t,t1)", 1).Take()
	c1 := query.Use(db).Card
	log.Printf("%p", &c1)
}
