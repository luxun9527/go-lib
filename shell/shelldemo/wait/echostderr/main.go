package main

import (
	"log"
	"os"
)

func main() {
	if _, err := os.Stderr.Write([]byte("stderr echo")); err != nil {
		log.Fatal(err)
	}
}
