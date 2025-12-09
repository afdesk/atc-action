package push

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/smartforce-io/atc/githubservice/fetcher"
	"github.com/smartforce-io/atc/githubservice/fetcher/customregex"
	"github.com/smartforce-io/atc/githubservice/fetcher/packagejson"
	"github.com/smartforce-io/atc/githubservice/fetcher/yaml/pluginyaml"
	"github.com/smartforce-io/atc/githubservice/fetcher/yaml/pubspecyaml"
	"github.com/smartforce-io/atc/githubservice/gitutil"
	"github.com/smartforce-io/atc/githubservice/provider"
	"github.com/smartforce-io/atc/githubservice/settings"
	"golang.org/x/oauth2"

	"github.com/google/go-github/v39/github"

	"github.com/smartforce-io/atc/githubservice/fetcher/buildgradle"
	"github.com/smartforce-io/atc/githubservice/fetcher/pomxml"
)

type TagContent struct {
	Version string
}

var autoFetchers = map[string]fetcher.VersionFetcher{
	"pom.xml":      &pomxml.Fetcher{},
	"build.gradle": &buildgradle.Fetcher{},
	"package.json": &packagejson.Fetcher{},
	"pubspec.yaml": &pubspecyaml.Fetcher{},
	"plugin.yaml":  &pluginyaml.Fetcher{},
}

func detectFetchType(path string) string {
	if path == "" {
		return ""
	}
	return filepath.Base(path)
}

func renderTagNameTemplate(templateString, version string) (string, error) {
	buf := new(bytes.Buffer)
	tagContent := TagContent{version}
	tmplFuncMap := template.FuncMap{
		"Time": func() time.Time { return time.Now() },
	}
	tmpl, err := template.New("template tagContent").Funcs(tmplFuncMap).Parse(templateString)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(buf, tagContent)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func AddTag() error {
	githubToken := os.Getenv("GITHUB_TOKEN")
	fullname := os.Getenv("GITHUB_REPOSITORY")
	commitSHA := os.Getenv("COMMIT_SHA")

	atcs := &settings.AtcSettings{
		Path:     os.Getenv("FILE_TYPE"),
		Behavior: os.Getenv("BEHAVIOR"),
		Template: os.Getenv("TEMPLATE"),
		RegexStr: os.Getenv("REGEX"),
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	s := strings.Split(fullname, "/")
	owner := s[0]
	repo := s[1]

	commit, _, err := client.Repositories.GetCommit(ctx, owner, repo, commitSHA, nil)
	if err != nil {
		return fmt.Errorf("error getting commit %s %v", commitSHA, err)
	}

	parents := commit.Parents
	if len(parents) == 0 {
		log.Printf("this branch has no older commits")
		return nil
	}

	ghOldContentProviderPtr := &provider.GhContentProvider{
		Owner:    owner,
		Repo:     repo,
		Ref:      parents[0].GetSHA(),
		Ctx:      ctx,
		GhClient: client,
	}
	ghNewContentProviderPtr := &provider.GhContentProvider{
		Owner:    owner,
		Repo:     repo,
		Ref:      commitSHA,
		Ctx:      ctx,
		GhClient: client,
	}

	var sha string

	caption, err := fetch(atcs, ghOldContentProviderPtr, ghNewContentProviderPtr, fullname)
	if err != nil {
		return fmt.Errorf("fetch version error: %v", err)
	}
	if caption == "" {
		log.Printf("Old and new versions are equal")
		return nil
	}

	if atcs.Behavior == settings.BehaviorAfter {
		sha = commit.GetSHA()
	} else {
		sha = parents[0].GetSHA()
	}

	objType := "commit"
	timestamp := time.Now()

	tag := &github.Tag{
		Tag:     &caption,
		Message: &caption,
		Tagger: &github.CommitAuthor{
			Date:  &timestamp,
			Name:  commit.Commit.Author.Name,
			Email: commit.Commit.Author.Email,
			Login: commit.Commit.Author.Login,
		},
		Object: &github.GitObject{
			Type: &objType,
			SHA:  &sha,
		},
	}

	if err = gitutil.AddTagToCommit(client, owner, repo, tag); err != nil {
		return fmt.Errorf("error when adding tag to commit %q: %v", fullname, err)
	}

	log.Printf("Added a new version for %q: %q", fullname, caption)
	return nil
}

func fetch(settings *settings.AtcSettings, ghOldContentProviderPtr,
	ghNewContentProviderPtr provider.ContentProvider, fullname string) (string, error) {
	fetchType := detectFetchType(settings.Path)
	var newVersion string
	var oldVersion string
	if fetchType != "" {
		var err error
		af := autoFetchers[fetchType]
		if af == nil {
			log.Printf("using custom fetcher")
			if settings.RegexStr == "" {
				return "", fmt.Errorf("don't have regexstr for not default package manager file %s", fetchType)
			}
			af = &customregex.Fetcher{}
		}

		oldVersion, err = af.GetVersion(ghOldContentProviderPtr, *settings)
		if err != nil && !errors.Is(err, provider.ErrHttpStatusCode) { //ignore http api error
			return "", fmt.Errorf("get prev version error for %q: %w", fullname, err)
		}

		log.Printf("old version %s", oldVersion)
		newVersion, err = af.GetVersion(ghNewContentProviderPtr, *settings)
		if err != nil {
			return "", fmt.Errorf("get new version error for %q: %w", fullname, err)
		}
	} else {
		fetched := false
		for defaultPath, af := range autoFetchers {
			var err error
			oldVersion, err = af.GetVersionUsingDefaultPath(ghOldContentProviderPtr)
			if err != nil && !errors.Is(err, provider.ErrHttpStatusCode) { //ignore http api error
				log.Printf("get prev version error for %q, default path: %s, err: %v", fullname, defaultPath, err)
				continue
			}

			newVersion, err = af.GetVersionUsingDefaultPath(ghNewContentProviderPtr)
			if err == nil {
				fetched = true
				break
			} else {
				return "", fmt.Errorf("autofetcher error for %q: %w", defaultPath, err)
			}
		}

		if !fetched {
			return "", fmt.Errorf("unable to fetch version using known methods")
		}
	}

	if newVersion != oldVersion {
		log.Printf("There is a new version for %q! Old version: %q, new version: %q", fullname, oldVersion, newVersion)
		caption, err := renderTagNameTemplate(settings.Template, newVersion)
		if err != nil {
			return "", fmt.Errorf("error in go templates: %v", err)
		}
		return caption, nil
	}

	return "", nil
}
