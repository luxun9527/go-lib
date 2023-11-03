package main

import (
	"bytes"
	"encoding/json"
	"log"
)

func main()  {
	f := NewFetcherHttps("goo.su")

	b := map[string]interface{}{"alias": "", "is_public": true, "group_id": 1, "url": "http://stg.kaiactivity.com"}
	d, _ := json.Marshal(b)
	f.Header.Set("x-goo-api-token","W4exIOZZym2nxYwr3puHz2orStoMjdM1QmQVXvZidskftEgOi4VLqfufRKv1")
	resp, body, err := f.Post("/api/links/create","application/json",bytes.NewBuffer(d))
	if err != nil {
		log.Printf("err",err)
		return
	}
	println("status:", resp.StatusCode)
	println("body:", string(body))
}
