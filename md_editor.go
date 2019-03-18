package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type markdown struct {
	Text string `json:"text"`
}

type Config struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Repo     string `json:"repo"`
}

var config Config

func md_editor(w http.ResponseWriter, r *http.Request) {
	fmt.Println("A client just sent a request")
	if r.Method == "GET" {
		// GET
		t, _ := template.ParseFiles("editor.html")
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		// POST
		decoder := json.NewDecoder(r.Body)
		var t markdown
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}
		log.Println(t.Text)

		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: config.Token},
		)
		tc := oauth2.NewClient(ctx, ts)

		client := github.NewClient(tc)

		// list all repositories for the authenticated user
		/*	repos, _, _ := client.Repositories.List(ctx, "", nil)
			fmt.Println(repos)*/
		fileContent := []byte(t.Text)

		// Note: the file needs to be absent from the repository as you are not
		// specifying a SHA reference here.

		opts := &github.RepositoryContentFileOptions{
			Message: github.String(time.Now().String()),
			Content: fileContent,
			Branch:  github.String("master"),
			//Committer: &github.CommitAuthor{Name: github.String("FirstName LastName"), Email: github.String("user@example.com")},
		}

		//user, _, _ :=

		client.Repositories.CreateFile(ctx, config.Username, config.Repo, time.Stamp+".md", opts)

	} else {
		fmt.Println("Unknown HTTP " + r.Method + "  Method")
	}
}

func main() {

	/*repo := &github.Repository{
		Name:    github.String("foo"),
		Private: github.Bool(true),
	}
	client.Repositories.Create(ctx, "", repo)*/
	jsonFile, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &config)

	fmt.Println(config.Username)
	fmt.Println(config.Token)
	fmt.Println(config.Repo)

	http.HandleFunc("/", md_editor)
	http.ListenAndServe(":8000", nil) // setting listening port
}
