package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)


type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewAccessToken(ctx *gin.Context) {
	var request renewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// verify refresh token
	refreshPayload, err := server.tokenMaker.VerifyToken(request.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	session, err := server.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// check session is blocked
	if session.IsBlocked {
		err := errors.New("blocked user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// check session user name
	if session.UserName != refreshPayload.Username {
		err := errors.New("incorrect session user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// check session refresh token
	if session.RefreshToken != request.RefreshToken {
		err := errors.New("mimatched refresh token")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// check expired 
	if time.Now().After(refreshPayload.ExpiredAt) {
		err := errors.New("refresh token is expired")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(session.UserName, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := renewAccessTokenResponse{
		AccessToken: accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}