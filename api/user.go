package api

import (
	"net/http"
	"time"

	db "github.com/amallick86/Go_bank/db/sqlc"
	"github.com/amallick86/Go_bank/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserRequest struct{
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"fullname" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Citizenship string `json:"citizenship" binding:"required"`
}
type cretedUserResponse struct{
	Username          string    `json:"username"`
	FullName          string    `json:"fullName"`
	Email             string    `json:"email"`
	Citizenship     string    `json:"citizenship"`
	PasswordChangedAt time.Time `json:"passwordChangedAt"`
	CreatedAt         time.Time `json:"createdAt"`
}
func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	arg := db.CreateUserParams{
		Username:         req.Username,
		HashedPassword:      hashedPassword,
		FullName: req.FullName,
		Email:  req.Email,
		Citizenship: req.Citizenship,
		
	}
	
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation" :
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := cretedUserResponse{
		Username: user.Username,
		FullName: user.FullName,
		Email:  user.Email,
		Citizenship: user.Citizenship,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt: user.CreatedAt,

	}
	ctx.JSON(http.StatusOK, rsp)
}
