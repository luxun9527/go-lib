package main

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//todo 空文件
//todo 忽略文件

// FileStatus 文件的各个状态
type FileStatus int32

const (
	Unchanged FileStatus = 1 << iota
	NotEXIST
	Newer
	New
	Delete
)

// OperateType 定义各种操作类型
type OperateType int32

const (
	NewNotEXIST      = OperateType(New | NotEXIST)
	NewerUnchanged   = OperateType(Newer | Unchanged)
	DeletedUnchanged = OperateType(Delete | Unchanged)
	NewNew           = OperateType(New | New)
	NewerNewer       = OperateType(Newer | Newer)
	NewerDeleted     = OperateType(Newer | Delete)
)

func (fs OperateType) Operate(src, dist *FileInfo, priority FileType) error {
	switch fs {
	case NewNotEXIST, NewerUnchanged, NewerDeleted:
		c := fmt.Sprintf("rclone -P copyto %s %s", src.FullPath, dist.FullPath)
		cmd := exec.Command("/bin/bash", "-c", c)
		if err := cmd.Run(); err != nil {
			log.Println("exec cmd err", err)
		}
		//自维护本地md5,避免从远程再拉一份,提升效率
		LastLocalEncodeList[src.FileName] = src.EncodeValue
		LastRemoteEncodeList[src.FileName] = src.EncodeValue
	case DeletedUnchanged:
		c := fmt.Sprintf("rclone -P delete %s ", dist.FullPath)
		cmd := exec.Command("/bin/bash", "-c", c)
		if err := cmd.Run(); err != nil {
			log.Println("exec cmd err", err)
		}
		delete(LastLocalEncodeList, src.FileName)
		delete(LastRemoteEncodeList, src.FileName)
	case NewNew, NewerNewer:
		//在这种情况，src一定是本地，dist一定是远程。
		//md5或sha1这些编码相同则不进行操作
		if src.EncodeValue == dist.EncodeValue {
			LastLocalEncodeList[src.FileName] = src.EncodeValue
			LastRemoteEncodeList[src.FileName] = src.EncodeValue
			return nil
		}

		var (
			local         string
			remote        string
			localToRemote string
			remoteToLocal string
		)
		if priority == Local {
			local = src.FullPath + ".local"
			remote = dist.FullPath + ".remote"
			localToRemote = dist.FullPath + ".local"
			remoteToLocal = src.FullPath + ".remote"
		} else {
			remote = src.FullPath + ".remote"
			local = dist.FullPath + ".local"
			remoteToLocal = dist.FullPath + ".remote"
			localToRemote = src.FullPath + ".local"
		}
		c := fmt.Sprintf("rclone -P moveto %s %s", src.FullPath, local)
		cmd := exec.Command("/bin/bash", "-c", c)
		if err := cmd.Run(); err != nil {
			log.Println("exec cmd err", err)
		}
		delete(LastLocalEncodeList, src.FileName)
		LastLocalEncodeList[src.FileName+".local"] = src.EncodeValue

		c = fmt.Sprintf("rclone -P moveto %s %s", dist.FullPath, remote)
		cmd = exec.Command("/bin/bash", "-c", c)
		if err := cmd.Run(); err != nil {
			log.Println("exec cmd err", err)
		}
		delete(LastRemoteEncodeList, dist.FileName)
		LastRemoteEncodeList[dist.FileName+".remote"] = dist.EncodeValue
		//上传到远程
		c = fmt.Sprintf("rclone -P copyto %s %s", local, localToRemote)
		cmd = exec.Command("/bin/bash", "-c", c)
		if err := cmd.Run(); err != nil {
			log.Println("exec cmd err", err)
		}
		LastRemoteEncodeList[src.FileName+".local"] = src.EncodeValue

		//下载到本地
		c = fmt.Sprintf("rclone -P copyto %s %s", remote, remoteToLocal)
		cmd = exec.Command("/bin/bash", "-c", c)
		if err := cmd.Run(); err != nil {
			log.Println("exec cmd err", err)
		}
		LastLocalEncodeList[dist.FileName+".remote"] = dist.EncodeValue

	}
	return nil
}

const blank = "  "

type FileType int32

const (
	Local FileType = iota
	Remote
)

const (
	MD5  = "md5sum"
	SHA1 = "sha1sum"
	Hash = "hash"
)

var LastLocalEncodeList = map[string]string{}
var LastRemoteEncodeList = map[string]string{}

var (
	LocalRoot  = "/home/deng/cloudsync"
	RemoteRoot = "test:/cloudsync2"
)

type FileInfo struct {
	FileName    string
	Status      FileStatus
	EncodeValue string
	FullPath    string
}

