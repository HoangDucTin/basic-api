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
	// Connect to your local mongo server, into database name "database"
	// with user "user" and password "password", and the timeout for
	// the connection will be 10 seconds.
	if err := mongo.NewMongoClient("127.0.0.1:27017", "database", "user", "password", 10 * time.Second); err != nil {
		panic(err)
	}
	
	// Do not forget to close the session to release it.
	defer mongo.Close()
}
```
This package provides methods to find a record (Find), find all records (FindAll), insert a record (Insert), insert many records (InsertAll), remove a latest record (Remove), remove all records (RemoveAll), update a latest record (Update) and update all records (UpdateAll).
### Redis
This library provides a package for connecting to Redis server, based on "github.com/go-redis/redis".
First thing first, you need to initiate the connection to Redis server for your application.
```go
package main

import (
	"github.com/tinwoan-go/basic-api/logger"
	"github.com/tinwoan-go/basic-api/redis"
)

func main() {
	// Connect to your local redis server
	// with user "user" and password "password".
	// If you don't use tunnels on your
	// Redis server, leave the master name empty.
	masterName := ""
	addrs := []string{"127.0.0.1:6379"}
	if err := redis.NewRedisClient("", "user", "password", addrs); err != nil {
		panic(err)
	}
	
	// Do not forget to close the connection
	// to redis server after using it.
	defer func(){
		if err := redis.Close(); err != nil {
			logger.Warn("Can not close connection to Redis server, error: %v", err)
	}()
}
```
This package provides 2 simple method to get data from Redis server (Get) and set a value to redis server with a specified key (Set).
### SQL
This library provides a package named "sql" for connecting and interacting with SQL server.
(Caution: Because this package is built on the purpose of making things generic, I've use the json format in some cases and I'm trying to implement it to a better phase.)
First thing first, you have to initiate the connection to SQL server if you want to use it.
(Notice: In my library, I used driver "github.com/denisenkom/go-mssqldb" for connect with local Microsoft SQL server.)
```go
package main

import (
	"github.com/tinwoan-go/basic-api/logger"
	"github.com/tinwoan-go/basic-api/sql"
)

func main() {
	// Connect to your local Microsoft SQL server
	// with user "user" and password "password",
	// and connect directly to database "database".
	if err := sql.NewSql("localhost", "user", "password", "database"); err != nil {
		panic(err)
	}

	// Do not forget to close the connection
	// to your SQL server after using it.
	defer func() {
		if err := sql.Close(); err != nil {
			logger.Error("Cannot close connection to SQL server, error: %v", err)
		}
	}()
}
```
This package provides methods to find one (Find), find all (FindAll), find with a limit (FindWithLimit), insert one (Insert), insert many (InsertMany), delete (Delete), update (Update).

Caution:
	To use this package, your structs must have json tag.
	Example:
```go
type Test struct {
	Message string `json:"message"`
}
```
	And the json tags should have the same names with the name of the columns in your SQL database.
	And the fields in the structs should have the same order with the order of the columns in your SQL database.
	
	Most of the functions can use parameters either in structs or in map[string]interface{}.
	(For now the generic implementation for SQL server has not been introduced yet, so I've used the json tag for this purpose.
	Sorry for the inconvinience and I'll improve it in the future.)
