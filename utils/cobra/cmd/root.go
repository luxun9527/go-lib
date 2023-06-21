package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var Verbose bool
var Source string

var rootCmd = &cobra.Command{
	Use:   "hugo [string to echo]",
	Short: "Hugo is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at http://hugo.spf13.com`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		fmt.Println("root file hogo commend invoke")
		fmt.Println(Verbose)
		fmt.Println(Source)
	},
}

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
