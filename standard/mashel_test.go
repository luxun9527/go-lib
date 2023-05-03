package standard

import (
	"testing"
)

func TestUnmarshal(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
	}

}
