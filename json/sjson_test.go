package json

import (
	"fmt"
	"github.com/tidwall/sjson"
	"testing"
)

func TestSetJson(t *testing.T) {
	const SetJson = `{"name":{"first":"li","last":"dj"},"age":18}`

	value, _ := sjson.Set(SetJson, "name.last", "dajun")
	value, _ = sjson.Set(SetJson, "sex", "男")
	fmt.Println(value)
}
