package friend

import (
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	rpc "Open_IM/pkg/proto/friend"
	"Open_IM/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// @Summary 添加好友
// @Description 添加好友
// @Tags 好友相关
// @ID AddFriend
// @Accept json
// @Param token header string true "im token"
// @Param req body api.AddFriendReq true "reqMsg为申请信息 <br> fromUserID为申请用户 <br> toUserID为被添加用户"
// @Produce json
// @Success 0 {object} api.AddFriendResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /friend/add_friend [post]
func AddFriend(c *gin.Context) {
	params := api.AddFriendReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.AddFriendReq{CommID: &rpc.CommID{}}
	utils.CopyStructFields(req.CommID, &params.ParamsCommFriend)
	req.ReqMsg = params.ReqMsg

	var ok bool
	var errInfo string
	ok, req.CommID.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.CommID.OperationID)
	if !ok {
		errMsg := req.CommID.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.CommID.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.CommID.OperationID, "AddFriend args ", req.String())

	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName, req.CommID.OperationID)
	if etcdConn == nil {
		errMsg := req.CommID.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.CommID.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewFriendClient(etcdConn)
	RpcResp, err := client.AddFriend(context.Background(), req)
	if err != nil {
		log.NewError(req.CommID.OperationID, "AddFriend failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call AddFriend rpc server failed"})
		return
	}
	resp := api.AddFriendResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}}
	log.NewInfo(req.CommID.OperationID, "AddFriend api return ", resp)
	c.JSON(http.StatusOK, resp)
}
