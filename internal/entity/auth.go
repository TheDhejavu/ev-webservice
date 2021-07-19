package entity

import (
	"context"

	"github.com/gin-gonic/gin"
)

type AuthUser struct {
	AccessToken string `json:"access_token"`
	User        User   `json:"user"`
}

type AuthIdentity struct {
	AccessToken string       `json:"access_token"`
	Identity    IdentityRead `json:"identity"`
}

// AuthService encapsulates the authentication logic.
type AuthService interface {
	// login to user account
	Login(ctx context.Context, username, password string) (AuthUser, error)
	// login to identity
	LoginIdentity(ctx context.Context, digits uint64, password string) (AuthIdentity, error)
}

type AuthMiddleware interface {
	AuthRequired() gin.HandlerFunc
	AdminRequired() gin.HandlerFunc
	HasIdentity(ctx *gin.Context) bool
	HasUser(ctx *gin.Context) bool
	GetIdentity(ctx *gin.Context) (IdentityRead, error)
	GetUser(ctx *gin.Context) (User, error)
}
