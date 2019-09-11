# Basic-api
Basic-api is a simple framework for beginners knowing how to build a simple Restful API system.
It also provides the ways to interact with MongoDB, SQL server and Redis.
Beside that, it gives users the way to build http client for communicating with other servers, or systems.
## Download
This is a library writen in Go languge, so for download it, you can you go get command.
```bash
go get "github.com/tinwoan-go/basic-api"
```
## Usage
This section will provide you some simple ways and steps for using this library.
### Sample routing
This library provides a simple routing method based on go-chi library.
This sample routing provides users 2 self-built middlewares for set the request header with no-cache, and log out screen every request and response came to the server.
You can initiate the routers as follow.
```go
package main

import "github.com/tinwoan-go/basic-api/handler"

func main() {
	// Initiate routers for application
	routers := handler.NewRouter()
	...
}
```
### Serve HTTP
This library provides a way to serve HTTP a lot easier than normal.
You don't need to create a server yourself and you don't need to handle graceful shutdown on your own.
You can just do it as follow.

(Note: I'll use the routers as above for serve HTTP in this section)
```go
package main

import (
	"github.com/tinwoan-go/basic-api/handler"
	"github.com/tinwoan-go/basic-api/server"
	"time"
)

func main() {
	// Initiate routers for application
	routers := handler.NewRouter()
	// Serve HTTP on address "localhost:3000", with 15 seconds of graceful shutdown time
	if err := server.ServeHTTP(routers, "localhost:3000", 15 * time.Second); err != nil {
		panic(err)
	}
}
```
### Logger
This library provides a realy simple way to log out the terminal the message in states of Warning, Error, Fatal or Information.
The using is as simple as it name.
```go
package main

impot "github.com/tinwoan-go/basic-api/logger"

func main() {
	logger.Info("You did it!")
}
```
This will print out the terminal
```bash
You did it!
```
### Http client
This library provides a package for sending the request from your application to other systems or services in http requests.
It support 4 methods for Post request with JSON format and XML format, Get request with JSON format and XML format too.

For initiating the http client of your application, you can do as below.
```go
package main

import (
	"github.com/tinwoan-go/basic-api/http"
	"time"
)

func main() {
	// For this example I don't have proxy URL so I'll leave as ""
	// And timeout for this client to send request will be 10 seconds.
	proxyURL := ""
	if err := http.NewHTTPClient(proxyURL, 10 * time.Second); err != nil {
		panic(err)
	}
}
```
### Mongo
This library provides a package for wrapping basic methods for interacting with MongoDB named mongo.
This package uses "github.com/globalsign/mgo" library to interact with MongoDB itself.
First thing first, you need to initiate the connection to MongoDB for your application.
```go
package main

import (
	"github.com/tinwoan-go/basic-api/mongo"
	"time"
)

func main() {
	// Connect to you local mongo server, into database name "database"
	// with user "user" and password "password", and the timeout for
	// the connection will be 10 seconds.
	if err := mongo.NewMongoClient("127.0.0.1:27017", "database", "user", "password", 10 * time.Second); err != nil {
		panic(err)
	}
	
	// Do not forget to close the session to release it.
	defer mongo.Close()
}
```
