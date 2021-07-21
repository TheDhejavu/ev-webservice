package chain

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/workspace/evoting/ev-webservice/internal/entity"
	customErr "github.com/workspace/evoting/ev-webservice/internal/errors"
	"github.com/workspace/evoting/ev-webservice/internal/utils"
	"github.com/workspace/evoting/ev-webservice/pkg/log"
)

// AccreditationHandler  represent the httphandler for Accreditations
type BlockchainHandler struct {
	logger            log.Logger
	authMiddleware    entity.AuthMiddleware
	blockchainService entity.BlockchainService
	v                 *utils.CustomValidator
}

// RegisterHandlers registers handlers for different HTTP requests.
func RegisterHandlers(
	router *gin.RouterGroup,
	authMiddleware entity.AuthMiddleware,
	blockchainService entity.BlockchainService,
	logger log.Logger,
) {
	handler := &BlockchainHandler{
		authMiddleware:    authMiddleware,
		blockchainService: blockchainService,
		logger:            logger,
		v:                 utils.CustomValidators(),
	}
	router.GET("/blockchain", handler.GetBlockchain)
}

func (handler BlockchainHandler) GetBlockchain(ctx *gin.Context) {

	res, err := handler.blockchainService.QueryBlockchain()
	if err != nil {
		utils.GinErrorResponse(
			ctx,
			customErr.InternalServerError(err.Error()),
		)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    res.Data,
		"message": "Success",
	})
}
