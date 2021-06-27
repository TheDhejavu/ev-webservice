package user

import (
	"context"
	"fmt"

	"github.com/workspace/evoting/ev-webservice/internal/entity"
	crypto "github.com/workspace/evoting/ev-webservice/pkg/crypto"
	"github.com/workspace/evoting/ev-webservice/pkg/log"
)

type userService struct {
	userRepo entity.UserRepository
	logger   log.Logger
}

// NewUserService creates a new user service.
func NewUserService(userRepo entity.UserRepository, logger log.Logger) entity.UserService {
	return &userService{
		userRepo: userRepo,
		logger:   logger,
	}
}

// Fetch returns the users with the specified filter.
func (service *userService) Fetch(ctx context.Context, filter interface{}) (res []entity.User, err error) {
	res, err = service.userRepo.Fetch(ctx, filter)
	if err != nil {
		return
	}
	return
}

// Get returns the user with the specified user ID.
func (service *userService) GetByID(ctx context.Context, id string) (res entity.User, err error) {
	res, err = service.userRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	return
}

// Update updates the user with the specified ID.
func (service *userService) Update(ctx context.Context, id string, user map[string]interface{}) (res entity.User, err error) {
	if val, ok := user["password"]; ok {
		value := fmt.Sprintf("%v", val)
		hashedPassword, err := crypto.HashPassword(value)
		user["password"] = fmt.Sprintf("%v", hashedPassword)
		if err != nil {
			return entity.User{}, err
		}
	}
	res, err = service.userRepo.Update(ctx, id, user)
	if err != nil {
		return
	}
	return
}

// Create creates a new user.
func (service *userService) Create(ctx context.Context, user map[string]interface{}) (res entity.User, err error) {
	hashedPassword, err := crypto.HashPassword(
		fmt.Sprintf("%v", user["password"]),
	)
	if err != nil {
		return
	}
	new_user := entity.User{
		Email:    fmt.Sprintf("%v", user["email"]),
		Username: fmt.Sprintf("%v", user["username"]),
		Password: hashedPassword,
		Role:     fmt.Sprintf("%v", user["role"]),
		FullName: fmt.Sprintf("%v", user["fullname"]),
	}

	res, err = service.userRepo.Create(ctx, new_user)
	if err != nil {
		return
	}
	return
}

// Delete deletes user with specified ID.
func (service *userService) Delete(ctx context.Context, id string) (err error) {
	value, _ := service.IdExists(ctx, id)
	if value == false {
		return entity.ErrNotFound
	}
	err = service.userRepo.Delete(ctx, id)
	if err != nil {
		return
	}
	return
}

// GetByEmail gets user by specified email
func (service *userService) GetByEmail(ctx context.Context, email string) (res entity.User, err error) {

	res, err = service.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return
	}
	return
}

// GetByUsername gets user by specified userbame
func (service *userService) GetByUsername(ctx context.Context, username string) (res entity.User, err error) {
	res, err = service.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return
	}
	return
}

// IsUsernameTaken checks if specified username exist already.
func (service *userService) IdExists(ctx context.Context, id string) (bool, error) {
	_, err := service.userRepo.GetByID(ctx, id)

	if err != nil {
		switch err {
		case entity.ErrNotFound:
			return false, nil
		default:
			return true, err
		}
	}

	return true, err
}

// IsUsernameTaken checks if specified username exist already.
func (service *userService) IsUsernameTaken(ctx context.Context, username string, id interface{}) (bool, error) {
	var err error
	if id == nil {
		_, err = service.userRepo.GetByUsername(ctx, username)
	} else {
		user := map[string]interface{}{
			"username": username,
		}
		exlude := map[string]interface{}{
			"id": id,
		}
		_, err = service.userRepo.GetWithExclude(ctx, user, exlude)
	}

	if err != nil {
		switch err {
		case entity.ErrNotFound:
			return false, nil
		default:
			return true, err
		}
	}

	return true, err
}

// IsEmailTaken checks if specified email exist already.
func (service *userService) IsEmailTaken(ctx context.Context, email string, id interface{}) (bool, error) {
	var err error
	if id == nil {
		_, err = service.userRepo.GetByEmail(ctx, email)
	} else {
		user := map[string]interface{}{
			"email": email,
		}
		exlude := map[string]interface{}{
			"id": id,
		}
		_, err = service.userRepo.GetWithExclude(ctx, user, exlude)
	}

	if err != nil {
		switch err {
		case entity.ErrNotFound:
			return false, nil
		default:
			return true, err
		}
	}

	return true, err
}
