package main

import (
	// Native packages
	"net/http"
	"time"

	// Third parties
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/go-chi/render"

	// Internal packages
	"github.com/tinwoan-go/basic-api/handler"
	"github.com/tinwoan-go/basic-api/mongo"
	"github.com/tinwoan-go/basic-api/redis"
	"github.com/tinwoan-go/basic-api/server"
	"github.com/tinwoan-go/basic-api/sql"
	"github.com/tinwoan-go/basic-api/tlog"
)

var (
	log tlog.Logger
)

func init() {
	log = tlog.WithPrefix("main")
}

func main() {
	// Connect Mongo
	mongoCfg := mongo.Configs{
		Addresses: "127.0.0.1:27017",
		Database:  "Testing",
		Username:  "tin",
		Password:  "tinwoan",
		Timeout:   5 * time.Second,
	}
	if err := mongo.NewMongoClient(mongoCfg); err != nil {
		panic(err)
	}
	defer mongo.Close()

	// Connect SQL
	sqlCfg := sql.Configs{
		Driver:   "mssql",
		Host:     "127.0.0.1",
		Port:     "1433",
		Username: "SA",
		Password: "TinWoan1234",
		Database: "information",
	}
	if err := sql.NewSQL(sqlCfg); err != nil {
		panic(err)
	}
	defer func() {
		if err := sql.Close(); err != nil {
			log.Warnf("Error close SQL connection: %v", err)
		}
	}()

	// Connect Redis
	redisCfg := redis.Configs{
		Addresses: []string{"127.0.0.1:6379"},
		Master:    "",
		Password:  "",
	}
	redis.NewRedisClient(redisCfg)
	defer func() {
		if err := redis.Close(); err != nil {
			log.Warnf("Error close Redis connection: %v", err)
		}
	}()

	r := handler.NewRouter()
	if err := server.ServeHTTP(r, "localhost:8000", 15*time.Second); err != nil {
		panic(err)
	}
}

func status() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		render.Status(r, http.StatusOK)
		render.JSON(w, r, struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}{
			Name: "Tin",
			Age:  24,
		})
	})
}
