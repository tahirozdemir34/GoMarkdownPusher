# ATTENTION: This project is highly experimental

This is a simple web application that allows you to create, update and delete markdown files in your GitHub repositories. 

Used libraries:

* Back-end
  * [go-github](https://github.com/google/go-github)
  * [oauth2](https://github.com/golang/oauth2)
* Front-end
  * [Axios](https://github.com/axios/axios)
  * [Showdown](https://github.com/showdownjs/showdown)
  * [Bootstrap](https://github.com/twbs/bootstrap)

### How to use:
* Get a token from Settings>Developer Settings>Personal access tokens
* create a config.json next to 'md_editor.go' as follow:

```json
{
    "token" : "your token",
    "username" : "your username",
}
```
* Run the program as follow:
```bash
go run md_editor.go
```
* go to localhost:8000

*This readme file created by using GoMarkdownPusher*