package country

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

// CountryHandler  represent the httphandler for Countrys
type countryHandler struct {
	service entity.CountryService
	logger  log.Logger
	v       *utils.CustomValidator
}

// RegisterHandlers will initialize the Countrys resources endpoint
func RegisterHandlers(router *gin.RouterGroup, service entity.CountryService, logger log.Logger) {
	handler := &countryHandler{
		service: service,
		logger:  logger,
		v:       utils.CustomValidators(),
	}

	router.GET("/countries", handler.GetCountries)
	router.POST("/country", handler.CreateCountry)
	router.GET("/country/:id", handler.GetCountry)
	router.DELETE("/country/:id", handler.DeleteCountry)
	router.PUT("/country/:id", handler.UpdateCountry)
}

// CreateCountry will create new Countrys
func (handler countryHandler) CreateCountry(ctx *gin.Context) {

	var countryRequest createCountryRequest
	if err := ctx.ShouldBindJSON(&countryRequest); err != nil {
		if errors.Is(err, io.EOF) {
			err = errors.New("Please Provide a valid Country information")
		}
		handler.logger.With(ctx).Error(err)
		utils.GinErrorResponse(
			ctx,
			customErr.BadRequest(err.Error()),
		)
		return
	}

	handler.logger.Info(countryRequest)
	err := countryRequest.Validate(ctx, handler)
	if err != nil {
		handler.logger.Error(err)
		utils.GinErrorResponse(
			ctx,
			customErr.InvalidRequestData(err, handler.v.Translator),
		)
		return
	}

	Country, err := handler.service.Store(ctx, utils.StructToMap(countryRequest))
	if err != nil {
		handler.logger.Error(err)
		utils.GinErrorResponse(
			ctx,
			customErr.InternalServerError(err.Error()),
		)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    Country,
		"message": "Successfully Created",
	})
}

// GetCountrys gets all Countrys
func (handler *countryHandler) GetCountries(ctx *gin.Context) {

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

// GetCountry get a Country with specified ID
func (handler countryHandler) GetCountry(ctx *gin.Context) {
	var params countryRequestParams
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
				customErr.NotFound("Country with provided ID does not exist"),
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

// UpdateCountry updates a Country with specified ID
func (handler countryHandler) UpdateCountry(ctx *gin.Context) {
	var body updateCountryRequest
	var params countryRequestParams
	if err := ctx.ShouldBindJSON(&body); err != nil {
		handler.logger.Error(err)
		if errors.Is(err, io.EOF) {
			err = errors.New("Please Provide a valid country information")
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

	country, err := handler.service.Update(ctx, params.Id, utils.StructToMap(body))
	if err != nil {
		handler.logger.Error(err)
		switch err {
		case entity.ErrNotFound:
			utils.GinErrorResponse(
				ctx,
				customErr.NotFound("Unable to update Country that does not exist"),
			)
			return
		default:
			utils.GinErrorResponse(ctx, customErr.InternalServerError(err.Error()))
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    country,
		"message": "Successfully Updated",
	})
}

// DeleteCountry deletes a Country with specified ID
func (handler countryHandler) DeleteCountry(ctx *gin.Context) {
	var params countryRequestParams
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
				customErr.NotFound("Country with provided ID does not exist"),
			)
			return
		default:
			utils.GinErrorResponse(ctx, customErr.InternalServerError(err.Error()))
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully deleted"})
}
