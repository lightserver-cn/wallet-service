package controller

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"server/app/request"
	"server/app/service"
	"server/pkg/consts"
)

func NewUser(serv service.UserInter) UserInter {
	return &UserCtrl{
		serv: serv,
	}
}

type UserInter interface {
	RegisterUser(ctx *gin.Context)
	GetUserByUID(ctx *gin.Context)
}

type UserCtrl struct {
	serv service.UserInter
}

func (c *UserCtrl) RegisterUser(ctx *gin.Context) {
	req := new(request.ReqRegisterUser)
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": consts.ErrValidationFailed, "details": err.Error()})
		return
	}

	if strings.TrimSpace(req.Username) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": consts.ErrUsernameRequired})
		return
	}

	if strings.TrimSpace(req.Email) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": consts.ErrEmailRequired})
		return
	}

	if strings.TrimSpace(req.Password) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": consts.ErrPasswordRequired})
		return
	}

	resUsername, err := c.serv.GetUserByUsername(ctx, req.Username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": consts.ErrInternalServer})
		return
	}

	if resUsername != nil && resUsername.ID > 0 {
		ctx.JSON(http.StatusConflict, gin.H{"error": consts.ErrUsernameAlreadyExists})
		return
	}

	resEmail, err := c.serv.GetUserByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": consts.ErrInternalServer})
		return
	}

	if resEmail != nil && resEmail.ID > 0 {
		ctx.JSON(http.StatusConflict, gin.H{"error": consts.ErrEmailAlreadyExists})
		return
	}

	user, err := c.serv.RegisterUser(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": consts.ErrInternalServer})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (c *UserCtrl) GetUserByUID(ctx *gin.Context) {
	req := new(request.ReqUID)
	if err := ctx.ShouldBindUri(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   consts.ErrValidationFailed,
			"details": err.Error(),
		})
		return
	}

	if req.UID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": consts.ErrInvalidUID})
		return
	}

	user, err := c.serv.GetUserByID(ctx, req.UID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": consts.ErrUserNotFound})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": consts.ErrInternalServer})
		return
	}

	ctx.JSON(http.StatusOK, user)
}
