package hyprland

import "github.com/sirupsen/logrus"

type client struct {
	logger *logrus.Logger
}

type hyprWorkspace struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
