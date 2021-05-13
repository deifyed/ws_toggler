package sway

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"

	"github.com/sirupsen/logrus"

	"github.com/deifyed/wstoggler/pkg/workspace"
)

const DefaultWorkspace = "1"

func (c client) GetFocusedWorkspace() (string, error) {
	raw, err := exec.Command("swaymsg", "-t", "get_workspaces", "-r").Output()
	if err != nil {
		return "", fmt.Errorf("acquiring workspace output: %w", err)
	}

	var result getWorkspaceResult

	err = json.Unmarshal(raw, &result)
	if err != nil {
		return "", fmt.Errorf("unmarshalling workspace result: %w", err)
	}

	for _, ws := range result {
		if ws.Focused {
			c.logger.Debug(fmt.Sprintf("found focused namespace: %s", ws.Name))

			return ws.Name, nil
		}
	}

	c.logger.Debug("could not find focused namespace")

	return "", errors.New("finding focused workspace")
}

func (c client) SetFocusedWorkspace(target string) error {
	c.logger.Debug(fmt.Sprintf("setting workspace to %s", target))

	cmd := exec.Command("swaymsg", "--", "workspace", target)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("running \"swaymsg -- workspace %s\" command: %w", target, err)
	}

	return nil
}

func NewWorkspaceClient(logger *logrus.Logger) workspace.Client {
	return &client{
		logger: logger,
	}
}
