package apiAuth

import (
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary 用户注册
// @Description 用户注册
// @Tags 鉴权认证
// @ID UserRegister
// @Accept json
// @Param req body api.UserRegisterReq true "secret为openIM密钥, 详细见服务端config.yaml secret字段 <br> platform为平台ID <br> ex为拓展字段 <br> gender为性别, 0为女, 1为男"
// @Produce json
// @Success 0 {object} api.UserRegisterResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /auth/user_register [post]
func UserRegister(c *gin.Context) {
	params := api.UserRegisterReq{}
	if err := c.BindJSON(&params); err != nil {
		errMsg := " BindJSON failed " + err.Error()
		log.NewError("0", errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}

}
