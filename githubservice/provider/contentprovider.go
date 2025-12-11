package provider

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/go-github/v39/github"
)

var ErrHttpStatusCode = errors.New("http status code error")

type ContentProvider struct {
	Owner    string
	Repo     string
	Ref      string
	Ctx      context.Context
	GhClient *github.Client
}

func (cp *ContentProvider) GetContents(path string) (string, error) {

	fileContent, _, response, err := cp.GhClient.Repositories.GetContents(cp.Ctx,
		cp.Owner, cp.Repo, path,
		&github.RepositoryContentGetOptions{Ref: cp.Ref})

	if err != nil {
		return "", err
	}
	content, _ := fileContent.GetContent()

	if response.StatusCode != http.StatusOK {
		return content, ErrHttpStatusCode
	}

	return content, nil
}
