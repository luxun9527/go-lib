package inject

import (
	"flag"
	"fmt"
	"os"
)

var (
	builtAt   string
	goVersion string
	gitAuthor string
	gitCommit string
	buildUser string
)

func init(){
	flag.Bool("version",false,"打印版本信息")
	Register("version",PrintVersionInfo)
}
func PrintVersionInfo() {
	fmt.Printf("%-20s %s\n", "builtAt", builtAt)
	fmt.Printf("%-20s %s\n", "buildUser", buildUser)
	fmt.Printf("%-20s %s\n", "goVersion", goVersion)
	fmt.Printf("%-20s %s\n", "gitAuthor", gitAuthor)
	fmt.Printf("%-20s %s\n", "gitCommit", gitCommit)
	os.Exit(1)
}

