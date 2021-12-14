package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v41/github"
	"golang.org/x/oauth2"
)

func main() {
	token := os.Getenv("TEST_GITHUB_TOKEN")
	if token == "" {
		println("BANG!!!")
		os.Exit(1)
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	var tag = "v0.2.0"

	var body *string
	chlog, err := changelog(ctx, client, "wesleimp-deprecated", "labs", "v0.1.0", tag)
	if err != nil {
		println("FAILED TO GET CHANGELOG", err.Error())
		body = github.String("generated from github client")
	} else {
		body = github.String(chlog)
	}

	r := &github.RepositoryRelease{TagName: github.String(tag), Name: github.String("sldkfjas"), Body: body}
	release, _, err := client.Repositories.GetReleaseByTag(
		ctx, "wesleimp-deprecated", "labs", tag,
	)
	if err != nil {
		println("CREATING RELEASE")
		release, _, err = client.Repositories.CreateRelease(ctx, "wesleimp-deprecated", "labs", r)
		if err != nil {
			println("ERRO ON CREATE", err.Error())
			return
		}
	} else {
		println("UPDATING RELEASE")
		release, _, err = client.Repositories.EditRelease(ctx, "wesleimp-deprecated", "labs", release.GetID(), r)
		if err != nil {
			println("ERRO ON EDIT", err.Error())
			return
		}
	}
	println("ID:", release.GetID(), "URL:", release.GetURL())
}

func changelog(ctx context.Context, client *github.Client, owner, repo, prev, current string) (string, error) {
	var log []string

	opts := &github.ListOptions{PerPage: 100}
	for {
		result, resp, err := client.Repositories.CompareCommits(ctx, owner, repo, prev, current, opts)
		if err != nil {
			return "", err
		}
		for _, commit := range result.Commits {
			log = append(log, fmt.Sprintf(
				"%s: %s",
				commit.GetSHA(),
				strings.Split(commit.Commit.GetMessage(), "\n")[0],
			))
		}
		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return strings.Join(log, "\n"), nil
}
