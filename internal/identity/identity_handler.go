package identity

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/workspace/evoting/ev-webservice/internal/config"
	"github.com/workspace/evoting/ev-webservice/internal/entity"
	"github.com/workspace/evoting/ev-webservice/internal/utils"

	customErr "github.com/workspace/evoting/ev-webservice/internal/errors"
	"github.com/workspace/evoting/ev-webservice/pkg/log"
)

// identityHandler   represent the httphandler for Identities
type identityHandler struct {
	service        entity.IdentityService
	countryService entity.CountryService
	logger         log.Logger
	v              *utils.CustomValidator
	conf           config.Config
}

// RegisterHandlers will initialize the Identities resources endpoint
func RegisterHandlers(
	router *gin.RouterGroup,
	service entity.IdentityService,
	countryService entity.CountryService,
	conf config.Config,
	logger log.Logger,
) {
	handler := &identityHandler{
		service:        service,
		countryService: countryService,
		logger:         logger,
		conf:           conf,
		v:              utils.CustomValidators(),
	}

	router.GET("/identity", handler.GetIdentities)
	router.POST("/identity", handler.CreateIdentity)
	router.GET("/identity/:id", handler.GetIdentity)
	router.DELETE("/identity/:id", handler.DeleteIdentity)
	router.PUT("/identity/:id", handler.UpdateIdentity)
}

// CreateIdentity will create new Identities
func (handler identityHandler) CreateIdentity(ctx *gin.Context) {

	var body createIdentityRequest
	if err := ctx.Bind(&body); err != nil {
		if errors.Is(err, io.EOF) {
			err = errors.New("Please Provide a valid Identity information")
		}
		handler.logger.With(ctx).Error(err)
		utils.GinErrorResponse(
			ctx,
			customErr.BadRequest(err.Error()),
		)
		return
	}

	// National ID Card
	nationalIdCard, err := ctx.FormFile("national_id_card")
	if err != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	ext := filepath.Ext(nationalIdCard.Filename)
	fileName := fmt.Sprintf("%s_%s%s", utils.GenKsuid(), "national_id_card", ext)
	fileDestination := filepath.Join(handler.conf.FileStoragePath, fileName)
	if err := ctx.SaveUploadedFile(nationalIdCard, fileDestination); err != nil {
		handler.logger.With(ctx).Error(err)
		utils.GinErrorResponse(
			ctx,
			customErr.BadRequest(fmt.Sprintf("upload file err: %s", err.Error())),
		)
		return
	}
	err = body.Validate(ctx, handler)
	if err != nil {
		handler.logger.Error(err)
		utils.GinErrorResponse(
			ctx,
			customErr.InvalidRequestData(err, handler.v.Translator),
		)
		return
	}
	identity := utils.StructToMap(body)
	identity["national_id_card"] = fileName
	identity["voter_card"] = fileName
	identity["birth_certificate"] = fileName

	Identity, err := handler.service.Create(ctx, identity)
	if err != nil {
		handler.logger.Error(err)
		utils.GinErrorResponse(
			ctx,
			customErr.InternalServerError(err.Error()),
		)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    Identity,
		"message": "Successfully Created",
	})
}

// GetIdentities gets all Identities
func (handler *identityHandler) GetIdentities(ctx *gin.Context) {

	handler.logger.Info("get Identities")
	result, err := handler.service.Fetch(ctx, nil)
	if err != nil {
		handler.logger.Error(err)
		utils.GinErrorResponse(
			ctx,
			customErr.InternalServerError(err.Error()),
		)
		return
	}

	handler.logger.Info(result)
	ctx.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}

// GetIdentity get a Identity with specified ID
func (handler identityHandler) GetIdentity(ctx *gin.Context) {
	var params requestParams
	if err := ctx.ShouldBindUri(&params); err != nil {
		handler.logger.Error(err)
		utils.GinErrorResponse(
			ctx,
			customErr.BadRequest(err.Error()),
		)
		return
	}

	result, err := handler.service.GetByID(ctx, params.Id)

	if err != nil {
		handler.logger.Error(err)
		switch err {
		case entity.ErrNotFound:
			utils.GinErrorResponse(
				ctx,
				customErr.NotFound("Identity with provided ID does not exist"),
			)
			return
		default:
			utils.GinErrorResponse(
				ctx,
				customErr.InternalServerError(err.Error()),
			)
			return
		}
	}

	handler.logger.Info(result)
	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// UpdateIdentity updates a Identity with specified ID
func (handler identityHandler) UpdateIdentity(ctx *gin.Context) {
	var body updateIdentityRequest
	var params requestParams
	if err := ctx.ShouldBindJSON(&body); err != nil {
		handler.logger.Error(err)
		if errors.Is(err, io.EOF) {
			err = errors.New("Please Provide a valid Identity information")
		}
		utils.GinErrorResponse(ctx, customErr.BadRequest(err.Error()))
		return
	}

	if err := ctx.ShouldBindUri(&params); err != nil {
		handler.logger.Error(err)
		utils.GinErrorResponse(ctx, customErr.BadRequest(err.Error()))
		return
	}

	err := body.Validate(ctx, handler, params)
	if err != nil {
		handler.logger.Error(err)
		utils.GinErrorResponse(
			ctx,
			customErr.InvalidRequestData(err, handler.v.Translator),
		)
		return
	}

	Identity, err := handler.service.Update(ctx, params.Id, utils.StructToMap(body))
	if err != nil {
		handler.logger.Error(err)
		switch err {
		case entity.ErrNotFound:
			utils.GinErrorResponse(
				ctx,
				customErr.NotFound("Unable to update Identity that does not exist"),
			)
			return
		default:
			utils.GinErrorResponse(ctx, customErr.InternalServerError(err.Error()))
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    Identity,
		"message": "Successfully Updated",
	})
}

// DeleteIdentity deletes a Identity with specified ID
func (handler identityHandler) DeleteIdentity(ctx *gin.Context) {
	var params requestParams
	if err := ctx.ShouldBindUri(&params); err != nil {
		handler.logger.Error(err)
		utils.GinErrorResponse(ctx, customErr.BadRequest(err.Error()))
		return
	}

	err := handler.service.Delete(ctx, params.Id)
	if err != nil {
		handler.logger.Error(err)
		switch err {
		case entity.ErrNotFound:
			utils.GinErrorResponse(
				ctx,
				customErr.NotFound("Identity with provided ID does not exist"),
			)
			return
		default:
			utils.GinErrorResponse(ctx, customErr.InternalServerError(err.Error()))
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully deleted"})
}
