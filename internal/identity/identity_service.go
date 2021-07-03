package identity

import (
	"context"
	"encoding/json"

	"github.com/workspace/evoting/ev-webservice/internal/entity"
	crypto "github.com/workspace/evoting/ev-webservice/pkg/crypto"
	"github.com/workspace/evoting/ev-webservice/pkg/facialrecognition"
	"github.com/workspace/evoting/ev-webservice/pkg/log"
)

type identityService struct {
	identityRepo entity.IdentityRepository
	logger       log.Logger
}

func NewIdentityService(identityRepo entity.IdentityRepository, logger log.Logger) entity.IdentityService {
	return &identityService{
		identityRepo: identityRepo,
		logger:       logger,
	}
}

func (service *identityService) Fetch(ctx context.Context, filter interface{}) (res []entity.IdentityRead, err error) {
	res, err = service.identityRepo.Fetch(ctx, filter)
	if err != nil {
		return
	}
	return
}
func (service *identityService) GetByID(ctx context.Context, id string) (res entity.IdentityRead, err error) {
	res, err = service.identityRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	return
}

func (service *identityService) GetByEmail(ctx context.Context, email string) (res entity.IdentityRead, err error) {
	res, err = service.identityRepo.GetByEmail(ctx, email)
	if err != nil {
		return
	}
	return
}
func (service *identityService) GetByDigits(ctx context.Context, id uint64) (res entity.IdentityRead, err error) {
	res, err = service.identityRepo.GetByDigits(ctx, id)
	if err != nil {
		return
	}
	return
}
func (service *identityService) Update(ctx context.Context, id string, data map[string]interface{}) (res entity.IdentityRead, err error) {
	res, err = service.identityRepo.Update(ctx, id, data)
	if err != nil {
		return
	}
	return
}

func (service *identityService) Create(ctx context.Context, data map[string]interface{}, images []string) (res entity.IdentityRead, err error) {
	jsonbody, err := json.Marshal(data)
	if err != nil {
		return
	}

	Identity := &entity.Identity{}
	if err = json.Unmarshal(jsonbody, &Identity); err != nil {
		return
	}
	hashedPassword, err := crypto.HashPassword(Identity.Password)
	Identity.Password = hashedPassword

	if err != nil {
		return
	}
	res, err = service.identityRepo.Create(ctx, *Identity)
	if err != nil {
		return
	}
	fg := facialrecognition.NewFacialRecogntion(service.logger)
	err = fg.Register(res.ID.Hex(), images)
	if err != nil {
		return
	}
	return
}
func (service *identityService) Delete(ctx context.Context, id string) (err error) {
	value, _ := service.Exists(ctx, map[string]interface{}{"_id": id}, nil)
	if value == false {
		return entity.ErrNotFound
	}

	err = service.identityRepo.Delete(ctx, id)
	if err != nil {
		return
	}
	return
}

func (service *identityService) Exists(ctx context.Context, filter, exclude map[string]interface{}) (res bool, err error) {
	if exclude == nil {
		_, err = service.identityRepo.Get(ctx, filter)
	} else {
		_, err = service.identityRepo.GetWithExclude(ctx, filter, exclude)
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
