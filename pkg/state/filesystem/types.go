package filesystem

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

const defaultStateFilePermissions = 0o644

type client struct {
	statePath string
	fs        *afero.Afero
	logger    *logrus.Logger
}

type fileState struct {
	PreviousWorkspace string `json:"previous_workspace"`
}
