package handler

import (
	"HBVocabulary/common"
	"HBVocabulary/internal/model"
	"HBVocabulary/token"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type createUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type createUserResponse struct {
	Code      string    `json:"code"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse(err))
		return
	}

	user := model.User{
		Username:  req.Username,
		Password:  req.Password,
		TestCount: 0,
		MaxScore:  0,
	}
	err := server.store.CreateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, createUserResponse{
		Code:      "success",
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	})
}

type loginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginUserResponse struct {
	Code                string    `json:"code"`
	Username            string    `json:"username"`
	AccessToken         string    `json:"access_token"`
	AccessTokenExpireAt time.Time `json:"access_token_expire_at"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse(err))
		return
	}

	user, err := server.store.GetUserByUsername(req.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse(err))
		return
	}
	if user.Password != req.Password {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "invalid password",
		})
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, loginUserResponse{
		Code:                "success",
		Username:            user.Username,
		AccessToken:         accessToken,
		AccessTokenExpireAt: accessPayload.ExpiredAt,
	})
}

type InfoUserResponse struct {
	Code      string    `json:"code"`
	Username  string    `json:"username"`
	TestCount int       `json:"test_count"`
	MaxScore  int       `json:"max_score"`
	FourScore int       `json:"four_score"`
	SixScore  int       `json:"six_score"`
	CreatedAt time.Time `json:"created_at"`
}

func (server *Server) infoUser(ctx *gin.Context) {
	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	user, err := server.store.GetUserByUsername(payload.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, InfoUserResponse{
		Code:      "success",
		Username:  user.Username,
		TestCount: user.TestCount,
		MaxScore:  user.MaxScore,
		FourScore: user.FourScore,
		SixScore:  user.SixScore,
		CreatedAt: user.CreatedAt,
	})
}

type SetGradesRequest struct {
	FourGrade int `json:"four"`
	SixGrade  int `json:"six"`
}

func (server *Server) setGradesUser(ctx *gin.Context) {
	var reqSetGrade SetGradesRequest
	if err := ctx.ShouldBindJSON(&reqSetGrade); err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse(err))
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	user, err := server.store.GetUserByUsername(payload.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse(err))
		return
	}

	err = server.store.SetGrades(user, reqSetGrade.FourGrade, reqSetGrade.SixGrade)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": "success",
	})
}
