package mgo

import (
	"context"
	"time"

	"github.com/OpenIMSDK/tools/mgoutil"

	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/db/table/relation"

	"github.com/OpenIMSDK/protocol/user"
	"github.com/OpenIMSDK/tools/pagination"

	"github.com/OpenIMSDK/tools/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewUserMongo - 初始化用户表
func NewUserMongo(db *mongo.Database) (relation.UserModelInterface, error) {
	coll := db.Collection("user")
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &UserMgo{coll: coll}, nil
}

type UserMgo struct {
	coll *mongo.Collection
}

func (u UserMgo) Create(ctx context.Context, users []*relation.UserModel) (err error) {
	return mgoutil.InsertMany(ctx, u.coll, users)
}

func (u UserMgo) UpdateByMap(ctx context.Context, userID string, args map[string]any) (err error) {
	//TODO implement me
	panic("implement me")
}

func (u UserMgo) Find(ctx context.Context, userIDs []string) (users []*relation.UserModel, err error) {
	return mgoutil.Find[*relation.UserModel](ctx, u.coll, bson.M{"user_id": bson.M{"$in": userIDs}})
}

func (u UserMgo) Take(ctx context.Context, userID string) (user *relation.UserModel, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UserMgo) TakeNotification(ctx context.Context, level int64) (user []*relation.UserModel, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UserMgo) TakeByNickname(ctx context.Context, nickname string) (user []*relation.UserModel, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UserMgo) Page(ctx context.Context, pagination pagination.Pagination) (count int64, users []*relation.UserModel, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UserMgo) PageFindUser(ctx context.Context, level1 int64, level2 int64, pagination pagination.Pagination) (count int64, users []*relation.UserModel, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UserMgo) PageFindUserWithKeyword(ctx context.Context, level1 int64, level2 int64, userID, nickName string, pagination pagination.Pagination) (count int64, users []*relation.UserModel, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UserMgo) Exist(ctx context.Context, userID string) (exist bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UserMgo) GetAllUserID(ctx context.Context, pagination pagination.Pagination) (count int64, userIDs []string, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UserMgo) GetUserGlobalRecvMsgOpt(ctx context.Context, userID string) (opt int, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UserMgo) CountTotal(ctx context.Context, before *time.Time) (count int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UserMgo) CountRangeEverydayTotal(ctx context.Context, start time.Time, end time.Time) (map[string]int64, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserMgo) AddUserCommand(ctx context.Context, userID string, Type int32, UUID string, value string, ex string) error {
	//TODO implement me
	panic("implement me")
}

func (u UserMgo) DeleteUserCommand(ctx context.Context, userID string, Type int32, UUID string) error {
	//TODO implement me
	panic("implement me")
}

func (u UserMgo) UpdateUserCommand(ctx context.Context, userID string, Type int32, UUID string, val map[string]any) error {
	//TODO implement me
	panic("implement me")
}

func (u UserMgo) GetUserCommand(ctx context.Context, userID string, Type int32) ([]*user.CommandInfoResp, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserMgo) GetAllUserCommand(ctx context.Context, userID string) ([]*user.AllCommandInfoResp, error) {
	//TODO implement me
	panic("implement me")
}
