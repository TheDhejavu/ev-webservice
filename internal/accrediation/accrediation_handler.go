package accrediation

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/workspace/evoting/ev-webservice/internal/entity"
	"github.com/workspace/evoting/ev-webservice/internal/utils"

	customErr "github.com/workspace/evoting/ev-webservice/internal/errors"
	"github.com/workspace/evoting/ev-webservice/pkg/log"
)

// AccreditationHandler  represent the httphandler for Accreditations
type AccreditationHandler struct {
	logger               log.Logger
	authMiddleware       entity.AuthMiddleware
	accreditationService entity.AccreditationService
	v                    *utils.CustomValidator
}

// RegisterHandlers registers handlers for different HTTP requests.
func RegisterHandlers(
	router *gin.RouterGroup,
	authMiddleware entity.AuthMiddleware,
	accreditationService entity.AccreditationService,
	logger log.Logger,
) {
	handler := &AccreditationHandler{
		authMiddleware:       authMiddleware,
		accreditationService: accreditationService,
		logger:               logger,
		v:                    utils.CustomValidators(),
	}
	router.POST("/accreditation/:id/start",
		authMiddleware.AuthRequired(),
		authMiddleware.AdminRequired(),
		handler.StartAccreditation,
	)
	router.POST("/accreditation/:id/stop",
		authMiddleware.AuthRequired(),
		authMiddleware.AdminRequired(),
		handler.StopAccreditation,
	)
	router.POST("/accreditation/:id/accredite",
		authMiddleware.AuthRequired(),
		handler.Accredite,
	)
}

type requestParams struct {
	Id string `uri:"id" validate:"required"`
}

func (handler AccreditationHandler) StartAccreditation(ctx *gin.Context) {
	var params requestParams
	if err := ctx.ShouldBindUri(&params); err != nil {
		handler.logger.Error(err)
		utils.GinErrorResponse(
			ctx,
			customErr.BadRequest(err.Error()),
		)
		return
	}

	result, err := handler.accreditationService.Start(ctx, params.Id)

	if err != nil {
		handler.logger.Error(err)
		switch err {
		case entity.ErrNotFound:
			utils.GinErrorResponse(
				ctx,
				customErr.NotFound("Election with provided ID does not exist"),
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
	ctx.JSON(http.StatusOK, gin.H{
		"data":    result,
		"message": "Successfully Started Accreditation process for election",
	})
}

func (handler AccreditationHandler) StopAccreditation(ctx *gin.Context) {
	var params requestParams
	if err := ctx.ShouldBindUri(&params); err != nil {
		handler.logger.Error(err)
		utils.GinErrorResponse(
			ctx,
			customErr.BadRequest(err.Error()),
		)
		return
	}

	result, err := handler.accreditationService.Stop(ctx, params.Id)

	if err != nil {
		handler.logger.Error(err)
		switch err {
		case entity.ErrNotFound:
			utils.GinErrorResponse(
				ctx,
				customErr.NotFound("Election with provided ID does not exist"),
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
	ctx.JSON(http.StatusOK, gin.H{
		"data":    result,
		"message": "Successfully Stoped Accreditation process for election",
	})
}

type AccrediteRequest struct {
	facialImage *multipart.FileHeader `form:"facial_image" validate:"required"`
}

func (handler AccreditationHandler) Accredite(ctx *gin.Context) {
	var accrediteRequest AccrediteRequest
	var params requestParams
	if err := ctx.ShouldBind(&accrediteRequest); err != nil {
		handler.logger.Error(err)
		if errors.Is(err, io.EOF) {
			err = errors.New("Please Provide a valid information")
		}
		utils.GinErrorResponse(ctx, customErr.BadRequest(err.Error()))
		return
	}

	if err := ctx.ShouldBindUri(&params); err != nil {
		handler.logger.Error(err)
		utils.GinErrorResponse(ctx, customErr.BadRequest(err.Error()))
		return
	}

	identity, err := handler.authMiddleware.GetIdentity(ctx)
	if err != nil {
		handler.logger.Error(err)
		switch err {
		case entity.ErrNotFound:
			utils.GinErrorResponse(
				ctx,
				customErr.NotFound("Invalid identity"),
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
	err = handler.v.Validator.Struct(accrediteRequest)
	if err != nil {
		utils.GinErrorResponse(
			ctx,
			customErr.InternalServerError(err.Error()),
		)
		return
	}
	file, err := ctx.FormFile("facial_image")
	if err != nil {
		handler.logger.With(ctx).Error(err)
		utils.GinErrorResponse(
			ctx,
			customErr.BadRequest(fmt.Sprintf("get form err: %s", err.Error())),
		)
	}
	ext := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%s_%s%s", utils.GenUUID(), "facial_image", ext)
	facialImagePath := filepath.Join("./tmp", fileName)
	if err := ctx.SaveUploadedFile(file, facialImagePath); err != nil {
		handler.logger.With(ctx).Error(err)
		utils.GinErrorResponse(
			ctx,
			customErr.BadRequest(fmt.Sprintf("upload file err: %s", err.Error())),
		)
		return
	}

	res, err := handler.accreditationService.CreateBallot(
		ctx,
		params.Id,
		identity.ID.Hex(),
		facialImagePath,
	)
	if err != nil {
		utils.GinErrorResponse(
			ctx,
			customErr.InternalServerError(err.Error()),
		)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data":    res,
		"message": "Successfully accredited",
	})
}
