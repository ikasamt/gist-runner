package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/google/go-github/github"
	homedir "github.com/mitchellh/go-homedir"
	"golang.org/x/oauth2"
)

var CacheDir string

func ListGists(client *github.Client, username string) {
	gists, _, err := client.Gists.List(context.Background(), username, nil)
	if err != nil {
		log.Println("Gists.List returned error: %v", err)
	}

	for _, gist := range gists {
		for _, file := range gist.Files {
			line := fmt.Sprintf("%32s	%10s	%10s", gist.GetID(), file.GetFilename(), gist.GetDescription())
			fmt.Println(line)
		}

	}
}

func RunCommand(filename string) {
	out, err := exec.Command("sh", filename).Output()
	if err != nil {
		log.Println(err)
	}
	log.Println("Output : ", string(out))
}

func RunGist(client *github.Client, gistID string) {
	gist, _, err := client.Gists.Get(context.Background(), gistID)
	if err != nil {
		log.Println("Gists.List returned error: %v", err)
	}

	log.Println(gist.GetDescription())
	for _, file := range gist.Files {
		line := fmt.Sprintf("[RUN] (%s) %s ...", file.GetLanguage(), file.GetFilename())
		fmt.Println(line)
		err := os.Mkdir(CacheDir, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
		cmdStr := file.GetContent()
		if cmdStr != `` {
			tmpFilename := fmt.Sprintf("%s/%s.%s", CacheDir, gistID, file.GetLanguage())
			log.Println(tmpFilename)
			f, err := os.Create(tmpFilename)
			if err != nil {
				log.Println(err)
				return
			}
			defer f.Close()
			f.WriteString(cmdStr)
			RunCommand(tmpFilename)
			//os.Remove(tmpFilename)
		}
	}
}

func main() {
	flag.Parse()

	userhomedir, _ := homedir.Dir()
	CacheDir = userhomedir + "/.gist-runner/cachedir"

	tokenPath := userhomedir + "/.gist-runner/token"
	accessTokenBytes, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		log.Println(`~/.gist-runner/token not found`)
		os.Exit(1)
	}
	accessToken := string(accessTokenBytes)

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	cmd := flag.Arg(0)
	switch cmd {
	case `run`:
		tmp := flag.Arg(1)
		gistID := strings.Split(tmp, "\t")[0]
		RunGist(client, gistID)
	default:
		username := cmd
		if username == `` {
			log.Println("username not found")
			return
		}
		ListGists(client, username)
	}
}
