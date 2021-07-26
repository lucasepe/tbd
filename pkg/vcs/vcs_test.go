package vcs

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGitRepo(t *testing.T) {
	// Plain
	plainRoot, err := ioutil.TempDir("", "")
	require.NoError(t, err)
	defer os.RemoveAll(plainRoot)

	_, err = InitGitRepo(plainRoot, false)
	require.NoError(t, err)
	plainGitDir := filepath.Join(plainRoot, ".git")

	tests := []struct {
		inPath  string
		outPath string
		err     bool
	}{
		// errors
		{"/", "", true},
		// parent dir of a repo
		{filepath.Dir(plainRoot), "", true},

		// Plain repo
		{plainRoot, plainGitDir, false},
		{plainGitDir, plainGitDir, false},
		{path.Join(plainGitDir, "objects"), plainGitDir, false},
	}

	for i, tc := range tests {
		dir, err := DetectGitPath(tc.inPath)
		if tc.err {
			require.Error(t, err, i)
		}

		_, err = OpenGitRepo(tc.inPath)
		if tc.err {
			require.Error(t, err, i)
		} else {
			require.NoError(t, err, i)
			assert.Equal(t, filepath.ToSlash(tc.outPath), filepath.Join(filepath.ToSlash(dir), ".git"), i)
		}
	}
}
