package query

import (
	"gorm.io/gen"
	"gorm.io/gorm/clause"
)

type StringBinaryLike struct {
	gen.DO
	TableName string
	Column  string
	Value  string
}
func (s StringBinaryLike)BeCond() interface{}{
	return s
}


func (like StringBinaryLike) Build(builder clause.Builder) {
	builder.WriteQuoted(like.TableName)
	builder.WriteByte(46)
	builder.WriteQuoted(like.Column)
	builder.WriteString(" LIKE binary ")
	builder.AddVar(builder, like.Value)
}


