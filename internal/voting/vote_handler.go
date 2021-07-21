package voting

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

// AccreditationHandler  represent the httphandler for Accreditations
type AccreditationHandler struct {
	authMiddleware entity.AuthMiddleware
	votingService  entity.VotingService
	v              *utils.CustomValidator
	logger         log.Logger
}

// RegisterHandlers registers handlers for different HTTP requests.
func RegisterHandlers(
	router *gin.RouterGroup,
	authMiddleware entity.AuthMiddleware,
	votingService entity.VotingService,
	logger log.Logger,
) {
	handler := &AccreditationHandler{
		authMiddleware: authMiddleware,
		logger:         logger,
		votingService:  votingService,
		v:              utils.CustomValidators(),
	}
	router.POST("/voting/:id/start",
		authMiddleware.AuthRequired(),
		authMiddleware.AdminRequired(),
		handler.StartVoting,
	)
	router.POST("/voting/:id/stop",
		authMiddleware.AuthRequired(),
		authMiddleware.AdminRequired(),
		handler.StopVoting,
	)
	router.POST("/voting/:id/cast-ballot",
		authMiddleware.AuthRequired(),
		handler.CastBallot,
	)
}

type requestParams struct {
	Id string `uri:"id" validate:"required"`
}

func (handler AccreditationHandler) StartVoting(ctx *gin.Context) {
	var params requestParams
	if err := ctx.ShouldBindUri(&params); err != nil {
		handler.logger.Error(err)
		utils.GinErrorResponse(
			ctx,
			customErr.BadRequest(err.Error()),
		)
		return
	}

	result, err := handler.votingService.Start(ctx, params.Id)

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
		"message": "Successfully Started voting process for election",
	})
}

func (handler AccreditationHandler) StopVoting(ctx *gin.Context) {
	var params requestParams
	if err := ctx.ShouldBindUri(&params); err != nil {
		handler.logger.Error(err)
		utils.GinErrorResponse(
			ctx,
			customErr.BadRequest(err.Error()),
		)
		return
	}

	result, err := handler.votingService.Stop(ctx, params.Id)

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
		"message": "Successfully Stopped voting process for election",
	})
}

type VotingRequest struct {
	Candidate string `json:"candidate"`
}

func (handler AccreditationHandler) CastBallot(ctx *gin.Context) {
	var params requestParams
	if err := ctx.ShouldBindUri(&params); err != nil {
		handler.logger.Error(err)
		utils.GinErrorResponse(
			ctx,
			customErr.BadRequest(err.Error()),
		)
		return
	}
	var votingRequest VotingRequest

	if err := ctx.ShouldBindJSON(&votingRequest); err != nil {
		handler.logger.Error(err)
		if errors.Is(err, io.EOF) {
			err = errors.New("Please Provide a valid information")
		}
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

	results, err := handler.votingService.CastVote(ctx,
		identity.ID.Hex(),
		params.Id,
		votingRequest.Candidate,
	)
	if err != nil {
		utils.GinErrorResponse(
			ctx,
			customErr.InternalServerError(err.Error()),
		)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data":    results,
		"message": "Success",
	})
}
