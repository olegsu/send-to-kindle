package kindle

import (
	"fmt"

	"github.com/spf13/cobra"
)

var sendCmd = &cobra.Command{
	Use: "send",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TBD")
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
}
