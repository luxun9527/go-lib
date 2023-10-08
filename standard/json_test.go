package standard

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"log"

	"testing"
)

func TestJson(t *testing.T) {
	h := gin.H{"data": "test"}
	d,_ := json.Marshal(h)
	str := "["
	for _,v := range d {
		str +=cast.ToString(v)+","
	}
	log.Println(str+"]")
	
}
