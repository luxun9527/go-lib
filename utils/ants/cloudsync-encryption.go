package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var CmdProcess *exec.Cmd

type Parameter struct {
	Download bool   `json:"download"` // 仅下载
	Upload   bool   `json:"upload"`   // 仅上传
	Threads  int    `json:"threads"`  // 下载线程数
	Time     int    `json:"time"`     // 滚动时间 默认 10
	Local    string `json:"local"`    // 本地同步路径
	Limit    string `json:"limit"`    // 带宽限制
	Process  string `json:"process"`  // 进程名称
	Size     string `json:"size"`     // 最大同步文件 5K, 5M, 5G
	Remote   string `json:"remote"`   // 同步远程路径
	Rule     string `json:"rule"`     // 过滤规则
}

func main()  {
	par := getArgs()
	shell := GetShell(par)
	statusPath := mkDir(par.Process)
	StartSync(shell, statusPath, par.Time)
}

func getArgs() Parameter {
	var par Parameter
	// 接收命令行参数
	fs := flag.NewFlagSet("cloudsync-encryption", flag.ExitOnError)
	fs.BoolVar(&par.Download,"d", false,"Download Only.")                      // 仅下载
	fs.BoolVar(&par.Upload,"u", false,"Upload Only.")                          // 仅上传
	fs.IntVar(&par.Threads, "h", 0,"download threads number for download")     // 下载线程数
	fs.StringVar(&par.Local, "l", "NOT NULL","Sync for local path.")           // 本地同步路径
	fs.StringVar(&par.Limit, "m", "","Bandwidth limit")                        // 带宽限制
	fs.StringVar(&par.Process, "p", "NOT NULL","Your Process Name.")           // 进程名称
	fs.StringVar(&par.Remote, "r", "NOT NULL","Sync for remote path.")         // 同步远程路径
	fs.StringVar(&par.Size, "s", "","max size type examples: 5K, 5M, 5G")      // 最大同步文件 5K, 5M, 5G
	fs.StringVar(&par.Rule, "f", "","filter rule file path")                   // 过滤规则
	fs.IntVar(&par.Time, "t", 10,"Rolling time(Minute unit)")                  // 滚动时间 默认 10
	// 解析命令行参数写入注册的flag里
	_ = fs.Parse(os.Args[1:])
	if len(fs.Args()) > 0 || par.Process == "NOT NULL" || par.Remote == "NOT NULL" || par.Local == "NOT NULL" {
		ExitMsg(fs)
	}
	if par.Download && par.Upload {
		ExitMsg(fs)
	}
	if !par.Download && !par.Upload {
		ExitMsg(fs)
	}
	return par
}

func GetShell(par Parameter) (shell string) {
	isUpload, allPath := GetSyncMod(par)
	var size string
	if len(par.Size) > 0 {
		size = fmt.Sprintf("--max-size %s", par.Size)
	}
	appRootPath := GetWorkPath()
	syncPath := path.Join(appRootPath,"/bin/cloudsync")
	logPath := path.Join(appRootPath,fmt.Sprintf("/log/%s.log", isUpload))
	progress := path.Join(appRootPath,fmt.Sprintf("/log/%s.progress", isUpload))
	str := "%s sync %s --ignore-errors --log-file='%s' --filter-from %s %s --create-empty-src-dirs --progress >> '%s'"
	shell = fmt.Sprintf(str, syncPath, allPath, logPath, par.Rule, size, progress)
	return
}

// GetSyncMod /** 获取同步模式
func GetSyncMod(par Parameter) (isUpload, allPath string) {
	if par.Download {
		isUpload = fmt.Sprintf("/%s_download", par.Process)
		allPath = fmt.Sprintf("'%s' '%s'", par.Remote, par.Local)
	} else if par.Upload {
		isUpload = fmt.Sprintf("/%s_upload", par.Process)
		allPath = fmt.Sprintf("'%s' '%s'",par.Local , par.Remote)
	}
	return isUpload, allPath
}

func mkDir(process string) string {
	appRootPath := GetWorkPath()
	statusPath := path.Join(appRootPath, "status")
	_ = os.MkdirAll(statusPath, os.ModePerm)
	statusPath = path.Join(statusPath, process)
	return statusPath
}

// StartSync /** 开始任务
func StartSync(shell, statusPath string, Time int) {
	_ = SetFile([]byte("processing"), statusPath)
	fmt.Println(shell)
	_, _ = CmdStart(shell)
	str := fmt.Sprintf("Wait %d minutes for the completion of the run", Time)
	fmt.Println(str)
	_ = SetFile([]byte("done"), statusPath)
	time.Sleep(time.Duration(Time*60)*time.Second)
	StartSync(shell, statusPath, Time)
}

// CmdStart /** 执行shell
func CmdStart(shell string) (string, error) {
	cmd := exec.Command("sh", "-c", shell)
	CmdProcess = cmd
	bytes, err := cmd.Output()
	if err != nil {
		return "null", err
	}
	return string(bytes), nil
}

// GetWorkPath /** 获取当前go程序工作目录
func GetWorkPath() string {
	dir, err := filepath.Abs(path.Dir(os.Args[0]))
	if err != nil {
		return "/"
	}
	dir = strings.Replace(dir,"\\","/", -1)
	return path.Dir(dir)
}

// Exists /** 查看某个目录或文件是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// SetFile /** 写入文件内容
func SetFile(data []byte, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func ExitMsg(fs *flag.FlagSet) {
	fs.Usage()
	log.Fatalf("./cloudsync [-p process_name] [-m 0.2M/2M] [-s 100K/100M/1G] [-h multi threads numbers] [-f --filter-from] [-e /filter/some/dir/or/files] [-i /filter/some/dir/or/files] [-u -d -ud] [-r remote_path] [-l local_path] [-t 10]")
}