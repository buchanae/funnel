package main

import (
  "context"
  _ "github.com/kr/pretty"
	"fmt"
  "strings"
	"gopkg.in/src-d/go-git.v4"
	. "gopkg.in/src-d/go-git.v4/_examples"
  "io"
  "golang.org/x/oauth2"
  "github.com/google/go-github/github"
  "os"
)

func main() {
  tok, ok := os.LookupEnv("GITHUB_TOKEN")
  if !ok {
    panic("GITHUB_TOKEN required")
  }

  ctx := context.Background()
  ts := oauth2.StaticTokenSource(
    &oauth2.Token{AccessToken: tok},
  )
  tc := oauth2.NewClient(ctx, ts)
  cl := github.NewClient(tc)

  r, err := git.PlainOpen(".")
	CheckIfError(err)

	start, err := r.ResolveRevision("master")
	CheckIfError(err)

  stop, err := r.ResolveRevision("0.5.0")
	CheckIfError(err)

	cIter, err := r.Log(&git.LogOptions{From: *start})
	CheckIfError(err)

  for {
    c, err := cIter.Next()
    if err == io.EOF {
      break
    }
	  CheckIfError(err)
    if c.Hash == *stop {
      break
    }

    sha := c.Hash.String()
		fmt.Println(sha[:6], strings.TrimRight(c.Message, "\n"))

    res, _, err := cl.Search.Issues(ctx, sha + " type:pr repo:ohsu-comp-bio/funnel", nil)
    if err != nil {
      panic(err)
    }

    //pretty.Println(res)
    if *res.Total != 1 {
      panic("expected one issue/pr")
    }
    pr := res.Issues[0]
    fmt.Printf("PR #%d, %s\n", *pr.Number, *pr.HTMLURL)
  }
}
