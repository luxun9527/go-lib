package snowflake

import (
	"github.com/yitter/idgenerator-go/idgen"
	"testing"
)

func TestSnowflake(t *testing.T) {
	var options = idgen.NewIdGeneratorOptions(1)
	//options.WorkerIdBitLength = 10
	idgen.SetIdGenerator(options)

	var newId = idgen.NextId()
	t.Log(newId)
}
