package filesystem

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/deifyed/wstoggler/pkg/state"
	"github.com/spf13/afero"
)

func (c client) GetPreviousWorkspace(defaultWorkspace string) (string, error) {
	if exists, _ := c.fs.Exists(c.statePath); !exists {
		return defaultWorkspace, nil
	}

	stateFile, err := c.fs.OpenFile(c.statePath, os.O_RDONLY, defaultStateFilePermissions)
	if err != nil {
		return "", fmt.Errorf("opening state file: %w", err)
	}

	defer func() {
		_ = stateFile.Close()
	}()

	raw, err := io.ReadAll(stateFile)
	if err != nil {
		return "", fmt.Errorf("reading state file: %w", err)
	}

	fState := fileState{}

	err = json.Unmarshal(heal(raw), &fState)
	if err != nil {
		return "", fmt.Errorf("unmarshalling state file: %w", err)
	}

	return fState.PreviousWorkspace, nil
}

func (c client) StorePreviousWorkspace(workspace string) error {
	raw, err := json.Marshal(&fileState{PreviousWorkspace: workspace})
	if err != nil {
		return fmt.Errorf("marshalling state: %w", err)
	}

	stateFile, err := c.fs.OpenFile(c.statePath, os.O_WRONLY|os.O_CREATE, defaultStateFilePermissions)
	if err != nil {
		return fmt.Errorf("opening state file: %w", err)
	}

	defer func() {
		_ = stateFile.Close()
	}()

	_, err = stateFile.Write(raw)
	if err != nil {
		return fmt.Errorf("writing to state file: %w", err)
	}

	return nil
}

func NewFilesystemStateClient(l *logrus.Logger, fs *afero.Afero, stateFilePath string) state.Client {
	return &client{
		logger:    l,
		statePath: stateFilePath,
		fs:        fs,
	}
}
