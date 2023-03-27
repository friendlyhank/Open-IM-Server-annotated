/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@open-im.io").
** time(2021/3/5 14:31).
 */
package logic

/*
 * 发送消息给用户逻辑
 */

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbPush "Open_IM/pkg/proto/push"
	"strings"
)

func MsgToUser(pushMsg *pbPush.PushMsgReq) {
	log.Debug(pushMsg.OperationID, "Get msg from msg_transfer And push msg", pushMsg.String())
	grpcCons := getcdv3.GetDefaultGatewayConn4Unique(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), pushMsg.OperationID)

}
