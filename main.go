package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/mitchellh/go-homedir"
	"log"
	"os"
	"os/exec"
	"strings"
	"golang.org/x/oauth2"
)

var  CacheDir string

func ListGists(client *github.Client, username string){
	gists, _, err := client.Gists.List(context.Background(), username, nil)
	if err != nil {
		log.Println("Gists.List returned error: %v", err)
	}

	for _, gist := range gists{
		for _, file := range gist.Files{
			line := fmt.Sprintf("%32s	%10s	%10s", gist.GetID(), file.GetFilename(), gist.GetDescription())
			fmt.Println(line)
		}

	}
}

func RunCommand(filename string){
	out, err := exec.Command("sh", filename).Output()
	if err != nil {
		log.Println(err)
	}
	log.Println("Output : ", string(out))
}

func RunGist(client *github.Client, gistID string){
	gist, _, err := client.Gists.Get(context.Background(), gistID)
	if err != nil {
		log.Println("Gists.List returned error: %v", err)
	}

	log.Println(gist.GetDescription())
	for _, file := range gist.Files{
		line := fmt.Sprintf("[RUN] (%s) %s ...", file.GetLanguage(), file.GetFilename())
		fmt.Println(line)
		err := os.Mkdir(CacheDir, os.ModePerm)
		if err != nil{
			log.Println(err)
		}
		cmdStr := file.GetContent()
		if cmdStr != ``{
			tmpFilename := fmt.Sprintf("%s/%s.%s", CacheDir, gistID, file.GetLanguage())
			log.Println(tmpFilename)
			f, err := os.Create(tmpFilename)
			if err!=nil{
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
	CacheDir = userhomedir +  "/.gistcachedir"


	accessToken := flag.Arg(0)
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	cmd := flag.Arg(1)
	switch cmd{
	case `list`:
		username := flag.Arg(2)
		if username == ``{
			log.Println("username not found")
			return
		}
		ListGists(client, username)
	case `run`:
		tmp := flag.Arg(2)
		gistID := strings.Split(tmp, "\t")[0]
		RunGist(client, gistID)
	}

}