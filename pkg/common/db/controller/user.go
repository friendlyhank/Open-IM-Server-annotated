package controller

import (
	"context"
	"time"

	"github.com/OpenIMSDK/protocol/user"
	"github.com/OpenIMSDK/tools/pagination"
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
}
