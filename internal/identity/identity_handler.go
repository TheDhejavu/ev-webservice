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
	service         entity.IdentityService
	countryService  entity.CountryService
	identityService entity.IdentityService
	authMiddleware  entity.AuthMiddleware
	logger          log.Logger
	v               *utils.CustomValidator
	conf            config.Config
}

// RegisterHandlers will initialize the Identities resources endpoint
func RegisterHandlers(
	router *gin.RouterGroup,
	service entity.IdentityService,
	countryService entity.CountryService,
	identityService entity.IdentityService,
	authMiddleware entity.AuthMiddleware,
	conf config.Config,
	logger log.Logger,
) {
	handler := &identityHandler{
		service:         service,
		countryService:  countryService,
		identityService: identityService,
		authMiddleware:  authMiddleware,
		logger:          logger,
		conf:            conf,
		v:               utils.CustomValidators(),
	}

	router.GET("/identity", handler.GetIdentities)
	router.POST("/identity", handler.CreateIdentity)
	router.GET("/identity/me",
		authMiddleware.AuthRequired(),
		handler.GetCurrentIdentity,
	)
	router.GET("/identity/:id", handler.GetIdentity)
	router.DELETE("/identity/:id", handler.DeleteIdentity)
	router.PUT("/identity/:id", handler.UpdateIdentity)
}

// CreateIdentity will create new Identities
func (handler identityHandler) CreateIdentity(ctx *gin.Context) {
	var body createIdentityRequest
	if err := ctx.ShouldBind(&body); err != nil {
		fmt.Println(err)
		if errors.Is(err, io.EOF) {
			err = errors.New("Please Provide a valid Identity information")
		}

		message := "Please provide a valid information"
		handler.logger.With(ctx).Error(err)
		utils.GinErrorResponse(
			ctx,
			customErr.BadRequest(message),
		)
		return
	}

	err := body.Validate(ctx, handler)
	if err != nil {
		handler.logger.Error(err)
		utils.GinErrorResponse(
			ctx,
			customErr.InvalidRequestData(err, handler.v.Translator),
		)
		return
	}

	identity := utils.StructToMap(body)
	fileNames := []string{
		"national_id_card",
		"voter_card",
		"birth_certificate",
	}
	for i := 0; i < len(fileNames); i++ {
		file, err := ctx.FormFile(fileNames[i])
		if err != nil {
			handler.logger.With(ctx).Error(err)
			utils.GinErrorResponse(
				ctx,
				customErr.BadRequest(fmt.Sprintf("get form err: %s", err.Error())),
			)
		}
		ext := filepath.Ext(file.Filename)
		fileName := fmt.Sprintf("%s_%s%s", utils.GenUUID(), fileNames[i], ext)
		fileDestination := filepath.Join(handler.conf.FileStoragePath, fileName)
		if err := ctx.SaveUploadedFile(file, fileDestination); err != nil {
			handler.logger.With(ctx).Error(err)
			utils.GinErrorResponse(
				ctx,
				customErr.BadRequest(fmt.Sprintf("upload file err: %s", err.Error())),
			)
			return
		}
		identity[fileNames[i]] = fileName
	}
	// Multipart form
	form, _ := ctx.MultipartForm()

	facialImages := form.File["facial_images"]
	facialImagesFiles := []string{}

	for i := 0; i < len(facialImages); i++ {
		file := facialImages[i]
		if err != nil {
			handler.logger.With(ctx).Error(err)
			utils.GinErrorResponse(
				ctx,
				customErr.BadRequest(fmt.Sprintf("get form err: %s", err.Error())),
			)
		}

		ext := filepath.Ext(file.Filename)
		fileName := fmt.Sprintf("%s_%s%s", utils.GenUUID(), body.FirstName, ext)
		fileDestination := filepath.Join("tmp", fileName)
		if err := ctx.SaveUploadedFile(file, fileDestination); err != nil {
			handler.logger.With(ctx).Error(err)
			utils.GinErrorResponse(
				ctx,
				customErr.BadRequest(fmt.Sprintf("upload file err: %s", err.Error())),
			)
			return
		}
		facialImagesFiles = append(facialImagesFiles, fileDestination)
	}

	Identity, err := handler.service.Create(ctx, identity, facialImagesFiles)
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

// GetCurrentIdentity get current loggedin identity
func (handler identityHandler) GetCurrentIdentity(ctx *gin.Context) {
	result, err := handler.authMiddleware.GetIdentity(ctx)

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
