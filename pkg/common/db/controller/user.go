package controller

import (
	"context"

	"github.com/OpenIMSDK/tools/errs"

	"github.com/OpenIMSDK/tools/tx"

	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/db/cache"

	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/db/table/relation"
)

type UserDatabase interface {
	// FindWithError -
	FindWithError(ctx context.Context, userIDs []string) (users []*relation.UserModel, err error)
	// Create -创建用户信息
	Create(ctx context.Context, users []*relation.UserModel) (err error)
	// IsExist - 用户是否存在
	IsExist(ctx context.Context, userIDs []string) (exist bool, err error)
}

// userDatabase -用户数据库
type userDatabase struct {
	tx     tx.CtxTx
	userDB relation.UserModelInterface
	cache  cache.UserCache
}

func NewUserDatabase(userDB relation.UserModelInterface, cache cache.UserCache, tx tx.CtxTx) UserDatabase {
	return &userDatabase{userDB: userDB, cache: cache, tx: tx}
}

func (u userDatabase) FindWithError(ctx context.Context, userIDs []string) (users []*relation.UserModel, err error) {
	users, err = u.cache.GetUsersInfo(ctx, userIDs)
	if err != nil {
		return
	}
	if len(users) != len(userIDs) {
		err = errs.ErrRecordNotFound.Wrap("userID not found")
	}
	return
}

func (u userDatabase) Create(ctx context.Context, users []*relation.UserModel) (err error) {
	return u.tx.Transaction(ctx, func(ctx context.Context) error {
		if err = u.userDB.Create(ctx, users); err != nil {
			return err
		}
		return nil
	})
}

func (u userDatabase) IsExist(ctx context.Context, userIDs []string) (exist bool, err error) {
	users, err := u.userDB.Find(ctx, userIDs)
	if err != nil {
		return false, err
	}
	if len(users) > 0 {
		return true, nil
	}
	return false, nil
}
