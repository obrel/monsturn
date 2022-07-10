package cmd

import (
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:     "monsturn",
	Short:   "Coturn statistic monitoring",
	Long:    "Coturn (STUN/TURN) statistic monitoring.",
	Version: "0.0",
	Run: func(command *cobra.Command, args []string) {
		command.Usage()
	},
}

func Execute() error {
	return root.Execute()
}

func Init() {
	root.AddCommand(
		monitorCmd,
	)
}
