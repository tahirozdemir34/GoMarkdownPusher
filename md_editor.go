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
	"path/filepath"
	"time"

	"encoding/base64"

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
type Markdown struct {
	Name string
}

var config Config
var activeRepo = "..."
var activeFile = "..."

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func md_editor(w http.ResponseWriter, r *http.Request) {
	fmt.Println("A client just sent a request")
	setupResponse(&w, r)
	if r.Method == "GET" {
		// GET
		println("GET")
		t, _ := template.ParseFiles("editor.html")
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		// POST
		ctx := context.Background() //TODO: Make ctx, ts, tc and client global
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: config.Token},
		)
		tc := oauth2.NewClient(ctx, ts)
		client := github.NewClient(tc)
		println("POST")
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
		} else if t.Operation == "listFiles" {
			_, content, _, err := client.Repositories.GetContents(ctx, config.Username, t.Content, "", nil)
			if err != nil {
				panic(err)
			}
			activeRepo = t.Content
			list := []Markdown{}
			list = append(list, Markdown{"..."})
			for _, element := range content {
				var extension = filepath.Ext(*element.Name)
				if extension == ".md" || extension == ".MD" {
					var temp Markdown
					temp.Name = *element.Name
					list = append(list, temp)
				}
			}

			data, err := json.Marshal(list)
			if err != nil {
				panic(err)
			}
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		} else if t.Operation == "listRepos" {
			repos, _, _ := client.Repositories.List(ctx, "", nil)
			list := []Markdown{}
			for _, element := range repos {
				var temp Markdown
				temp.Name = *element.Name
				list = append(list, temp)
			}

			data, err := json.Marshal(list)
			if err != nil {
				panic(err)
			}
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		} else if t.Operation == "getFileContent" {
			content, _, _, err := client.Repositories.GetContents(ctx, config.Username, activeRepo, t.Content, nil)
			if err != nil {
				panic(err)
			}
			activeFile = t.Content
			decoded, err := base64.StdEncoding.DecodeString(*content.Content)
			println(string(decoded))
			if err != nil {
				panic(err)
			}
			data, err := json.Marshal(string(decoded))
			if err != nil {
				panic(err)
			}
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		} else {
			updateFile(t.Content, activeRepo, activeFile) //TODO: File name is placeholder. Please be sure that it exists
		}

	} else if (*r).Method == "OPTIONS" {
		return

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
	client.Repositories.CreateFile(ctx, config.Username, config.Repo, time.Now().String()+".md", opts)

}

func updateFile(content, reponame string, filename string) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	optsForGet := &github.RepositoryContentGetOptions{
		Ref: "master",
	}
	file, _, _, err := client.Repositories.GetContents(ctx, config.Username, reponame, filename, optsForGet)
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
	_, _, err = client.Repositories.UpdateFile(ctx, config.Username, config.Repo, filename, optsForUpdate)
	if err != nil {
		panic(err)
	}
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
