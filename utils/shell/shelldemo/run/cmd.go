package main

import (
	"context"
	"log"
	"os/exec"
)

func main()  {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd := exec.CommandContext(ctx, "sh", "-c", `dd if=/dev/zero of=dest.txt count=0 bs=1MB seek=10`)
	data, err := cmd.CombinedOutput()
	if err!=nil{
		log.Printf("call CombinedOutput failed %v\n",err)
		return
	}
	log.Println(string(data))
}




