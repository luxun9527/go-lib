package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)

}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
		log.Println(args)
		log.Println(Verbose)
	},
}

//command 命令 Use
//args 参数  ./cobrademo times 11   输出 Echo: 11 \n false
//flag ./cobrademo version -v=false
