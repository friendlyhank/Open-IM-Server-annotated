package api

import (
	"github.com/OpenIMSDK/protocol/user"
	"github.com/OpenIMSDK/tools/a2r"
	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/rpcclient"
	"github.com/gin-gonic/gin"
)

/*
 * user api调用实现
 */

type UserApi rpcclient.User

func NewUserApi(client rpcclient.User) UserApi {
	return UserApi(client)
}

func (u *UserApi) UserRegister(c *gin.Context) {
	a2r.Call(user.UserClient.UserRegister, u.Client, c)
}
