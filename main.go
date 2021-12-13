package main

import (
	"context"

	"github.com/google/go-github/v41/github"
)

func main() {
	client := github.NewClient(nil)

	r := &github.RepositoryRelease{TagName: github.String("v0.1.0"), Name: github.String("teeest"), Body: github.String("generated from github client")}

	release, _, err := client.Repositories.CreateRelease(context.Background(), "wesleimp-deprecated", "todo-gtk", r)
	if err != nil {
		println("Ooooops", err.Error())
		return
	}
	println(release.GetID(), release.GetURL())
}
