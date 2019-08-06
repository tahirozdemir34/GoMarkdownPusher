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
* Create "config.json" next to 'md_editor.go' as follow:

```json
{
    "token" : "your token",
    "username" : "your username",
}
```
* Alternatively, define environment variables "GITHUB_USERNAME" and "GITHUB_TOKEN" instead of creating "config.json"
* Run the program as follow:
```bash
go run md_editor.go
```
* go to localhost:8000

### With Docker:
* Run the following command:
```bash
docker run -p 8000:8000 -e GITHUB_USERNAME='username' \
-e GITHUB_TOKEN='token' \
--name go-mdp tahirozdemir34/go-markdown-pusher
```

You can watch a sample usage here:
[![Youtube](https://i.ibb.co/8XmRRS0/image.png)](https://youtu.be/EwRYA8RIWLo)

*This readme file is created by using GoMarkdownPusher*
