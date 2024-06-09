package main

import (
	"github.com/sudhir512kj/ecommerce_backend/config"
	"github.com/sudhir512kj/ecommerce_backend/database"
	"github.com/sudhir512kj/ecommerce_backend/server"
)

var conf *config.Config

func main() {
	conf = config.GetConfig()
	db := database.NewPostgresDatabase(conf)
	server.NewEchoServer(conf, db).Start()
}
