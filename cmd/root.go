package cmd

import (
	"os"

	"github.com/deifyed/wstoggler/cmd/root"
	"github.com/spf13/cobra"
)

var (
	rootCmdOpts = root.RootCmdOptions{}
	rootCmd     = &cobra.Command{
		Use:     "wstoggler",
		Short:   "Returns to previous workspace when set",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: root.PreRunE(&rootCmdOpts),
		RunE:    root.RunE(&rootCmdOpts),
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		rootCmdOpts.Logger.Errorf(err.Error())

		os.Exit(1)
	}
}
