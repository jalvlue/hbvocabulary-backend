package handler

import (
	"HBVocabulary/config"
	"HBVocabulary/internal/model"
	"HBVocabulary/token"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store      *model.Store
	router     *gin.Engine
	config     *config.Config
	tokenMaker token.JWTMaker
}

func NewServer(config *config.Config, store *model.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.JWTSecretKey)
	if err != nil {
		return nil, fmt.Errorf("create token maker failed: %w", err)
	}

	server := &Server{
		store:      store,
		config:     config,
		tokenMaker: *tokenMaker,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	// TODO: setup router
	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello hbvocabulary")
	})

	v1 := router.Group("/user")
	v1.POST("/login", server.loginUser)
	v1.POST("/register", server.createUser)
	// auth route /user/info
	v1.Group("/info", authMiddleware(server.tokenMaker)).GET("", server.infoUser)

	// v2 := router.Group("/word")

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
