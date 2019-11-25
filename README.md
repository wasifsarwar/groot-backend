# Groot

Groot is a web application that connects to a Slack workspace and displays a list of the workspace users and their profile information

### Tech

Groot uses a number of open source projects to work properly:

* [React.js] - Interactive UI that displays the result
* [Golang] - Primarily written in go
* [nlopes/slack] - Go library that supports most api.slack.com REST calls.
* [Mongo] - Existing users are persisted in a mongodb
* [mongo-go-driver] - The MongoDB supported driver for Go

### Local Installation

You'll need a mongodb connection string and a slack OAUTH key to start groot

```sh
$ cd groot-backend
$ export MONGO_STRING="YOUR_MONGO_TOKEN"
$ export SLACK_TOKEN="YOUR_SLACK_TOKEN"
$ go run server/main.go
```

Head over to [groot frontend] to install the front end UI and get started with the web app



[//]: # (These are reference links used in the body of this note and get stripped out when the markdown processor does its job. There is no need to format nicely because it shouldn't be seen. Thanks SO - http://stackoverflow.com/questions/4823468/store-comments-in-markdown-syntax)


   [React.js]:<https://reactjs.org/>
   [Golang]:<https://golang.org/>
   [nlopes/slack]:<https://github.com/nlopes/slack>
   [Mongo]:<https://www.mongodb.com/>
   [mongo-go-driver]:<https://github.com/mongodb/mongo-go-driver>
   [groot frontend]:<https://github.com/wasifsarwar/groot-ui>
   
