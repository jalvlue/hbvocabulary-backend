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

	userAPI := router.Group("/user")
	userAPI.POST("/login", server.loginUser)
	userAPI.POST("/register", server.createUser)
	// auth route /user/info
	userAPI.Group("/info", authMiddleware(server.tokenMaker)).GET("", server.infoUser)
	userAPI.Group("/grade", authMiddleware(server.tokenMaker)).POST("", server.setGradesUser)

	wordAPI := router.Group("/word")
	wordAPI.Use(authMiddleware(server.tokenMaker))
	wordAPI.GET("/roundOne", server.getWordListRoundOne)
	wordAPI.POST("/roundTwo", server.getWordListRoundTwo)
	wordAPI.POST("/getResult", server.getResult)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
