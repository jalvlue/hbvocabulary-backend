package main

import (
	"HBVocabulary/config"
	"HBVocabulary/internal/handler"
	"HBVocabulary/internal/model"
	"log"
)

func init() {
	var err error
	config.Conf, err = config.LoadConfig("./etc/.env")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	model.DB, err = model.InitDB()
	if err != nil {
		log.Fatal("cannot connect to DB: ", err)
	}
}

func main() {
	server, err := handler.NewServer(config.Conf, model.DB)
	if err != nil {
		log.Fatal("cannot create http server: ", err)
	}

	err = server.Start(config.Conf.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
