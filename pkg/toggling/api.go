package toggling

import (
	"fmt"
)

func (c Toggle) To(targetWorkspace string) error {
	currentWorkspace, err := c.WorkspaceClient.GetFocusedWorkspace()
	if err != nil {
		return fmt.Errorf("getting current workspace: %w", err)
	}

	if targetWorkspace == currentWorkspace {
		c.Logger.Debug(fmt.Sprintf(
			"target workspace(%s) is equal to current namespace(%s), do nothing",
			targetWorkspace,
			currentWorkspace,
		))

		return nil
	}

	err = c.WorkspaceClient.SetFocusedWorkspace(targetWorkspace)
	if err != nil {
		return fmt.Errorf("setting focused workspace: %w", err)
	}

	err = c.StateClient.StorePreviousWorkspace(currentWorkspace)
	if err != nil {
		return fmt.Errorf("storing previous workspace: %w", err)
	}

	return nil
}

func (c Toggle) Back() error {
	previousWorkspace, err := c.StateClient.GetPreviousWorkspace(c.WorkspaceClient.GetDefaultWorkspace())
	if err != nil {
		return fmt.Errorf("retrieving previous workspace: %w", err)
	}

	currentWorkspace, err := c.WorkspaceClient.GetFocusedWorkspace()
	if err != nil {
		return fmt.Errorf("retrieving current workspace: %w", err)
	}

	err = c.StateClient.StorePreviousWorkspace(currentWorkspace)
	if err != nil {
		return fmt.Errorf("storing current namespace: %w", err)
	}

	err = c.WorkspaceClient.SetFocusedWorkspace(previousWorkspace)
	if err != nil {
		return fmt.Errorf("setting focused namespace: %w", err)
	}

	return nil
}
