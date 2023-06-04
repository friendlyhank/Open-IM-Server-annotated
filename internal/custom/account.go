package custom

import (
	"Open_IM/internal/demo/register"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	rpc "Open_IM/pkg/proto/auth"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Login(c *gin.Context) {
	params := register.ParamsLogin{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}

	user, err := im_mysql_model.GetUserByPhoneNumber(params.PhoneNumber)
	if err != nil {
		errMsg := params.OperationID + " no this user"
		log.NewError(params.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	req := &rpc.UserTokenReq{Platform: params.Platform, FromUserID: user.UserID, OperationID: params.OperationID, LoginIp: ""}
	log.NewInfo(req.OperationID, "UserToken args ", req.String())
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAuthName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + " getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewAuthClient(etcdConn)
	reply, err := client.UserToken(context.Background(), req)
	if err != nil {
		errMsg := req.OperationID + " UserToken failed " + err.Error() + " req: " + req.String()
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	type loginData struct {
		ChatToken string `json:"chatToken"`
		ImToken   string `json:"imToken"`
		UserID    string `json:"userID"`
	}
	type UserTokenResp struct {
		ErrCode   int32     `json:"errCode"`
		ErrMsg    string    `json:"errMsg"`
		LoginData loginData `json:"data"`
	}
	resp := UserTokenResp{ErrCode: reply.CommonResp.ErrCode, ErrMsg: reply.CommonResp.ErrMsg,
		LoginData: loginData{UserID: req.FromUserID, ChatToken: reply.Token, ImToken: reply.Token}}
	log.NewInfo(req.OperationID, "UserToken return ", resp)
	c.JSON(http.StatusOK, resp)
}
