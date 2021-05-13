package cmd

import (
	"fmt"
	"os"

	"github.com/deifyed/wstoggler/pkg/state"
	"github.com/deifyed/wstoggler/pkg/workspace"
	"github.com/sirupsen/logrus"

	"github.com/deifyed/wstoggler/pkg/state/filesystem"
	"github.com/deifyed/wstoggler/pkg/workspace/sway"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

const (
	defaultStatePath = "/tmp/wstoggler.state"
	defaultLogPath   = "/tmp/wstoggler.log"
)

var logger *logrus.Logger

var rootCmd = &cobra.Command{
	Use:   "wstoggler",
	Short: "Returns to previous workspace when set",
	Args:  cobra.MaximumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		logFile, err := os.OpenFile(defaultLogPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o644)
		if err != nil {
			return fmt.Errorf("opening log file: %w", err)
		}

		logger = &logrus.Logger{
			Out: logFile,
			Formatter: &logrus.JSONFormatter{
				PrettyPrint: true,
			},
			Level: logrus.DebugLevel,
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		fs := afero.Afero{Fs: afero.NewOsFs()}

		workspaceClient := sway.NewWorkspaceClient(logger)
		stateClient := filesystem.NewFilesystemStateClient(logger, &fs, defaultStatePath)

		moveTo := generateMoveFunction(stateClient, workspaceClient)

		if len(args) == 1 {
			logger.Debug(fmt.Sprintf("found argument \"%s\", attempting move to desired workspace", args[0]))

			desiredWorkspace := args[0]

			return moveTo(desiredWorkspace)
		}

		logger.Debug("no argument found, moving to previous workspace")

		previousWorkspace, err := stateClient.GetPreviousWorkspace(sway.DefaultWorkspace)
		if err != nil {
			return fmt.Errorf("getting previous workspace: %w", err)
		}

		return moveTo(previousWorkspace)
	},
}

func generateMoveFunction(stateClient state.Client, workspaceClient workspace.Client) func(string) error {
	return func(target string) error {
		currentWorkspace, err := workspaceClient.GetFocusedWorkspace()
		if err != nil {
			return fmt.Errorf("getting current workspace: %w", err)
		}

		if target == currentWorkspace {
			logger.Debug(fmt.Sprintf(
				"target workspace(%s) is equal to current namespace(%s), do nothing",
				target,
				currentWorkspace,
			))

			return nil
		}

		err = workspaceClient.SetFocusedWorkspace(target)
		if err != nil {
			return fmt.Errorf("setting focused workspace: %w", err)
		}

		err = stateClient.StorePreviousWorkspace(currentWorkspace)
		if err != nil {
			return fmt.Errorf("storing previous workspace: %w", err)
		}

		return nil
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Errorf(err.Error())

		os.Exit(1)
	}
}
