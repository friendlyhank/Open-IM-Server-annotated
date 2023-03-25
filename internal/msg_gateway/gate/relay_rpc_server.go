package gate

import "Open_IM/pkg/common/config"

// rpc服务相关逻辑
type RPCServer struct {
	rpcPort    int // rpc端口
	etcdSchema string
	etcdAddr   []string // etcd地址
}

// rpc服务初始化
func (r *RPCServer) onInit(rpcPort int) {
	r.rpcPort = rpcPort
	r.etcdSchema = config.Config.Etcd.EtcdSchema
	r.etcdAddr = config.Config.Etcd.EtcdAddr
}

func (r *RPCServer) run() {

}
