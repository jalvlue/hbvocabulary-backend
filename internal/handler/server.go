package handler

import (
	"HBVocabulary/config"
	"HBVocabulary/token"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	db         *gorm.DB
	router     *gin.Engine
	config     *config.Config
	tokenMaker token.JWTMaker
}

func NewServer(config *config.Config, db *gorm.DB) (*Server, error) {
	log.Println(config.JWTSecretKey)
	tokenMaker, err := token.NewJWTMaker(config.JWTSecretKey)
	if err != nil {
		return nil, fmt.Errorf("create token maker failed: %w", err)
	}

	server := &Server{
		db:         db,
		router:     setupRouter(),
		config:     config,
		tokenMaker: *tokenMaker,
	}

	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
