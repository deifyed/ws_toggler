package toggling

import (
	"io"
	"testing"

	"github.com/deifyed/wstoggler/pkg/state/filesystem"
	"github.com/deifyed/wstoggler/pkg/workspace/memory"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func mockLogger() *logrus.Logger {
	logger := logrus.New()
	logger.Out = io.Discard

	return logger
}

func mockFs() *afero.Afero {
	return &afero.Afero{Fs: afero.NewMemMapFs()}
}

type target struct {
	workspace string
}

func TestWorkspace(t *testing.T) {
	testCases := []struct {
		name                 string
		withToggler          Toggle
		withTargetSequence   []target
		expectPreviousTarget string
	}{
		{
			name: "Should work",
			withToggler: Toggle{
				Logger:          mockLogger(),
				StateClient:     filesystem.NewFilesystemStateClient(mockLogger(), mockFs(), "mockState"),
				WorkspaceClient: memory.NewWorkspaceClient(),
			},
			withTargetSequence: []target{
				{workspace: "1"},
				{workspace: "3"},
				{workspace: "7"},
			},
			expectPreviousTarget: "3",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			for _, target := range tc.withTargetSequence {
				err := tc.withToggler.To(target.workspace)
				assert.NoError(t, err)
			}

			result, err := tc.withToggler.StateClient.GetPreviousWorkspace(
				tc.withToggler.WorkspaceClient.GetDefaultWorkspace(),
			)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectPreviousTarget, result)
		})
	}
}

func TestBack(t *testing.T) {
	testCases := []struct {
		name                 string
		withToggler          Toggle
		withTargetSequence   []target
		withBackCalls        int
		expectPreviousTarget string
	}{
		{
			name: "Should work",
			withToggler: Toggle{
				Logger:          mockLogger(),
				StateClient:     filesystem.NewFilesystemStateClient(mockLogger(), mockFs(), "mockState"),
				WorkspaceClient: memory.NewWorkspaceClient(),
			},
			withTargetSequence: []target{
				{workspace: "1"},
				{workspace: "3"},
				{workspace: "7"},
			},
			withBackCalls:        1,
			expectPreviousTarget: "7",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			for _, target := range tc.withTargetSequence {
				err := tc.withToggler.To(target.workspace)
				assert.NoError(t, err)
			}

			for i := 0; i < tc.withBackCalls; i++ {
				err := tc.withToggler.Back()
				assert.NoError(t, err)
			}

			defaultWorkspace := tc.withToggler.WorkspaceClient.GetDefaultWorkspace()

			result, err := tc.withToggler.StateClient.GetPreviousWorkspace(defaultWorkspace)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectPreviousTarget, result)
		})
	}
}
