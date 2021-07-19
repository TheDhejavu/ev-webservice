package election

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

// electionHandler   represent the httphandler for Elections
type electionHandler struct {
	service         entity.ElectionService
	countryService  entity.CountryService
	partyService    entity.PoliticalPartyService
	identityService entity.IdentityService
	authMiddleware  entity.AuthMiddleware
	logger          log.Logger
	v               *utils.CustomValidator
}

// RegisterHandlers will initialize the Elections resources endpoint
func RegisterHandlers(
	router *gin.RouterGroup,
	service entity.ElectionService,
	countryService entity.CountryService,
	partyService entity.PoliticalPartyService,
	identityService entity.IdentityService,
	authMiddleware entity.AuthMiddleware,
	logger log.Logger,
) {
	handler := &electionHandler{
		service:         service,
		countryService:  countryService,
		partyService:    partyService,
		identityService: identityService,
		authMiddleware:  authMiddleware,
		logger:          logger,
		v:               utils.CustomValidators(),
	}

	router.GET("/elections",
		authMiddleware.AuthRequired(),
		handler.GetElections,
	)
	router.POST("/elections",
		authMiddleware.AuthRequired(),
		authMiddleware.AdminRequired(),
		handler.CreateElection,
	)
	router.GET("/elections/:id",
		handler.GetElection,
	)
	router.DELETE("/elections/:id",
		authMiddleware.AuthRequired(),
		authMiddleware.AdminRequired(),
		handler.DeleteElection,
	)
	router.PUT("/elections/:id",
		authMiddleware.AuthRequired(),
		authMiddleware.AdminRequired(),
		handler.UpdateElection,
	)

	router.POST("/elections/:id/start",
		authMiddleware.AuthRequired(),
		authMiddleware.AdminRequired(),
		handler.StartElection,
	)

	router.POST("/elections/:id/stop",
		authMiddleware.AuthRequired(),
		authMiddleware.AdminRequired(),
		handler.StopElection,
	)
}

// CreateElection will create new Elections
func (handler electionHandler) CreateElection(ctx *gin.Context) {

	var body createElectionRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		if errors.Is(err, io.EOF) {
			err = errors.New("Please Provide a valid election information")
		}
		handler.logger.With(ctx).Error(err)
		utils.GinErrorResponse(
			ctx,
			customErr.BadRequest(err.Error()),
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

	Election, err := handler.service.Create(ctx, utils.StructToMap(body))
	if err != nil {
		handler.logger.Error(err)
		utils.GinErrorResponse(
			ctx,
			customErr.InternalServerError(err.Error()),
		)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    Election,
		"message": "Successfully Created",
	})
}

// GetElections gets all Elections
func (handler *electionHandler) GetElections(ctx *gin.Context) {
	var result []*entity.ElectionRead
	var err error
	var filter map[string]interface{}

	hasIdentity := handler.authMiddleware.HasIdentity(ctx)
	if hasIdentity {
		identity, _ := handler.authMiddleware.GetIdentity(ctx)
		filter = map[string]interface{}{
			"country": identity.Origin.Country.ID,
		}
	}
	result, err = handler.service.Fetch(ctx, filter)
	if err != nil {
		handler.logger.Error(err)
		utils.GinErrorResponse(
			ctx,
			customErr.InternalServerError(err.Error()),
		)
	}

	handler.logger.Info(result)
	ctx.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}

// GetElection get a Election with specified ID
func (handler electionHandler) GetElection(ctx *gin.Context) {
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
	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// UpdateElection updates a Election with specified ID
func (handler electionHandler) UpdateElection(ctx *gin.Context) {
	var body updateElectionRequest
	var params requestParams
	if err := ctx.ShouldBindJSON(&body); err != nil {
		handler.logger.Error(err)
		if errors.Is(err, io.EOF) {
			err = errors.New("Please Provide a valid election information")
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
	// fmt.Println(utils.StructToMap(body))
	Election, err := handler.service.Update(ctx, params.Id, body)
	if err != nil {
		handler.logger.Error(err)
		switch err {
		case entity.ErrNotFound:
			utils.GinErrorResponse(
				ctx,
				customErr.NotFound("Unable to update election that does not exist"),
			)
			return
		default:
			utils.GinErrorResponse(ctx, customErr.InternalServerError(err.Error()))
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    Election,
		"message": "Successfully Updated",
	})
}

// DeleteElection deletes a Election with specified ID
func (handler electionHandler) DeleteElection(ctx *gin.Context) {
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
				customErr.NotFound("Election with provided ID does not exist"),
			)
			return
		default:
			utils.GinErrorResponse(ctx, customErr.InternalServerError(err.Error()))
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully deleted"})
}

// StartElection starts election
func (handler *electionHandler) StartElection(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{
		// "data": result,
	})
}

// StopElection stops election
func (handler *electionHandler) StopElection(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{
		// "data": result,
	})
}
