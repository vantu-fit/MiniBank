package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpriedAt time.Time `json:"access_token_expried_at"`
}

func (server *Server) renewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	// verify refresh token
	refreshPayload, err := server.maker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}
	// kiem tra session
	session, err := server.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
	}
	// check block session
	if session.IsBlocked {
		err := fmt.Errorf("block session")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	if session.Username != refreshPayload.Username {
		err := fmt.Errorf("incorrect session")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("block incrrect session user")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}
	// check expried session
	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("sesion expried")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}
	// parser time suration
	accesstokenDuration, err := time.ParseDuration(server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	// create new accesstoken
	accesstoken, accessPayload, err := server.maker.CreateToken(session.Username,refreshPayload.Role ,  accesstokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	// create response
	response := renewAccessTokenResponse{
		AccessToken:          accesstoken,
		AccessTokenExpriedAt: accessPayload.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, response)
	return

}
