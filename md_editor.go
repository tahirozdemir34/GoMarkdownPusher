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
	Operation string `json:"operation"`
	Content   string `json:"content"`
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
		log.Println(t.Operation)
		log.Println(t.Content)
		if t.Operation == "Create" {
			createFile(t.Content)
		} else {
			updateFile(t.Content, "Jan _2 15:04:05") //TODO: File name is placeholder. Please be sure that it exists
		}

	} else {
		fmt.Println("Unknown HTTP " + r.Method + "  Method")
	}
}
func createFile(content string) {
	ctx := context.Background() //TODO: Make ctx, ts, tc and client global
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	fileContent := []byte(content)
	opts := &github.RepositoryContentFileOptions{
		Message: github.String(time.Now().String()),
		Content: fileContent,
		Branch:  github.String("master"),
		//Committer: &github.CommitAuthor{Name: github.String("FirstName LastName"), Email: github.String("user@example.com")},
	}
	client.Repositories.CreateFile(ctx, config.Username, config.Repo, time.Stamp+".md", opts)
}

func updateFile(content, filename string) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	optsForGet := &github.RepositoryContentGetOptions{
		Ref: "master",
	}
	file, _, _, err := client.Repositories.GetContents(ctx, config.Username, config.Repo, filename, optsForGet)
	if err != nil {
		panic(err)
	}
	fmt.Println(file.SHA)
	fileContent := []byte(content)
	optsForUpdate := &github.RepositoryContentFileOptions{
		Message: github.String(time.Now().String()),
		Content: fileContent,
		Branch:  github.String("master"),
		SHA:     file.SHA,
		//Committer: &github.CommitAuthor{Name: github.String("FirstName LastName"), Email: github.String("user@example.com")},
	}
	client.Repositories.UpdateFile(ctx, config.Username, config.Repo, filename, optsForUpdate)
}

func main() {
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
	fmt.Println("Listenning and serving on port 8000. Please visit 127.0.0.1:8000...")
	http.ListenAndServe(":8000", nil) // setting listening port

}
