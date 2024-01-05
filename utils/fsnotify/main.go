package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path/filepath"
	"time"
)

type NotifyFile struct {
	watch *fsnotify.Watcher
}

func NewNotifyFile() *NotifyFile {
	w := new(NotifyFile)
	w.watch, _ = fsnotify.NewWatcher()
	return w
}

//监控目录
func (this *NotifyFile) WatchDir(dir string) {
	//通过Walk来遍历目录下的所有子目录
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		//判断是否为目录，监控目录,目录下文件也在监控范围内，不需要加
		if info.IsDir() {
			path, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			err = this.watch.Add(path)
			if err != nil {
				return err
			}
			fmt.Println("监控 : ", path)
		}
		return nil
	})

	go this.WatchEvent() //协程
}

func (this *NotifyFile) WatchEvent() {
	for {
		select {
		case ev := <-this.watch.Events:
			{
				if ev.Op&fsnotify.Create == fsnotify.Create {
					fmt.Println("创建文件 : ", ev.Name)
					//获取新创建文件的信息，如果是目录，则加入监控中
					//file, err := os.Stat(ev.varName)
					//if err == nil && file.IsDir() {
					//	this.watch.Add(ev.varName)
					//	fmt.Println("添加监控 : ", ev.varName)
					//}
				}

				if ev.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("写入文件 : ", ev.Name)
				}

				if ev.Op&fsnotify.Remove == fsnotify.Remove {
					fmt.Println("删除文件 : ", ev.Name)
					//如果删除文件是目录，则移除监控
					fi, err := os.Stat(ev.Name)
					if err == nil && fi.IsDir() {
						this.watch.Remove(ev.Name)
						fmt.Println("删除监控 : ", ev.Name)
					}
				}

				if ev.Op&fsnotify.Rename == fsnotify.Rename {
					//如果重命名文件是目录，则移除监控 ,注意这里无法使用os.Stat来判断是否是目录了
					//因为重命名后，go已经无法找到原文件来获取信息了,所以简单粗爆直接remove
					fmt.Println("重命名文件 : ", ev.Name)
					this.watch.Remove(ev.Name)
				}
				if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
					fmt.Println("修改权限 : ", ev.Name)
				}
			}
		case err := <-this.watch.Errors:
			{
				fmt.Println("error : ", err)
				return
			}
		}
	}

}

func main() {
	//watch := NewNotifyFile()
	//watch.WatchDir("/home/deng/smb/go-lib/fsnotify/watchfoler")
	watchFile()
	select {}

}
func watchFolder() {
	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	var i int32
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				i++
				log.Printf("i = %v op = %v", i, event.Op)

				if !ok {
					return
				}

				log.Println("event:", event)
				if event.Has(fsnotify.Write) {
					//log.Println("modified file:", event.varName)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add a path.
	err = watcher.Add("/home/deng/test.txt")
}
func watchFile() {
	go func() {
		time.Sleep(time.Second * 3)
		fs, err := os.OpenFile("S:\\go-lib\\fsnotify\\test.txt", os.O_RDWR, os.ModePerm)
		if err != nil {
			log.Printf("open file failed %v\n", err)
		}
		log.Println("open file")
		for i := 0; i < 10; i++ {
			fs.WriteAt([]byte("ffasdfsdfasdfsadf"), int64(i))
			fs.WriteAt([]byte("\n"), int64(i))
			time.Sleep(time.Millisecond * 3)
			//fs.Close()
		}
		fs.Close()
		log.Println("close file")
	}()

	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	var i int32
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				i++
				log.Printf("i = %v op = %v", i, event.Op)

				if !ok {
					return
				}
				log.Printf("event %+v", event.Name)
				if event.Has(fsnotify.Write) {
					//log.Println("modified file:", event.varName)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add a path.
	err = watcher.Add("S:\\go-lib\\fsnotify\\test.txt")
	if err != nil {
		log.Fatal("watcher.add ", err)
	}

	// Block main goroutine forever.
	<-make(chan struct{})
}
