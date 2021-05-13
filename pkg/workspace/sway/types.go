package sway

import "github.com/sirupsen/logrus"

type client struct {
	logger *logrus.Logger
}

type swayWorkspace struct {
	Name    string `json:"name"`
	Focused bool   `json:"focused"`
}

type getWorkspaceResult []swayWorkspace
