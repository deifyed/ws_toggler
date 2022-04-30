package toggling

import (
	"github.com/deifyed/wstoggler/pkg/state"
	"github.com/deifyed/wstoggler/pkg/workspace"
	"github.com/sirupsen/logrus"
)

type Toggle struct {
	Logger          *logrus.Logger
	StateClient     state.Client
	WorkspaceClient workspace.Client
}
