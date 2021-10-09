# Appointy_Task

User Creation and Image Posting API made in GoLang.
 


### _Recruitment task for Appointy_
[Task can be found at this link](https://docs.google.com/document/d/1sFhVumoczf_PmaL_R__Rm9AHqaHsUWgj1x9YcQP6Is4/preview?pru=AAABfIQl1lQ*_YCVPzug85vn7PP9nddWJQ)

## Installation

Use the [GIT](https://git-scm.com/downloads) / Github Desktop to clone this API.

```bash
git clone https://github.com/sudhu-1003/Appointy_Task
```

## Pre-requisites

Initialize the project by typing in the following command in the project directory
```bash
go mod init Appointy_Task
```

These packages must be downloaded/installed after initialization for this code to run successfully
```bash
go get go.mongodb.org/mongo-driver/bson
go get go.mongodb.org/mongo-driver/bson/primitive
go get go.mongodb.org/mongo-driver/mongo
go get go.mongodb.org/mongo-driver/mongo/options
go get golang.org/x/crypto/bcrypt
```

## Usage
Make sure you have MongoDB servers running.

```Go
go run main.go
```



## Comments

Third-Party Packages were not allowed thus not using Gorilla MUX or httprouter. Designed a custom router which can be accessed in /router/router.go
Passwords are stored securely and are encrypted using the [bcrypt libarary](https://pkg.go.dev/golang.org/x/crypto/bcrypt)

## API Commands

- Create an User - POST '/users' with raw json request body. <br />
Sample request body:
{
     "name":"SecretUser",
     "email":"su@gmail.com",
     "password":"secret"
}

- Get a user using id - GET '/users/<id here>'

- Create a Post - POST '/posts' with raw json request body. <br />
Sample request body:
{
	"userID":"6161579e32c758215fc5a575",
	"Caption":"Hello World",
	"ImageURL":"https://cdn.pixabay.com/photo/2015/04/19/08/32/marguerite-729510__480.jpg"
}

- Get a post using id - GET '/posts/<id here>'

- List all posts of a user - GET '/posts/user/<id here>'
