package auth

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/workspace/evoting/ev-webservice/internal/entity"
	customErr "github.com/workspace/evoting/ev-webservice/internal/errors"
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
			utils.GinAbortResponse(
				ctx,
				customErr.Unauthorized("authorization header is not provided"),
			)
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			utils.GinAbortResponse(
				ctx,
				customErr.Unauthorized("invalid authorization header format"),
			)
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			utils.GinAbortResponse(
				ctx,
				customErr.Unauthorized(fmt.Sprintf("unsupported authorization type %s", authorizationType)),
			)
			return
		}

		accessToken := fields[1]
		payload, err := m.tokenMaker.VerifyToken(accessToken)
		if err != nil {
			utils.GinAbortResponse(
				ctx,
				customErr.Unauthorized(err.Error()),
			)
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
			utils.GinAbortResponse(
				ctx,
				customErr.Unauthorized("Unauthorized Access."),
			)
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
			utils.GinAbortResponse(
				ctx,
				customErr.Unauthorized("Unauthorized Access."),
			)
			return
		}
		user, err := m.userService.GetByUsername(ctx, authPayload.Data)

		if err != nil {
			m.logger.Error("Authentication failed.")
			utils.GinAbortResponse(
				ctx,
				customErr.Unauthorized("Authentication failed."),
			)
			return
		}
		if m.isAdmin(user) == false {
			m.logger.Error("Authentication failed.")
			utils.GinAbortResponse(
				ctx,
				customErr.Unauthorized("Unauthorized Access."),
			)
			return
		}

		return
	}
}

func (m authMiddleware) HasIdentity(ctx *gin.Context) bool {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	return authPayload.Identity
}

func (m authMiddleware) HasUser(ctx *gin.Context) bool {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	return authPayload.Identity == false
}

func (m authMiddleware) GetUser(ctx *gin.Context) (user entity.User, err error) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	user, err = m.userService.GetByID(ctx, authPayload.Data)
	if err != nil {
		return
	}
	return user, nil
}

func (m authMiddleware) GetIdentity(ctx *gin.Context) (identity entity.IdentityRead, err error) {
	var digits uint64
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Identity {
		digits, err = strconv.ParseUint(authPayload.Data, 0, 64)
		if err != nil {
			return
		}
		identity, err = m.identityService.GetByDigits(ctx, digits)
		if err != nil {
			return
		}
		return identity, nil
	}

	return
}
