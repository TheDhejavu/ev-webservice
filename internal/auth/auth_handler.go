package auth

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/workspace/evoting/ev-webservice/internal/entity"
	"github.com/workspace/evoting/ev-webservice/internal/utils"

	customErr "github.com/workspace/evoting/ev-webservice/internal/errors"
	"github.com/workspace/evoting/ev-webservice/pkg/log"
)

// AuthHandler  represent the httphandler for Auths
type AuthHandler struct {
	service entity.AuthService
	logger  log.Logger
	v       *utils.CustomValidator
}

// RegisterHandlers registers handlers for different HTTP requests.
func RegisterHandlers(
	router *gin.RouterGroup,
	service entity.AuthService,
	logger log.Logger,
) {
	handler := &AuthHandler{
		service: service,
		logger:  logger,
		v:       utils.CustomValidators(),
	}
	router.POST("/login", handler.Login)
	router.POST("/login/identity", handler.LoginIdentity)
}

type loginUser struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required"`
}

// login returns a handler that handles user login request.
func (handler AuthHandler) Login(ctx *gin.Context) {
	var body loginUser
	if err := ctx.ShouldBindJSON(&body); err != nil {
		if errors.Is(err, io.EOF) {
			err = errors.New("Please Provide a valid login information")
		}
		handler.logger.With(ctx).Error(err)
		utils.GinErrorResponse(
			ctx,
			customErr.BadRequest(err.Error()),
		)
		return
	}

	err := handler.v.Validator.Struct(body)
	if err != nil {
		handler.logger.Error(err)
		utils.GinErrorResponse(
			ctx,
			customErr.InvalidRequestData(err, handler.v.Translator),
		)
		return
	}

	user, err := handler.service.Login(ctx, body.Username, body.Password)
	if err != nil {
		switch err {
		case entity.ErrInvalidUser:
			utils.GinErrorResponse(
				ctx,
				customErr.Unauthorized(err.Error()),
			)
			return
		default:
			handler.logger.Error(err)
			utils.GinErrorResponse(
				ctx,
				customErr.InternalServerError(err.Error()),
			)
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    user,
		"message": "Successfully logged in",
	})
}

type loginIdentity struct {
	Digits   uint64 `json:"digits" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginIdentity Log in to user identity
func (handler AuthHandler) LoginIdentity(ctx *gin.Context) {
	var body loginIdentity
	if err := ctx.ShouldBindJSON(&body); err != nil {
		if errors.Is(err, io.EOF) {
			err = errors.New("Please Provide a valid identity")
		}
		handler.logger.With(ctx).Error(err)
		utils.GinErrorResponse(
			ctx,
			customErr.BadRequest(err.Error()),
		)
		return
	}

	err := handler.v.Validator.Struct(body)
	if err != nil {
		handler.logger.Error(err)
		utils.GinErrorResponse(
			ctx,
			customErr.InvalidRequestData(err, handler.v.Translator),
		)
		return
	}

	identity, err := handler.service.LoginIdentity(ctx, body.Digits, body.Password)
	if err != nil {
		switch err {
		case entity.ErrInvalidIdentity:
			utils.GinErrorResponse(
				ctx,
				customErr.Unauthorized(err.Error()),
			)
			return
		default:
			handler.logger.Error(err)
			utils.GinErrorResponse(
				ctx,
				customErr.InternalServerError(err.Error()),
			)
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    identity,
		"message": "Successfully logged in",
	})
}
