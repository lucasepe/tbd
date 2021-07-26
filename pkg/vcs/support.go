package vcs

import (
	"net/url"
	"strings"
)

const (
	RepoCommit   = "REPO_COMMIT"
	RepoTag      = "REPO_TAG"
	RepoTagClean = "REPO_TAG_CLEAN"
	RepoURL      = "REPO_URL"
	RepoHost     = "REPO_HOST"
	RepoName     = "REPO_NAME"
	RepoRoot     = "REPO_ROOT"
)

func GitRepoMetadata(path string, meta map[string]string) error {
	repo, err := OpenGitRepo(path)
	if err != nil {
		return err
	}

	commit, err := CurrentCommitFromGitRepo(repo)
	if err != nil {
		return err
	}
	meta[RepoCommit] = commit

	tag, err := LatestTagFromGitRepo(repo)
	if err != nil {
		return err
	}
	idx := strings.LastIndex(tag, "/")
	if idx != -1 {
		tag = tag[idx+1:]
	}
	meta[RepoTag] = tag

	if strings.HasPrefix(tag, "v") {
		meta[RepoTagClean] = tag[1:]
	}

	repoURL, err := GitRepoURL(repo)
	if err != nil {
		return err
	}
	meta[RepoURL] = repoURL

	if u, err := url.Parse(repoURL); err == nil {
		idx := strings.Index(u.Path[1:], "/")
		if idx != -1 {
			meta[RepoRoot] = u.Path[1 : idx+1]
		}

		meta[RepoHost] = u.Host
	}

	meta[RepoName] = repoURL[strings.LastIndex(repoURL, "/")+1:]

	return nil
}
