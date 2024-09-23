package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/longln/simplebank/db/sqlc"
	"github.com/longln/simplebank/utils"
)



type createUserRequest struct {
	UserName    string `json:"user_name" binding:"required,alphanum"`
	Password 	string `json:"password" binding:"required,min=6"`
	FullName	string `json:"full_name" binding:"required"`
	Email		string `json:"email" binding:"required,email"`
}

type createUserResponse struct {
	UserName string `json:"user_name"`
	FullName string `json:"full_name"`
	Email string `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var request createUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedpassword, err := utils.HashPassword(request.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	arg := db.CreateUserParams{
		UserName: request.UserName,
		HashedPassword: hashedpassword,
		FullName: request.FullName,
		Email: request.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	response := createUserResponse{
		UserName: user.UserName,
		FullName: user.FullName,
		Email: user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt: user.CreatedAt,

	}
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, response)
}
