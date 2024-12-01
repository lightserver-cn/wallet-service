package repository

import (
	"database/sql"
	"errors"

	"server/app/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewUser(db *sql.DB, logger *zap.SugaredLogger) UserInter {
	return &UserRepo{
		db:     db,
		logger: logger,
	}
}

type UserInter interface {
	CreateUser(ctx *gin.Context, mod *model.User) (*model.User, error)
	UpdateUser(ctx *gin.Context, mod *model.User) error
	GetUserByID(ctx *gin.Context, id int64) (*model.User, error)
	GetUserByUsername(ctx *gin.Context, username string) (*model.User, error)
	GetUserByEmail(ctx *gin.Context, email string) (*model.User, error)
}

type UserRepo struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

func (u *UserRepo) CreateUser(ctx *gin.Context, mod *model.User) (*model.User, error) {
	u.logger.Infof(model.LogUserInsert, mod.Username, mod.Email, mod.PasswordHash)

	var id int64
	err := u.db.QueryRowContext(ctx, model.QueryUserInsert, mod.Username, mod.Email, mod.PasswordHash).Scan(&id)
	if err != nil {
		u.logger.Errorf("CreateUser QueryUserInsert err: %s", err.Error())
		return mod, err
	}

	u.logger.Infof("User created with ID: %d", id)

	mod.ID = id

	return mod, err
}

func (u *UserRepo) UpdateUser(ctx *gin.Context, mod *model.User) error {
	u.logger.Infof(model.LogUserUpdate, mod.Username, mod.Email, mod.ID)

	_, err := u.db.ExecContext(ctx, model.QueryUserUpdate, mod.Username, mod.Email, mod.ID)
	if err != nil {
		u.logger.Errorf("UpdateUser error: %s", err.Error())
	}

	return err
}

func (u *UserRepo) GetUserByID(ctx *gin.Context, id int64) (*model.User, error) {
	return u.queryModelByField(ctx, "id", id)
}

func (u *UserRepo) GetUserByUsername(ctx *gin.Context, username string) (*model.User, error) {
	return u.queryModelByField(ctx, "username", username)
}

func (u *UserRepo) GetUserByEmail(ctx *gin.Context, email string) (*model.User, error) {
	return u.queryModelByField(ctx, "email", email)
}

// queryModelByField is a reusable function to query a model by a field.
func (u *UserRepo) queryModelByField(ctx *gin.Context, field string, value any) (*model.User, error) {
	u.logger.Infof(model.LogUserByField, field, value)

	mod := &model.User{}
	err := u.db.QueryRowContext(ctx, model.GetQueryByField(field), value).
		Scan(&mod.ID, &mod.Username, &mod.Email, &mod.Status, &mod.CreatedAt, &mod.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return mod, err
		}

		u.logger.Errorf("queryModelByField error: %s", err.Error())

		return mod, err
	}

	return mod, nil
}
