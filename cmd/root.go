package cmd

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/deifyed/wstoggler/pkg/toggling"

	"github.com/sirupsen/logrus"

	"github.com/deifyed/wstoggler/pkg/state/filesystem"
	"github.com/deifyed/wstoggler/pkg/workspace/sway"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

const (
	defaultStateFilename = "wstoggler.state"
	defaultLogFilename   = "wstoggler.log"
)

type rootCmdOptions struct {
	logger *logrus.Logger
	fs     *afero.Afero
}

var (
	rootCmdOpts = rootCmdOptions{}
	rootCmd     = &cobra.Command{
		Use:   "wstoggler",
		Short: "Returns to previous workspace when set",
		Args:  cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			user, err := user.Current()
			if err != nil {
				return fmt.Errorf("getting current user: %w", err)
			}

			logfilePath := path.Join(os.TempDir(), fmt.Sprintf("%s-%s", user.Username, defaultLogFilename))

			logFile, err := os.OpenFile(logfilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o644)
			if err != nil {
				return fmt.Errorf("opening log file: %w", err)
			}

			rootCmdOpts.logger = &logrus.Logger{
				Out: logFile,
				Formatter: &logrus.JSONFormatter{
					PrettyPrint: true,
				},
				Level: logrus.DebugLevel,
			}

			rootCmdOpts.fs = &afero.Afero{Fs: afero.NewOsFs()}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			user, err := user.Current()
			if err != nil {
				return fmt.Errorf("getting current user: %w", err)
			}

			stateFilePath := path.Join(os.TempDir(), fmt.Sprintf("%s-%s", user.Username, defaultStateFilename))

			toggle := toggling.Toggle{
				Logger:          rootCmdOpts.logger,
				WorkspaceClient: sway.NewWorkspaceClient(rootCmdOpts.logger),
				StateClient: filesystem.NewFilesystemStateClient(
					rootCmdOpts.logger,
					rootCmdOpts.fs,
					stateFilePath,
				),
			}

			if len(args) == 1 {
				rootCmdOpts.logger.Debug(fmt.Sprintf("found argument \"%s\", attempting move to desired workspace", args[0]))

				desiredWorkspace := args[0]

				return toggle.To(desiredWorkspace)
			}

			rootCmdOpts.logger.Debug("no argument found, moving to previous workspace")

			return toggle.Back()
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		rootCmdOpts.logger.Errorf(err.Error())

		os.Exit(1)
	}
}
