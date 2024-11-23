package root

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type RootCmdOptions struct {
	Logger           *logrus.Logger
	Fs               *afero.Afero
	WorkspaceBackend string
}
