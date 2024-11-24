package hyprland

import "github.com/sirupsen/logrus"

type client struct {
	logger *logrus.Logger
}

type hyprWorkspace struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
