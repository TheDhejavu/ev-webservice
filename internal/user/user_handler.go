package user

import (
	"errors"
	"io"
	"net/http"

	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"

	"github.com/workspace/evoting/ev-webservice/internal/entity"
	"github.com/workspace/evoting/ev-webservice/internal/utils"

	customErr "github.com/workspace/evoting/ev-webservice/internal/errors"
	"github.com/workspace/evoting/ev-webservice/pkg/log"
)

// UserHandler  represent the httphandler for Users
type userHandler struct {
	service entity.UserService
	logger  log.Logger
	v       *utils.CustomValidator
}

// RegisterHandlers will initialize the Users resources endpoint
func RegisterHandlers(router *gin.RouterGroup, service entity.UserService, logger log.Logger) {
	handler := &userHandler{
		service: service,
		logger:  logger,
		v:       utils.CustomValidators(),
	}

	router.GET("/users", handler.GetUsers)
	router.POST("/users", handler.CreateUser)
	router.GET("/users/:id", handler.GetUser)
	router.DELETE("/users/:id", handler.DeleteUser)
	router.PUT("/users/:id", handler.UpdateUser)
}

// CreateUser will create new users
func (handler userHandler) CreateUser(ctx *gin.Context) {

	var userRequest createUserRequest
	if err := ctx.ShouldBindJSON(&userRequest); err != nil {
		handler.logger.Error(err)
		if errors.Is(err, io.EOF) {
			err = errors.New("Please Provide a valid user information")
		}
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	handler.logger.Info(userRequest)
	err := userRequest.Validate(ctx, handler)
	if err != nil {
		handler.logger.Error(err)
		errs := err.(validator.ValidationErrors)
		ctx.JSON(http.StatusBadRequest, utils.FormatError(errs.Translate(handler.v.Translator)))
		return
	}

	User, err := handler.service.Create(ctx, structs.Map(userRequest))
	if err != nil {
		handler.logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    User,
		"message": "Successfully Created",
	})
}

// GetUsers gets all users
func (handler *userHandler) GetUsers(ctx *gin.Context) {

	handler.logger.Info("get users")
	result, err := handler.service.Fetch(ctx, nil)
	if err != nil {
		handler.logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	handler.logger.Info(result)
	ctx.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}

// GetUser get a user with specified ID
func (handler userHandler) GetUser(ctx *gin.Context) {
	var params userRequestParams
	if err := ctx.ShouldBindUri(&params); err != nil {
		handler.logger.Error(err)
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	result, err := handler.service.GetByID(ctx, params.Id)

	if err != nil {
		handler.logger.Error(err)
		switch err {
		case customErr.ErrEntityDoesNotExist:
			ctx.JSON(
				http.StatusNotFound,
				utils.ErrorResponse(errors.New("User with provided ID does not exist")),
			)
			return
		default:
			ctx.JSON(
				http.StatusInternalServerError,
				utils.ErrorResponse(err),
			)
			return
		}
	}

	handler.logger.Info(result)
	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// UpdateUser updates a user with specified ID
func (handler userHandler) UpdateUser(ctx *gin.Context) {
	var body updateUserRequest
	var params userRequestParams
	if err := ctx.ShouldBindJSON(&body); err != nil {
		handler.logger.Error(err)
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	if err := ctx.ShouldBindUri(&params); err != nil {
		handler.logger.Error(err)
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	err := body.Validate(ctx, handler, params)
	if err != nil {
		handler.logger.Error(err)
		errs := err.(validator.ValidationErrors)
		ctx.JSON(http.StatusBadRequest, utils.FormatError(errs.Translate(handler.v.Translator)))
		return
	}

	user, err := handler.service.Update(ctx, params.Id, structs.Map(body))
	if err != nil {
		handler.logger.Error(err)
		switch err {
		case customErr.ErrEntityDoesNotExist:
			ctx.JSON(
				http.StatusNotFound,
				utils.ErrorResponse(errors.New("Unable to update user that does not exist")),
			)
			return
		default:
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    user,
		"message": "Successfully Updated",
	})
}

// DeleteUser deletes a user with specified ID
func (handler userHandler) DeleteUser(ctx *gin.Context) {
	var params userRequestParams
	if err := ctx.ShouldBindUri(&params); err != nil {
		handler.logger.Error(err)
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	err := handler.service.Delete(ctx, params.Id)
	if err != nil {
		handler.logger.Error(err)
		switch err {
		case customErr.ErrEntityDoesNotExist:
			ctx.JSON(
				http.StatusNotFound,
				utils.ErrorResponse(errors.New("User with provided ID does not exist")),
			)
			return
		default:
			ctx.JSON(
				http.StatusInternalServerError,
				utils.ErrorResponse(err),
			)
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully deleted"})
}
