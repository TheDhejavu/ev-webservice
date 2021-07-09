package auth

import (
	"context"
	"errors"
	"strconv"

	"github.com/workspace/evoting/ev-webservice/internal/config"
	"github.com/workspace/evoting/ev-webservice/internal/entity"
	crypto "github.com/workspace/evoting/ev-webservice/pkg/crypto"
	"github.com/workspace/evoting/ev-webservice/pkg/log"
	"github.com/workspace/evoting/ev-webservice/pkg/token"
)

type authService struct {
	identityService entity.IdentityService
	userService     entity.UserService
	tokenMaker      token.Maker
	logger          log.Logger
	config          config.Config
}

// NewService creates a new authentication service.
func NewIdentityService(
	identityService entity.IdentityService,
	userService entity.UserService,
	logger log.Logger,
	config config.Config,
	tokenMaker token.Maker,
) entity.AuthService {
	return authService{
		identityService: identityService,
		userService:     userService,
		tokenMaker:      tokenMaker,
		config:          config,
		logger:          logger,
	}
}

// Login handles user login
func (s authService) Login(ctx context.Context, username, password string) (entity.AuthUser, error) {
	user, err := s.userService.GetByUsername(ctx, username)
	authUser := entity.AuthUser{}
	if err != nil {
		switch err {
		case entity.ErrNotFound:
			return authUser, entity.ErrInvalidUser
		default:
			return authUser, err
		}
	}

	err = crypto.CheckPassword(password, user.Password)
	if err != nil {
		return authUser, entity.ErrInvalidUser
	}

	accessToken, err := s.tokenMaker.CreateToken(
		user.Username,
		false,
		s.config.TokenDuration,
	)
	if err != nil {
		return authUser, errors.New("Something went wrong while generating login token")
	}
	user.Password = ""
	authUser.AccessToken = accessToken
	authUser.User = user

	return authUser, nil
}
func (s authService) LoginIdentity(ctx context.Context, digits uint64, password string) (entity.AuthIdentity, error) {
	identity, err := s.identityService.GetByDigits(ctx, digits)
	AuthIdentity := entity.AuthIdentity{}

	if err != nil {
		switch err {
		case entity.ErrNotFound:
			return AuthIdentity, entity.ErrInvalidIdentity
		default:
			return AuthIdentity, err
		}
	}

	err = crypto.CheckPassword(password, identity.Password)
	if err != nil {
		return AuthIdentity, entity.ErrInvalidIdentity
	}

	accessToken, err := s.tokenMaker.CreateToken(
		strconv.Itoa(int(identity.Digits)),
		true,
		s.config.TokenDuration,
	)
	if err != nil {
		return AuthIdentity, errors.New("Something went wrong while generating login token")
	}
	identity.Password = ""
	AuthIdentity.AccessToken = accessToken
	AuthIdentity.Identity = identity

	return AuthIdentity, nil
}
