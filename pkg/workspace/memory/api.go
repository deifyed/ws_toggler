package memory

import "github.com/deifyed/wstoggler/pkg/workspace"

type client struct {
	focused string
}

func (c *client) GetDefaultWorkspace() string {
	return "default"
}

func (c *client) GetFocusedWorkspace() (string, error) {
	return c.focused, nil
}

func (c *client) SetFocusedWorkspace(s string) error {
	c.focused = s

	return nil
}

func NewWorkspaceClient() workspace.Client {
	return &client{}
}