func main() {
	//获取远程文件的sha1

	//空文件

	engine := gin.New()
	engine.GET("/test", func(c *gin.Context) {
		//step1 获取文件的各个状态
		localFileStatus, err := MarkStatus(LocalRoot, SHA1, LastLocalEncodeList)
		if err != nil {
			return
		}

		remoteFileStatus, err := MarkStatus(RemoteRoot, SHA1, LastRemoteEncodeList)
		if err != nil {
			return
		}
		//step2 比较并且按照指定的策略来执行。
		//三种情况， local文件有操作 vs remote 无操作， local文件有操作 vs remote有操作 local文件无操作vs remote文件有操作。
		//local文件有操作 vs remote 无操作， local文件有操作 vs remote有操作
		for k, v := range localFileStatus {
			remoteFileChanges, ok := remoteFileStatus[k]
			var operate FileStatus
			var (
				src,
				dist *FileInfo
			)
			//不存在表示对方没有改变,或对方不存在。
			if !ok {
				if v.Status == New {
					//v.Status NotEXIST
					operate = v.Status | NotEXIST
				} else {
					//v.Status Unchanged
					operate = v.Status | Unchanged
				}
				dist = &FileInfo{
					FileName: v.FileName,
					FullPath: filepath.Join(RemoteRoot, v.FileName),
				}
			} else {
				//v.Status vs  fileInfo.Status 都存在则执行策略
				operate = v.Status | remoteFileChanges.Status
				dist = &FileInfo{
					FileName:    remoteFileChanges.FileName,
					EncodeValue: remoteFileChanges.EncodeValue,
					FullPath:    filepath.Join(RemoteRoot, remoteFileChanges.FileName),
				}
			}
			src = &FileInfo{
				FileName:    v.FileName,
				EncodeValue: v.EncodeValue,
				FullPath:    filepath.Join(LocalRoot, v.FileName),
			}
			if err := OperateType(operate).Operate(src, dist, Local); err != nil {
				log.Println("err", err)
			}
		}
		//local文件无操作vs remote文件有操作。
		for k, v := range remoteFileStatus {
			_, ok := localFileStatus[k]
			var operate FileStatus
			if !ok {
				if v.Status == New {
					operate = v.Status | NotEXIST
				} else {
					//操作Unchanged
					operate = v.Status | Unchanged
				}
				//dist := filepath.Join(RemoteRoot, v.FileName)
				//src := filepath.Join(LocalRoot, v.FileName)
				src := &FileInfo{
					FileName:    v.FileName,
					EncodeValue: v.EncodeValue,
					FullPath:    filepath.Join(RemoteRoot, v.FileName),
				}
				dist := &FileInfo{
					FileName:    v.FileName,
					EncodeValue: v.EncodeValue,
					FullPath:    filepath.Join(LocalRoot, v.FileName),
				}
				if err := OperateType(operate).Operate(src, dist, Remote); err != nil {
					log.Println("err", err)
				}
			}
		}
		//step3 更新双方的md5
		//LastLocalEncodeList, err = getEncodeList(LocalRoot, SHA1)
		//if err != nil {
		//	log.Println("err", err)
		//}
		//
		//LastRemoteEncodeList, err = getEncodeList(RemoteRoot, SHA1)
		//if err != nil {
		//	log.Println("err", err)
		//}
	})
	engine.Run(":9999")
}

// MarkStatus 标记状态
func MarkStatus(path, encodeMethod string, lastCache map[string]string) (map[string]*FileInfo, error) {
	currentEncodeList, err := getEncodeList(path, encodeMethod)
	if err != nil {
		return nil, err
	}

	changedFile := make(map[string]*FileInfo, 10)
	for k, v := range currentEncodeList {

		lastEncodeValue, ok := lastCache[k]
		//这次有上传没有，新建
		if !ok {
			f := &FileInfo{
				FileName:    k,
				Status:      New,
				EncodeValue: v,
			}
			changedFile[k] = f
			continue
		}
		//相比于上一次，编码不一样
		if lastEncodeValue != v {
			f := &FileInfo{
				FileName:    k,
				Status:      Newer,
				EncodeValue: v,
			}
			changedFile[k] = f
			continue
		}

	}
	//在上次的缓存中的map中存在但是在当前的时候不存在
	for key := range lastCache {
		if _, ok := currentEncodeList[key]; !ok {
			f := &FileInfo{
				FileName: key,
				Status:   Delete,
			}
			changedFile[key] = f
		}
	}

	return changedFile, nil
}

// GetEncodeList 获取列表的编码集合
func getEncodeList(path string, encodeMethod string) (map[string]string, error) {

	outputFile := "/home/deng/smb/go-lib/rclone/" + filepath.Base(path) + "_" + encodeMethod
	//空文件
	log.Println(outputFile)
	command := exec.Command("/bin/bash", "-c", fmt.Sprintf(" rclone %s %s --output-file %s ", encodeMethod, path, outputFile))
	if err := command.Run(); err != nil {
		log.Println("debug1 err", err)
		return nil, err
	}

	fi, err := os.Open(outputFile)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return nil, err
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	m := make(map[string]string, 10)
	for {
		data, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		fmt.Println(string(data))
		value := strings.Split(string(data), blank)
		if len(value) < 2 {
			continue
		}
		v := value[0]
		k := value[1]
		m[k] = v
	}
	return m, nil
}
