package politicalparty

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

// politicalPartyHandler  represent the httphandler for PoliticalPartys
type politicalPartyHandler struct {
	service        entity.PoliticalPartyService
	countryService entity.CountryService
	logger         log.Logger
	v              *utils.CustomValidator
}

// RegisterHandlers will initialize the PoliticalPartys resources endpoint
func RegisterHandlers(
	router *gin.RouterGroup,
	service entity.PoliticalPartyService,
	countryService entity.CountryService,
	logger log.Logger,
) {
	handler := &politicalPartyHandler{
		service:        service,
		countryService: countryService,
		logger:         logger,
		v:              utils.CustomValidators(),
	}

	router.GET("/political-party", handler.GetPoliticalParties)
	router.POST("/political-party", handler.CreatePoliticalParty)
	router.GET("/political-party/:id", handler.GetPoliticalParty)
	router.DELETE("/political-party/:id", handler.DeletePoliticalParty)
	router.PUT("/political-party/:id", handler.UpdatePoliticalParty)
}

// CreatePoliticalParty will create new PoliticalPartys
func (handler politicalPartyHandler) CreatePoliticalParty(ctx *gin.Context) {

	var body createPoliticalPartyRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		if errors.Is(err, io.EOF) {
			err = errors.New("Please Provide a valid Political party information")
		}
		handler.logger.With(ctx).Error(err)
		utils.GinErrorResponse(
			ctx,
			customErr.BadRequest(err.Error()),
		)
		return
	}

	handler.logger.Info(body)
	err := body.Validate(ctx, handler)
	if err != nil {
		handler.logger.Error(err)
		utils.GinErrorResponse(
			ctx,
			customErr.InvalidRequestData(err, handler.v.Translator),
		)
		return
	}

	PoliticalParty, err := handler.service.Store(ctx, utils.StructToMap(body))
	if err != nil {
		handler.logger.Error(err)
		utils.GinErrorResponse(
			ctx,
			customErr.InternalServerError(err.Error()),
		)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    PoliticalParty,
		"message": "Successfully Created",
	})
}

// GetPoliticalPartys gets all PoliticalPartys
func (handler *politicalPartyHandler) GetPoliticalParties(ctx *gin.Context) {

	handler.logger.Info("get Countries")
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

// GetPoliticalParty get a PoliticalParty with specified ID
func (handler politicalPartyHandler) GetPoliticalParty(ctx *gin.Context) {
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
				customErr.NotFound("Political party with provided ID does not exist"),
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

// UpdatePoliticalParty updates a PoliticalParty with specified ID
func (handler politicalPartyHandler) UpdatePoliticalParty(ctx *gin.Context) {
	var body updatePoliticalPartyRequest
	var params requestParams
	if err := ctx.ShouldBindJSON(&body); err != nil {
		handler.logger.Error(err)
		if errors.Is(err, io.EOF) {
			err = errors.New("Please Provide a valid Political party information")
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

	PoliticalParty, err := handler.service.Update(ctx, params.Id, utils.StructToMap(body))
	if err != nil {
		handler.logger.Error(err)
		switch err {
		case entity.ErrNotFound:
			utils.GinErrorResponse(
				ctx,
				customErr.NotFound("Unable to update Political party that does not exist"),
			)
			return
		default:
			utils.GinErrorResponse(ctx, customErr.InternalServerError(err.Error()))
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    PoliticalParty,
		"message": "Successfully Updated",
	})
}

// DeletePoliticalParty deletes a PoliticalParty with specified ID
func (handler politicalPartyHandler) DeletePoliticalParty(ctx *gin.Context) {
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
				customErr.NotFound("Political party with provided ID does not exist"),
			)
			return
		default:
			utils.GinErrorResponse(ctx, customErr.InternalServerError(err.Error()))
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully deleted"})
}
