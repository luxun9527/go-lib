package json

import (
	"fmt"
	"github.com/tidwall/gjson"
	"testing"
)

func TestGetJson(t *testing.T) {
	json := `{"name":{"first":"li","last":"dj"},"age":18}`
	lastName := gjson.Get(json, "name.last")
	fmt.Println("last name:", lastName.String())

	age := gjson.Get(json, "age")
	fmt.Println("age:", age.Int())
}
