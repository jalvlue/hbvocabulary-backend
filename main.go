package main

import (
	"HBVocabulary/config"
	"HBVocabulary/internal/handler"
	"HBVocabulary/internal/model"
	"log"
	"os"
)

func init() {
	config.InitConfig()
	// model.InitDB(config.Conf.DBSource)
	model.InitDB(os.Getenv("DB_SOURCE"))
}

func main() {
	server, err := handler.NewServer(config.Conf, model.NewStore(model.DB))
	if err != nil {
		log.Fatal("cannot create http server: ", err)
	}

	err = server.Start(config.Conf.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
