package main

import (
	"context"
	"go-lib/sdk/gorm/gen/dao/query"
	"gorm.io/gen/field"
	"log"
	"testing"
)

func TestGen(t *testing.T) {
	c := query.Use(db).Card
	log.Printf("%p", &c)

	c.WithContext(context.Background()).Select(field.NewField("", "test")).Where(field.NewField("", "test")).Find()

}
