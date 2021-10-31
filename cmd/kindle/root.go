package kindle

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "send-to-kindle",
	Version: "0.1.0",
}

func Build() *cobra.Command {
	return rootCmd
}
