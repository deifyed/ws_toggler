package workspace

type Client interface {
	GetDefaultWorkspace() string
	GetFocusedWorkspace() (string, error)
	SetFocusedWorkspace(string) error
}
