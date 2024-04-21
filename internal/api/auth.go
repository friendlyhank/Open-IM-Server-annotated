package api

import (
	"github.com/OpenIMSDK/protocol/auth"
	"github.com/OpenIMSDK/tools/a2r"
	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/rpcclient"
	"github.com/gin-gonic/gin"
)

type AuthApi rpcclient.Auth

func NewAuthApi(client rpcclient.Auth) AuthApi {
	return AuthApi(client)
}

func (o *AuthApi) UserToken(c *gin.Context) {
	a2r.Call(auth.AuthClient.UserToken, o.Client, c)
}
