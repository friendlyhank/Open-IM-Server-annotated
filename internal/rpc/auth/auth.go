package auth

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	pbAuth "Open_IM/pkg/proto/auth"
	"Open_IM/pkg/utils"
	"context"
)

func (rpc *rpcAuth) UserRegister(_ context.Context, req *pbAuth.UserRegisterReq) (*pbAuth.UserRegisterResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc args ", req.String())
	var user db.User
	utils.CopyStructFields(&user, req.UserInfo)
	if req.UserInfo.BirthStr != "" {
		time, err := utils.TimeStringToTime(req.UserInfo.BirthStr)
		if err != nil {
			log.NewError(req.OperationID, "TimeStringToTime failed ", err.Error(), req.UserInfo.BirthStr)
			return &pbAuth.UserRegisterResp{CommonResp: &pbAuth.CommonResp{ErrCode: constant.ErrArgs.ErrCode, ErrMsg: "TimeStringToTime failed:" + err.Error()}}, nil
		}
		user.Birth = time
	}
	log.Debug(req.OperationID, "copy ", user, req.UserInfo)
	err := imdb.UserRegister(user)
	if err != nil {
		errMsg := req.OperationID + " imdb.UserRegister failed " + err.Error() + user.UserID
		log.NewError(req.OperationID, errMsg, user)
		return &pbAuth.UserRegisterResp{CommonResp: &pbAuth.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}}, nil
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc return ", pbAuth.UserRegisterResp{CommonResp: &pbAuth.CommonResp{}})
	return &pbAuth.UserRegisterResp{CommonResp: &pbAuth.CommonResp{}}, nil
}

type rpcAuth struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}
