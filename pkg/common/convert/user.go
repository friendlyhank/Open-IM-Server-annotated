package convert

import (
	"github.com/OpenIMSDK/protocol/sdkws"
	relationtb "github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/db/table/relation"
)

func UsersDB2Pb(users []*relationtb.UserModel) []*sdkws.UserInfo {
	result := make([]*sdkws.UserInfo, 0, len(users))
	for _, user := range users {
		userPb := &sdkws.UserInfo{
			UserID:     user.UserID,
			Nickname:   user.Nickname,
			FaceURL:    user.FaceURL,
			CreateTime: user.CreateTime.UnixMilli(),
		}
		result = append(result, userPb)
	}
	return result
}
