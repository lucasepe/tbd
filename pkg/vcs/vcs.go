package vcs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pkg/errors"

	gitUrls "github.com/whilp/git-urls"
)

func InitGitRepo(path string, isBare bool) (*git.Repository, error) {
	repo, err := git.PlainInit(path, true)
	if err != nil {
		if err == git.ErrRepositoryAlreadyExists {
			return git.PlainOpen(path)
		}
		return nil, err
	}
	return repo, nil
}

func OpenGitRepo(path string) (*git.Repository, error) {
	path, err := DetectGitPath(path)
	if err != nil {
		return nil, err
	}

	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("error opening repository '%s'", path))
	}

	return repo, nil
}

func CurrentBranchFromGitRepo(repository *git.Repository) (string, error) {
	branchRefs, err := repository.Branches()
	if err != nil {
		return "", err
	}

	headRef, err := repository.Head()
	if err != nil {
		return "", err
	}

	var currentBranchName string
	err = branchRefs.ForEach(func(branchRef *plumbing.Reference) error {
		if branchRef.Hash() == headRef.Hash() {
			currentBranchName = branchRef.Name().String()

			return nil
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	return currentBranchName, nil
}

func CurrentCommitFromGitRepo(repository *git.Repository) (string, error) {
	headRef, err := repository.Head()
	if err != nil {
		return "", err
	}
	headSha := headRef.Hash().String()

	return headSha, nil
}

func LatestTagFromGitRepo(repository *git.Repository) (string, error) {
	tagRefs, err := repository.Tags()
	if err != nil {
		return "", err
	}

	var latestTagCommit *object.Commit
	var latestTagName string
	err = tagRefs.ForEach(func(tagRef *plumbing.Reference) error {
		revision := plumbing.Revision(tagRef.Name().String())
		tagCommitHash, err := repository.ResolveRevision(revision)
		if err != nil {
			return err
		}

		commit, err := repository.CommitObject(*tagCommitHash)
		if err != nil {
			return err
		}

		if latestTagCommit == nil {
			latestTagCommit = commit
			latestTagName = tagRef.Name().String()
		}

		if commit.Committer.When.After(latestTagCommit.Committer.When) {
			latestTagCommit = commit
			latestTagName = tagRef.Name().String()
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	return latestTagName, nil
}

func GitRepoURL(repo *git.Repository) (string, error) {
	cfg, err := repo.Config()
	if err != nil {
		return "", err
	}

	var res string

	for k, v := range cfg.Remotes {
		if k == "origin" && len(v.URLs) > 0 {
			u, err := gitUrls.Parse(v.URLs[0])
			if err != nil {
				return "", err
			}

			res = u.String()

			if strings.HasPrefix(u.Scheme, "ssh") {
				var sb strings.Builder
				sb.WriteString("https://")
				sb.WriteString(u.Host)
				sb.WriteString("/")
				sb.WriteString(u.Path)

				res = strings.TrimSuffix(sb.String(), ".git")
			}

			break
		}
	}

	return res, nil
}

func DetectGitPath(path string) (string, error) {
	// normalize the path
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	for {
		fi, err := os.Stat(filepath.Join(path, ".git"))
		if err == nil {
			if !fi.IsDir() {
				return "", fmt.Errorf(".git exist but '%s' is not a directory", path)
			}
			return filepath.Join(path, ".git"), nil
		}
		if !os.IsNotExist(err) {
			// unknown error
			return "", err
		}

		// detect bare repo
		ok, err := IsGitDir(path)
		if err != nil {
			return "", err
		}
		if ok {
			return path, nil
		}

		if parent := filepath.Dir(path); parent == path {
			return "", fmt.Errorf(".git not found in '%s'", path)
		} else {
			path = parent
		}
	}
}

func IsGitDir(path string) (bool, error) {
	markers := []string{"HEAD", "objects", "refs"}

	for _, marker := range markers {
		_, err := os.Stat(filepath.Join(path, marker))
		if err == nil {
			continue
		}
		if !os.IsNotExist(err) {
			// unknown error
			return false, err
		} else {
			return false, nil
		}
	}

	return true, nil
}
