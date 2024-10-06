package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

var echoTimes int

// go run main.go times -t=10 34
var cmdTimes = &cobra.Command{
	Use:   "times [string to echo]",
	Short: "Echo anything to the screen more times",
	Long: `echo things multiple times back to the user by providing
a count and a string.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for i := 0; i < echoTimes; i++ {
			fmt.Println("Echo: " + strings.Join(args, " "))
		}
	},
}

func init() {
	//全局flag
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&Source, "source", "s", "test1212", "verbose output")
	//单个命令的flag
	cmdTimes.Flags().IntVarP(&echoTimes, "times", "t", 1, "times to echo the input")
	rootCmd.AddCommand(cmdTimes)
}
