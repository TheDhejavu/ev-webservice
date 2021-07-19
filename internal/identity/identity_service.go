package identity

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/workspace/evoting/ev-webservice/internal/entity"
	crypto "github.com/workspace/evoting/ev-webservice/pkg/crypto"
	"github.com/workspace/evoting/ev-webservice/pkg/facialrecognition"
	"github.com/workspace/evoting/ev-webservice/pkg/log"
	"github.com/workspace/evoting/ev-webservice/wallet"
)

type identityService struct {
	identityRepo entity.IdentityRepository
	wallets      *wallet.Wallets
	logger       log.Logger
}

func NewIdentityService(identityRepo entity.IdentityRepository, logger log.Logger) entity.IdentityService {
	// Initialize system identity wallet
	wallets, _ := wallet.InitializeWallets()

	return &identityService{
		identityRepo: identityRepo,
		wallets:      wallets,
		logger:       logger,
	}
}

func (service *identityService) populateKeys(identity *entity.IdentityRead) (err error) {
	w, err := service.wallets.GetWallet(identity.ID.Hex())

	if err != nil {
		fmt.Println("HEKKEKEKEK")
		return
	}

	identity.Wallet.Certificate = string(w.Certificate[:])
	identity.Wallet.PublicMainKey = string(wallet.Base58Encode(w.Main.PublicKey))
	identity.Wallet.PublicViewKey = string(wallet.Base58Encode(w.View.PublicKey))

	return nil
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
	err = service.populateKeys(&res)
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
	err = service.populateKeys(&res)

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

func (service *identityService) Create(ctx context.Context, data map[string]interface{}, facialImages []string) (res entity.IdentityRead, err error) {
	jsonbody, err := json.Marshal(data)
	if err != nil {
		return
	}

	Identity := &entity.Identity{}
	if err = json.Unmarshal(jsonbody, &Identity); err != nil {
		return
	}
	// Hash user password
	hashedPassword, err := crypto.HashPassword(Identity.Password)
	Identity.Password = hashedPassword

	if err != nil {
		return
	}
	// Add new identity to the underlying data layer
	res, err = service.identityRepo.Create(ctx, *Identity)
	if err != nil {
		return
	}
	// Register Captured facial images of the user
	fg := facialrecognition.NewFacialRecogntion(service.logger)
	err = fg.Register(res.ID.Hex(), facialImages)
	if err != nil {
		service.identityRepo.Delete(ctx, res.ID.Hex())
		return
	}

	// Add new identity to the wallet with the User ID
	service.wallets.AddWallet(res.ID.Hex())
	service.wallets.Save()
	err = service.populateKeys(&res)
	if err != nil {
		service.identityRepo.Delete(ctx, res.ID.Hex())
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
