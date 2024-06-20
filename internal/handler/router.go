package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()

	// TODO: setup router
	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello hbvocabulary")
	})

	return router
}
