package service

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"

	"server/app/model"
	"server/app/repository"
	"server/app/request"
)

func NewUser(repo repository.UserInter, repoWallet repository.WalletInter) UserInter {
	return &UserServ{
		repo:       repo,
		repoWallet: repoWallet,
	}
}

type UserInter interface {
	RegisterUser(ctx *gin.Context, req *request.ReqRegisterUser) (*model.User, error)
	UpdateUser(ctx *gin.Context, mod *model.User) error
	GetUserByID(ctx *gin.Context, id int64) (*model.User, error)
	GetUserByUsername(ctx *gin.Context, username string) (*model.User, error)
	GetUserByEmail(ctx *gin.Context, email string) (*model.User, error)
}

type UserServ struct {
	repo       repository.UserInter
	repoWallet repository.WalletInter
}

func (s *UserServ) RegisterUser(ctx *gin.Context, req *request.ReqRegisterUser) (*model.User, error) {
	mod := &model.User{}

	if req.Password == "" {
		return nil, errors.New("password cannot be empty")
	}

	if req.Email == "" {
		return nil, errors.New("email cannot be empty")
	}

	pwdHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return mod, err
	}

	mod = &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: pwdHash,
		Status:       model.UserStatusInvalid,
	}

	mod, err = s.repo.CreateUser(ctx, mod)
	if err != nil {
		return mod, err
	}

	modWallet := &model.Wallet{
		UID:     mod.ID,
		Balance: decimal.New(0, 0),
	}

	_, err = s.repoWallet.CreateWallet(ctx, modWallet)
	if err != nil {
		return nil, err
	}

	return mod, nil
}

func (s *UserServ) UpdateUser(ctx *gin.Context, mod *model.User) error {
	return s.repo.UpdateUser(ctx, mod)
}

func (s *UserServ) GetUserByID(ctx *gin.Context, id int64) (*model.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *UserServ) GetUserByUsername(ctx *gin.Context, username string) (*model.User, error) {
	return s.repo.GetUserByUsername(ctx, username)
}

func (s *UserServ) GetUserByEmail(ctx *gin.Context, email string) (*model.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}
