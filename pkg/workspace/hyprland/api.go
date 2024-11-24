package hyprland

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/sirupsen/logrus"

	"github.com/deifyed/wstoggler/pkg/workspace"
)

func (c client) GetFocusedWorkspace() (string, error) {
	raw, err := exec.Command("hyprctl", "activeworkspace", "-j").Output()
	if err != nil {
		return "", fmt.Errorf("acquiring workspace output: %w", err)
	}

	var result hyprWorkspace

	err = json.Unmarshal(raw, &result)
	if err != nil {
		return "", fmt.Errorf("unmarshalling workspace result: %w", err)
	}

	return result.Name, nil
}

func (c client) SetFocusedWorkspace(target string) error {
	c.logger.Debug(fmt.Sprintf("setting workspace to %s", target))

	cmd := exec.Command("hyprctl", "dispatch", "workspace", target)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("running command %s: %w", cmd.String(), err)
	}

	return nil
}

func (c client) GetDefaultWorkspace() string {
	return "1"
}

func NewWorkspaceClient(logger *logrus.Logger) workspace.Client {
	return &client{
		logger: logger,
	}
}
