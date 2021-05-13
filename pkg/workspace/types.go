package workspace

type Client interface {
	GetFocusedWorkspace() (string, error)
	SetFocusedWorkspace(string) error
}
