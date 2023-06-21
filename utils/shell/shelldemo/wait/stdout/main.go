package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	//cmd := exec.CommandContext(ctx, "/bin/bash", "-c", `/home/deng/smb/go-lib/shell/shelldemo/wait/output/output`)
	//cmd.Stdout = os.Stdout
	//if err := cmd.Run(); err != nil {
	//	panic(err)
	//}
	echo()
}
func echo() {
	out, err := exec.Command("date").Output()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("The date is %s\n", out)
}