package main

import (
	"context"
	"os"

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

	var tag = "v0.1.0"

	var body *string
	notes, _, err := client.Repositories.GenerateReleaseNotes(ctx, "wesleimp-deprecated", "labs", &github.GenerateNotesOptions{
		TagName:         tag,
		PreviousTagName: github.String(tag),
	})
	if err != nil {
		println("FAILED TO GET CHANGELOG", err.Error())
		body = github.String("generated from github client")
	} else {
		body = github.String(notes.Body)
	}

	r := &github.RepositoryRelease{TagName: github.String(tag), Name: github.String("experimental"), Body: body}
	release, _, err := client.Repositories.GetReleaseByTag(
		ctx, "wesleimp-deprecated", "labs", "v0.1.0",
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
