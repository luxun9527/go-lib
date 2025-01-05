package main

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"testing"
	"time"
)

func TestWatchFile(t *testing.T) {
	watchFile("./test.txt")
}

// 监控文件
func watchFile(fileName string) {
	go func() {
		time.Sleep(time.Second * 3)
		fs, err := os.OpenFile(fileName, os.O_RDWR, os.ModePerm)
		if err != nil {
			log.Printf("open file failed %v\n", err)
		}
		fs.Write([]byte("hello world"))
		time.Sleep(time.Second * 3)
		fs.Close()
		//os.Rename(fileName, fileName+"_bak")
		//os.Chmod(fileName+"_bak", 0666)
		os.Remove(fileName)
		log.Println("close file")
	}()

	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	err = watcher.Add(fileName)
	if err != nil {
		log.Panicf("watcher add failed %v ", err)
	}

	for {
		select {
		//监控文件，有变化时执行，create 不会执行。
		case event := <-watcher.Events:
			log.Printf("event %+v", event)
		case err := <-watcher.Errors:
			log.Println("error:", err)
		}
	}

}
func TestWatchFolder(t *testing.T) {
	watchFolder("E:\\openproject\\bestpractice\\server\\accountApi\\lang")
}
func watchFolder(folderName string) {
	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()
	watcher.Add(folderName)
	// Start listening for events.
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Printf("event: %v", event)

			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()
	go func() {
		time.Sleep(time.Second * 3)
		fs, err := os.OpenFile(folderName+"/test.txt", os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
		if err != nil {
			log.Printf("open file failed %v\n", err)
		}
		fs.Write([]byte("hello world"))
		time.Sleep(time.Second * 3)
		fs.Close()
		os.Remove(folderName + "/test.txt")
		log.Println("close file")
	}()
	time.Sleep(time.Hour)
}
