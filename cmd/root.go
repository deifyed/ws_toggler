package cmd

import (
	"os"
	"strings"

	"github.com/deifyed/wstoggler/cmd/root"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	rootCmdOpts = root.RootCmdOptions{
		Logger: logrus.StandardLogger(),
		Fs:     &afero.Afero{Fs: afero.NewOsFs()},
	}
	rootCmd = &cobra.Command{
		Use:     "wstoggler",
		Short:   "Returns to previous workspace when set",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: root.PreRunE(&rootCmdOpts),
		RunE:    root.RunE(&rootCmdOpts),
	}
)

func init() {
	defaultBackend, ok := os.LookupEnv("XDG_CURRENT_DESKTOP")
	if !ok {
		defaultBackend = "sway"
	}

	defaultBackend = strings.ToLower(defaultBackend)

	rootCmd.Flags().StringVarP(&rootCmdOpts.WorkspaceBackend, "backend", "b", defaultBackend, "Workspace backend to use (sway, hyprland)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		rootCmdOpts.Logger.Errorf(err.Error())

		os.Exit(1)
	}
}
