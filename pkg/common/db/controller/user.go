package controller

import (
	"context"
	"time"

	"github.com/OpenIMSDK/tools/errs"

	"github.com/OpenIMSDK/tools/tx"

	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/db/cache"

	"github.com/OpenIMSDK/protocol/user"
	"github.com/OpenIMSDK/tools/pagination"
	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/db/table/relation"
)

type UserDatabase interface {
	// FindWithError Get the information of the specified user. If the userID is not found, it will also return an error
	FindWithError(ctx context.Context, userIDs []string) (users []*relation.UserModel, err error)
	// Find Get the information of the specified user If the userID is not found, no error will be returned
	Find(ctx context.Context, userIDs []string) (users []*relation.UserModel, err error)
	// Find userInfo By Nickname
	FindByNickname(ctx context.Context, nickname string) (users []*relation.UserModel, err error)
	// Find notificationAccounts
	FindNotification(ctx context.Context, level int64) (users []*relation.UserModel, err error)
	// Create Insert multiple external guarantees that the userID is not repeated and does not exist in the db
	Create(ctx context.Context, users []*relation.UserModel) (err error)
	// Update update (non-zero value) external guarantee userID exists
	//Update(ctx context.Context, user *relation.UserModel) (err error)
	// UpdateByMap update (zero value) external guarantee userID exists
	UpdateByMap(ctx context.Context, userID string, args map[string]any) (err error)
	// FindUser
	PageFindUser(ctx context.Context, level1 int64, level2 int64, pagination pagination.Pagination) (count int64, users []*relation.UserModel, err error)
	//FindUser with keyword
	PageFindUserWithKeyword(ctx context.Context, level1 int64, level2 int64, userID string, nickName string, pagination pagination.Pagination) (count int64, users []*relation.UserModel, err error)
	// Page If not found, no error is returned
	Page(ctx context.Context, pagination pagination.Pagination) (count int64, users []*relation.UserModel, err error)
	// IsExist true as long as one exists
	IsExist(ctx context.Context, userIDs []string) (exist bool, err error)
	// GetAllUserID Get all user IDs
	GetAllUserID(ctx context.Context, pagination pagination.Pagination) (int64, []string, error)
	// Get user by userID
	GetUserByID(ctx context.Context, userID string) (user *relation.UserModel, err error)
	// InitOnce Inside the function, first query whether it exists in the db, if it exists, do nothing; if it does not exist, insert it
	InitOnce(ctx context.Context, users []*relation.UserModel) (err error)
	// CountTotal Get the total number of users
	CountTotal(ctx context.Context, before *time.Time) (int64, error)
	// CountRangeEverydayTotal Get the user increment in the range
	CountRangeEverydayTotal(ctx context.Context, start time.Time, end time.Time) (map[string]int64, error)
	// SubscribeUsersStatus Subscribe a user's presence status
	SubscribeUsersStatus(ctx context.Context, userID string, userIDs []string) error
	// UnsubscribeUsersStatus unsubscribe a user's presence status
	UnsubscribeUsersStatus(ctx context.Context, userID string, userIDs []string) error
	// GetAllSubscribeList Get a list of all subscriptions
	GetAllSubscribeList(ctx context.Context, userID string) ([]string, error)
	// GetSubscribedList Get all subscribed lists
	GetSubscribedList(ctx context.Context, userID string) ([]string, error)
	// GetUserStatus Get the online status of the user
	GetUserStatus(ctx context.Context, userIDs []string) ([]*user.OnlineStatus, error)
	// SetUserStatus Set the user status and store the user status in redis
	SetUserStatus(ctx context.Context, userID string, status, platformID int32) error

	//CRUD user command
	AddUserCommand(ctx context.Context, userID string, Type int32, UUID string, value string, ex string) error
	DeleteUserCommand(ctx context.Context, userID string, Type int32, UUID string) error
	UpdateUserCommand(ctx context.Context, userID string, Type int32, UUID string, val map[string]any) error
	GetUserCommands(ctx context.Context, userID string, Type int32) ([]*user.CommandInfoResp, error)
	GetAllUserCommands(ctx context.Context, userID string) ([]*user.AllCommandInfoResp, error)
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

func (u userDatabase) Find(ctx context.Context, userIDs []string) (users []*relation.UserModel, err error) {
	//TODO implement me
	panic("implement me")
}

func (u userDatabase) FindByNickname(ctx context.Context, nickname string) (users []*relation.UserModel, err error) {
	//TODO implement me
	panic("implement me")
}

func (u userDatabase) FindNotification(ctx context.Context, level int64) (users []*relation.UserModel, err error) {
	//TODO implement me
	panic("implement me")
}

func (u userDatabase) Create(ctx context.Context, users []*relation.UserModel) (err error) {
	return u.tx.Transaction(ctx, func(ctx context.Context) error {
		if err = u.userDB.Create(ctx, users); err != nil {
			return err
		}
		return nil
	})
}

func (u userDatabase) UpdateByMap(ctx context.Context, userID string, args map[string]any) (err error) {
	//TODO implement me
	panic("implement me")
}

func (u userDatabase) PageFindUser(ctx context.Context, level1 int64, level2 int64, pagination pagination.Pagination) (count int64, users []*relation.UserModel, err error) {
	//TODO implement me
	panic("implement me")
}

func (u userDatabase) PageFindUserWithKeyword(ctx context.Context, level1 int64, level2 int64, userID string, nickName string, pagination pagination.Pagination) (count int64, users []*relation.UserModel, err error) {
	//TODO implement me
	panic("implement me")
}

func (u userDatabase) Page(ctx context.Context, pagination pagination.Pagination) (count int64, users []*relation.UserModel, err error) {
	//TODO implement me
	panic("implement me")
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

func (u userDatabase) GetAllUserID(ctx context.Context, pagination pagination.Pagination) (int64, []string, error) {
	//TODO implement me
	panic("implement me")
}

func (u userDatabase) GetUserByID(ctx context.Context, userID string) (user *relation.UserModel, err error) {
	//TODO implement me
	panic("implement me")
}

func (u userDatabase) InitOnce(ctx context.Context, users []*relation.UserModel) (err error) {
	//TODO implement me
	panic("implement me")
}

func (u userDatabase) CountTotal(ctx context.Context, before *time.Time) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (u userDatabase) CountRangeEverydayTotal(ctx context.Context, start time.Time, end time.Time) (map[string]int64, error) {
	//TODO implement me
	panic("implement me")
}

func (u userDatabase) SubscribeUsersStatus(ctx context.Context, userID string, userIDs []string) error {
	//TODO implement me
	panic("implement me")
}

func (u userDatabase) UnsubscribeUsersStatus(ctx context.Context, userID string, userIDs []string) error {
	//TODO implement me
	panic("implement me")
}

func (u userDatabase) GetAllSubscribeList(ctx context.Context, userID string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (u userDatabase) GetSubscribedList(ctx context.Context, userID string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (u userDatabase) GetUserStatus(ctx context.Context, userIDs []string) ([]*user.OnlineStatus, error) {
	//TODO implement me
	panic("implement me")
}

func (u userDatabase) SetUserStatus(ctx context.Context, userID string, status, platformID int32) error {
	//TODO implement me
	panic("implement me")
}

func (u userDatabase) AddUserCommand(ctx context.Context, userID string, Type int32, UUID string, value string, ex string) error {
	//TODO implement me
	panic("implement me")
}

func (u userDatabase) DeleteUserCommand(ctx context.Context, userID string, Type int32, UUID string) error {
	//TODO implement me
	panic("implement me")
}

func (u userDatabase) UpdateUserCommand(ctx context.Context, userID string, Type int32, UUID string, val map[string]any) error {
	//TODO implement me
	panic("implement me")
}

func (u userDatabase) GetUserCommands(ctx context.Context, userID string, Type int32) ([]*user.CommandInfoResp, error) {
	//TODO implement me
	panic("implement me")
}

func (u userDatabase) GetAllUserCommands(ctx context.Context, userID string) ([]*user.AllCommandInfoResp, error) {
	//TODO implement me
	panic("implement me")
}
