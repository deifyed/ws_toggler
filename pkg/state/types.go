package state

type Client interface {
	GetPreviousWorkspace(defaultWorkspace string) (string, error)
	StorePreviousWorkspace(string) error
}
