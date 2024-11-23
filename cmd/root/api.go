package root

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/deifyed/wstoggler/pkg/state/filesystem"
	"github.com/deifyed/wstoggler/pkg/toggling"
	"github.com/deifyed/wstoggler/pkg/workspace/sway"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	defaultStateFilename = "wstoggler.state"
	defaultLogFilename   = "wstoggler.log"
)

func RunE(opts *RootCmdOptions) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		user, err := user.Current()
		if err != nil {
			return fmt.Errorf("getting current user: %w", err)
		}

		stateFilePath := path.Join(os.TempDir(), fmt.Sprintf("%s-%s", user.Username, defaultStateFilename))

		toggle := toggling.Toggle{
			Logger:          opts.Logger,
			WorkspaceClient: sway.NewWorkspaceClient(opts.Logger),
			StateClient: filesystem.NewFilesystemStateClient(
				opts.Logger,
				opts.Fs,
				stateFilePath,
			),
		}

		if len(args) == 1 {
			opts.Logger.Debug(fmt.Sprintf("found argument \"%s\", attempting move to desired workspace", args[0]))

			desiredWorkspace := args[0]

			return toggle.To(desiredWorkspace)
		}

		opts.Logger.Debug("no argument found, moving to previous workspace")

		return toggle.Back()
	}
}

func PreRunE(opts *RootCmdOptions) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		user, err := user.Current()
		if err != nil {
			return fmt.Errorf("getting current user: %w", err)
		}

		logfilePath := path.Join(os.TempDir(), fmt.Sprintf("%s-%s", user.Username, defaultLogFilename))

		logFile, err := os.OpenFile(logfilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o644)
		if err != nil {
			return fmt.Errorf("opening log file: %w", err)
		}

		opts.Logger = &logrus.Logger{
			Out: logFile,
			Formatter: &logrus.JSONFormatter{
				PrettyPrint: true,
			},
			Level: logrus.DebugLevel,
		}

		return nil
	}
}
