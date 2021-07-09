package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/workspace/evoting/ev-webservice/internal/entity"
	"github.com/workspace/evoting/ev-webservice/internal/utils"
	"github.com/workspace/evoting/ev-webservice/pkg/log"
	"github.com/workspace/evoting/ev-webservice/pkg/token"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

type authMiddleware struct {
	userService     entity.UserService
	identityService entity.IdentityService
	tokenMaker      token.Maker
	logger          log.Logger
}

func NewAuthMiddleware(userService entity.UserService, identityService entity.IdentityService, tokenMaker token.Maker, logger log.Logger) entity.AuthMiddleware {
	return &authMiddleware{userService, identityService, tokenMaker, logger}
}

// HasAdmin checks that user has an admin permission.
func (m authMiddleware) isAdmin(user entity.User) bool {
	m.logger.Info("Checking Role....")
	return user.Role == "admin"
}

// AuthMiddleware creates a gin middleware for authorization
func (m authMiddleware) AuthRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := m.tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}

		if payload.Identity == false {
			_, err = m.userService.GetByUsername(ctx, payload.Data)
		} else {
			digits, _ := strconv.ParseUint(payload.Data, 10, 64)
			_, err = m.identityService.GetByDigits(ctx, digits)
		}

		if err != nil {
			m.logger.Error("Authentication failed.")
			err = errors.New("Unauthorized Access.")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)

		return
	}
}

// AuthMiddleware creates a gin middleware for admin authorization
func (m authMiddleware) AdminRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

		if authPayload.Identity {
			m.logger.Error("Authentication failed.")
			err = errors.New("Unauthorized Access.")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}
		user, err := m.userService.GetByUsername(ctx, authPayload.Data)

		if err != nil {
			m.logger.Error("Authentication failed.")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}
		if m.isAdmin(user) == false {
			m.logger.Error("Unauthorized access.")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}

		return
	}
}

func (m authMiddleware) GetUser(ctx *gin.Context) *token.Payload {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	return authPayload
}
