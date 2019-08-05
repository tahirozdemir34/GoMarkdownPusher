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
	"sort"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Request struct {
	Operation string `json:"operation"`
	Content   string `json:"content"`
}

type Config struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}
type File struct {
	Name string
	Type string
	Path string
}

var config Config
var activeRepo = "..."
var activeDir = "..."
var ctx context.Context
var ts oauth2.TokenSource
var tc *http.Client
var client *github.Client

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
		println("POST")
		decoder := json.NewDecoder(r.Body)
		var req Request
		err := decoder.Decode(&req)
		if err != nil {
			panic(err)
		}
		log.Println(req.Operation)
		log.Println(req.Content)
		operation := strings.Split(req.Operation, "//")
		if operation[0] == "Create" {
			go createFile(activeRepo, operation[1], req.Content)
		} else if operation[0] == "Delete" {
			go deleteFile(req.Content, activeRepo, activeDir)
		} else if operation[0] == "Update" {
			go updateFile(req.Content, activeRepo, activeDir)
		} else if operation[0] == "listFilesFromDir" {
			println("fromdir")
			_, content, _, err := client.Repositories.GetContents(ctx, config.Username, activeRepo, req.Content, nil)
			if err != nil {
				panic(err)
			}
			activeDir = req.Content
			list := []File{}
			for _, element := range content {

				var extension = filepath.Ext(*element.Name)
				if strings.EqualFold(extension, ".md") || element.GetType() == "dir" || strings.EqualFold(extension, "") {
					var temp File
					temp.Name = element.GetName()
					temp.Type = element.GetType()
					temp.Path = element.GetPath()
					list = append(list, temp)
				}
			}
			sort.Slice(list, func(i, j int) bool {
				if list[i].Type < list[j].Type {
					return true
				}
				if list[i].Type > list[j].Type {
					return false
				}
				return strings.ToLower(list[i].Name) < strings.ToLower(list[j].Name)
			})
			data, err := json.Marshal(list)
			if err != nil {
				panic(err)
			}
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		} else if operation[0] == "listFilesFromRepo" {
			_, content, _, err := client.Repositories.GetContents(ctx, config.Username, req.Content, "", nil)
			if err != nil {
				panic(err)
			}
			activeRepo = req.Content
			list := []File{}
			for _, element := range content {

				var extension = filepath.Ext(*element.Name)
				if strings.EqualFold(extension, ".md") || element.GetType() == "dir" || strings.EqualFold(extension, "") {
					var temp File
					temp.Name = element.GetName()
					temp.Type = element.GetType()
					temp.Path = element.GetPath()
					list = append(list, temp)
				}
			}
			sort.Slice(list, func(i, j int) bool {
				if list[i].Type < list[j].Type {
					return true
				}
				if list[i].Type > list[j].Type {
					return false
				}
				return strings.ToLower(list[i].Name) < strings.ToLower(list[j].Name)
			})
			data, err := json.Marshal(list)
			if err != nil {
				panic(err)
			}
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		} else if operation[0] == "listRepos" {
			opt := &github.RepositoryListOptions{Type: "owner"}
			repos, _, err := client.Repositories.List(ctx, "", opt)
			list := []File{}
			for _, element := range repos {
				var temp File
				temp.Name = element.GetName()
				temp.Path = element.GetCloneURL()
				temp.Type = "Repo"
				list = append(list, temp)
			}
			sort.Slice(list, func(i, j int) bool {
				if list[i].Type < list[j].Type {
					return true
				}
				if list[i].Type > list[j].Type {
					return false
				}
				return strings.ToLower(list[i].Name) < strings.ToLower(list[j].Name)
			})
			data, err := json.Marshal(list)
			if err != nil {
				panic(err)
			}
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		} else if operation[0] == "getFileContent" {
			content, _, _, err := client.Repositories.GetContents(ctx, config.Username, activeRepo, req.Content, nil)
			if err != nil {
				panic(err)
			}

			decoded, err := content.GetContent()
			println(string(decoded))
			if err != nil {
				panic(err)
			}

			data, err := json.Marshal(decoded)
			if err != nil {
				panic(err)
			}
			activeDir = req.Content
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		}
	} else if (*r).Method == "OPTIONS" {
		return

	} else {
		fmt.Println("Unknown HTTP " + r.Method + "  Method")
	}
}
func createFile(reponame string, filename string, content string) {
	println("CREATE")

	fileContent := []byte(content)
	opts := &github.RepositoryContentFileOptions{
		Message: github.String(filename + " created."),
		Content: fileContent,
		Branch:  github.String("master"),
		//Committer: &github.CommitAuthor{Name: github.String("FirstName LastName"), Email: github.String("user@example.com")},
	}
	_, _, err := client.Repositories.CreateFile(ctx, config.Username, reponame, filename, opts)
	if err != nil {
		panic(err)
	}

}

func deleteFile(content, reponame string, filename string) {
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
		Message: github.String(filename + " deleted."),
		Content: fileContent,
		Branch:  github.String("master"),
		SHA:     file.SHA,
		//Committer: &github.CommitAuthor{Name: github.String("FirstName LastName"), Email: github.String("user@example.com")},
	}
	_, _, err = client.Repositories.DeleteFile(ctx, config.Username, reponame, filename, optsForUpdate)
	if err != nil {
		panic(err)
	}
}

func updateFile(content, reponame string, filename string) {
	optsForGet := &github.RepositoryContentGetOptions{
		Ref: "master",
	}
	println(reponame)
	println(filename)
	file, _, _, err := client.Repositories.GetContents(ctx, config.Username, reponame, filename, optsForGet)
	if err != nil {
		panic(err)
	}

	fmt.Println(file.SHA)
	fileContent := []byte(content)
	optsForUpdate := &github.RepositoryContentFileOptions{
		Message: github.String(filename + " updated."),
		Content: fileContent,
		Branch:  github.String("master"),
		SHA:     file.SHA,
		//Committer: &github.CommitAuthor{Name: github.String("FirstName LastName"), Email: github.String("user@example.com")},
	}
	_, _, err = client.Repositories.UpdateFile(ctx, config.Username, reponame, filename, optsForUpdate)
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

	ctx = context.Background()
	ts = oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Token},
	)
	tc = oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)

	http.HandleFunc("/", md_editor)
	fmt.Println("Listenning and serving on port 8000. Please visit 127.0.0.1:8000...")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":8000", nil) // setting listening port

}
