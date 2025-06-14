package snowflake

import (
	"github.com/gookit/goutil/strutil"
	"github.com/yitter/idgenerator-go/idgen"
	"testing"
)

func TestSnowflake(t *testing.T) {
	var options = idgen.NewIdGeneratorOptions(1)
	//options.WorkerIdBitLength = 10
	idgen.SetIdGenerator(options)

	var newId = idgen.NextId()
	t.Log(newId)

	md5 := strutil.Md5("12345678")
	t.Log(md5)
}
