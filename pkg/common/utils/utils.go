package utils

import (
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
)

// FriendRequestDBCopyOpenIM - 好友请求信息拷贝
func FriendRequestDBCopyOpenIM(dst *open_im_sdk.FriendRequest, src *db.FriendRequest) error {
	utils.CopyStructFields(dst, src)
	user, err := imdb.GetUserByUserID(src.FromUserID)
	if err != nil {
		return utils.Wrap(err, "")
	}
	dst.FromNickname = user.Nickname
	dst.FromFaceURL = user.FaceURL
	dst.FromGender = user.Gender
	user, err = imdb.GetUserByUserID(src.ToUserID)
	if err != nil {
		return utils.Wrap(err, "")
	}
	dst.ToNickname = user.Nickname
	dst.ToFaceURL = user.FaceURL
	dst.ToGender = user.Gender
	dst.CreateTime = uint32(src.CreateTime.Unix())
	dst.HandleTime = uint32(src.HandleTime.Unix())
	return nil
}
