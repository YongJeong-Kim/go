package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	db "github.com/yongjeong-kim/go/gogin/db/sqlc"
	"github.com/yongjeong-kim/go/gogin/util"
	"log"
	"net/http"
	"time"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type createUserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	password, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: password,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	err = server.store.CreateUser(ctx, arg)
	if err != nil {
		/*if pgErr, ok := err.(*pg.Error); ok {
			log.Println(pgErr.Code.Name())
		}*/
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			log.Println(mysqlErr.Number)
			log.Println(mysqlErr.Error())
			log.Println(mysqlErr.Message)
			switch mysqlErr.Number {
			case 1062:
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//user, err := server.store.GetUser(ctx, req.Username)
	//if err != nil {
	//	ctx.JSON(http.StatusNotFound, errorResponse(err))
	//	return
	//}

	res := createUserResponse{
		Username: arg.Username,
		FullName: arg.FullName,
		Email:    arg.Email,
		//PasswordChangedAt: arg.PasswordChangedAt,
		//CreatedAt:         arg.CreatedAt,
	}

	ctx.JSON(http.StatusCreated, res)
}
