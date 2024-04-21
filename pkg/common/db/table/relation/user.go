package relation

import (
	"context"
	"time"

	"github.com/OpenIMSDK/protocol/user"
	"github.com/OpenIMSDK/tools/pagination"
)

type UserModel struct {
	UserID     string    `bson:"user_id"`
	Nickname   string    `bson:"nickname"`
	FaceURL    string    `bson:"face_url"`
	CreateTime time.Time `bson:"create_time"`
}

type UserModelInterface interface {
	Create(ctx context.Context, users []*UserModel) (err error)
	UpdateByMap(ctx context.Context, userID string, args map[string]any) (err error)
	Find(ctx context.Context, userIDs []string) (users []*UserModel, err error)
	Take(ctx context.Context, userID string) (user *UserModel, err error)
	TakeNotification(ctx context.Context, level int64) (user []*UserModel, err error)
	TakeByNickname(ctx context.Context, nickname string) (user []*UserModel, err error)
	Page(ctx context.Context, pagination pagination.Pagination) (count int64, users []*UserModel, err error)
	PageFindUser(ctx context.Context, level1 int64, level2 int64, pagination pagination.Pagination) (count int64, users []*UserModel, err error)
	PageFindUserWithKeyword(ctx context.Context, level1 int64, level2 int64, userID, nickName string, pagination pagination.Pagination) (count int64, users []*UserModel, err error)
	Exist(ctx context.Context, userID string) (exist bool, err error)
	GetAllUserID(ctx context.Context, pagination pagination.Pagination) (count int64, userIDs []string, err error)
	GetUserGlobalRecvMsgOpt(ctx context.Context, userID string) (opt int, err error)
	// Get user total quantity
	CountTotal(ctx context.Context, before *time.Time) (count int64, err error)
	// Get user total quantity every day
	CountRangeEverydayTotal(ctx context.Context, start time.Time, end time.Time) (map[string]int64, error)
	//CRUD user command
	AddUserCommand(ctx context.Context, userID string, Type int32, UUID string, value string, ex string) error
	DeleteUserCommand(ctx context.Context, userID string, Type int32, UUID string) error
	UpdateUserCommand(ctx context.Context, userID string, Type int32, UUID string, val map[string]any) error
	GetUserCommand(ctx context.Context, userID string, Type int32) ([]*user.CommandInfoResp, error)
	GetAllUserCommand(ctx context.Context, userID string) ([]*user.AllCommandInfoResp, error)
}
