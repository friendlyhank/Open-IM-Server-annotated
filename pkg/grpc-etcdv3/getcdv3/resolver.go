package getcdv3

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"context"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
	"strings"
	"sync"
	"time"
)

// 用grpc+etcd实现的服务发现与注册
// todo hank 后面要重点看一下这个服务发现与注册服务

// Resolver需要实现grpc build的接口
type Resolver struct {
	cc                 resolver.ClientConn
	serviceName        string           // 服务名称
	grpcClientConn     *grpc.ClientConn // grpc连接
	cli                *clientv3.Client // etcd客户端
	schema             string           // etcd schema
	etcdAddr           string           // etcd地址
	watchStartRevision int64
}

var (
	nameResolver        = make(map[string]*Resolver)
	rwNameResolverMutex sync.RWMutex
)

func NewResolver(schema, etcdAddr, serviceName string, operationID string) (*Resolver, error) {
	// 初始化etcd客户端
	etcdCli, err := clientv3.New(clientv3.Config{
		Endpoints: strings.Split(etcdAddr, ","),
		Username:  config.Config.Etcd.UserName,
		Password:  config.Config.Etcd.Password,
	})
	if err != nil {
		log.Error(operationID, "etcd client v3 failed")
		return nil, utils.Wrap(err, "")
	}

	var r Resolver
	r.serviceName = serviceName
	r.cli = etcdCli
	r.schema = schema
	r.etcdAddr = etcdAddr

	// 注册服务
	resolver.Register(&r)
	//
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	conn, err := grpc.DialContext(ctx, GetPrefix(schema, serviceName),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithInsecure())
	log.Debug(operationID, "etcd key ", GetPrefix(schema, serviceName))
	if err == nil {
		r.grpcClientConn = conn
	}
	return &r, utils.Wrap(err, "")
}

func (r1 *Resolver) ResolveNow(rn resolver.ResolveNowOptions) {
}

func (r1 *Resolver) Close() {
}

// getConn - 获取rpc连接
func getConn(schema, etcdaddr, serviceName string, operationID string) *grpc.ClientConn {
	rwNameResolverMutex.RLock()
	r, ok := nameResolver[schema+serviceName]
	rwNameResolverMutex.RUnlock()
	if ok {
		log.Debug(operationID, "etcd key ", schema+serviceName, "value ", *r.grpcClientConn, *r)
		return r.grpcClientConn
	}

	rwNameResolverMutex.Lock()
	r, ok = nameResolver[schema+serviceName]
	if ok {
		rwNameResolverMutex.Unlock()
		log.Debug(operationID, "etcd key ", schema+serviceName, "value ", *r.grpcClientConn, *r)
		return r.grpcClientConn
	}

	r, err := NewResolver(schema, etcdaddr, serviceName, operationID)
	if err != nil {
		log.Error(operationID, "etcd failed ", schema, etcdaddr, serviceName, err.Error())
		rwNameResolverMutex.Unlock()
		return nil
	}

	log.Debug(operationID, "etcd key ", schema+serviceName, "value ", *r.grpcClientConn, *r)
	nameResolver[schema+serviceName] = r
	rwNameResolverMutex.Unlock()
	return r.grpcClientConn
}

// GetConfigConn - 兜底从配置中获取连接
func GetConfigConn(serviceName string, operationID string) *grpc.ClientConn {
	rpcRegisterIP := config.Config.RpcRegisterIP
	var err error
	if config.Config.RpcRegisterIP == "" {
		rpcRegisterIP, err = utils.GetLocalIP()
		if err != nil {
			log.Error(operationID, "GetLocalIP failed ", err.Error())
			return nil
		}
	}

	var configPortList []int
	// todo hank 缺少配置判断
	//3
	if config.Config.RpcRegisterName.OpenImMsgName == serviceName {
		configPortList = config.Config.RpcPort.OpenImMessagePort
	}
	//4
	if config.Config.RpcRegisterName.OpenImPushName == serviceName {
		configPortList = config.Config.RpcPort.OpenImPushPort
	}
	//5
	if config.Config.RpcRegisterName.OpenImRelayName == serviceName {
		configPortList = config.Config.RpcPort.OpenImMessageGatewayPort
	}
	//7
	if config.Config.RpcRegisterName.OpenImAuthName == serviceName {
		configPortList = config.Config.RpcPort.OpenImAuthPort
	}
	if len(configPortList) == 0 {
		log.Error(operationID, "len(configPortList) == 0  ")
		return nil
	}
	target := rpcRegisterIP + ":" + utils.Int32ToString(int32(configPortList[0]))
	log.Info(operationID, "rpcRegisterIP ", rpcRegisterIP, " port ", configPortList, " grpc target: ", target, " serviceName: ", serviceName)
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		log.Error(operationID, "grpc.Dail failed ", err.Error())
		return nil
	}
	log.NewDebug(operationID, utils.GetSelfFuncName(), serviceName, conn)
	return conn
}

// GetDefaultConn - 获取默认连接
func GetDefaultConn(schema, etcdaddr, serviceName string, operationID string) *grpc.ClientConn {
	con := getConn(schema, etcdaddr, serviceName, operationID)
	if con != nil {
		return con
	}
	// 兜底从配置中获取连接
	log.NewWarn(operationID, utils.GetSelfFuncName(), "conn is nil !!!!!", schema, etcdaddr, serviceName, operationID)
	con = GetConfigConn(serviceName, operationID)
	return con
}

// Build - resolver.Builder
func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	if r.cli == nil {
		return nil, fmt.Errorf("etcd clientv3 client failed, etcd:%s", target)
	}
	r.cc = cc
	log.Debug("", "Build..")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//     "%s:///%s"
	prefix := GetPrefix(r.schema, r.serviceName)
	// get key first
	resp, err := r.cli.Get(ctx, prefix, clientv3.WithPrefix())
	if err == nil {
		var addrList []resolver.Address
		for i := range resp.Kvs {
			log.Debug("", "etcd init addr: ", string(resp.Kvs[i].Value))
			addrList = append(addrList, resolver.Address{Addr: string(resp.Kvs[i].Value)})
		}
		r.cc.UpdateState(resolver.State{Addresses: addrList})
		r.watchStartRevision = resp.Header.Revision + 1
		go r.watch(prefix, addrList)
	} else {
		return nil, fmt.Errorf("etcd get failed, prefix: %s", prefix)
	}
	return r, nil
}

func (r *Resolver) Scheme() string {
	return r.schema
}

// 判断地址是否存在
func exists(addrList []resolver.Address, addr string) bool {
	for _, v := range addrList {
		if v.Addr == addr {
			return true
		}
	}
	return false
}

// remove - 移除服务
func remove(s []resolver.Address, addr string) ([]resolver.Address, bool) {
	for i := range s {
		if s[i].Addr == addr {
			s[i] = s[len(s)-1]
			return s[:len(s)-1], true
		}
	}
	return nil, false
}

// watch - 监视服务变化
func (r *Resolver) watch(prefix string, addrList []resolver.Address) {
	rch := r.cli.Watch(context.Background(), prefix, clientv3.WithPrefix(), clientv3.WithPrefix())
	for n := range rch {
		flag := 0
		for _, ev := range n.Events {
			switch ev.Type {
			case mvccpb.PUT: // 添加服务
				if !exists(addrList, string(ev.Kv.Value)) {
					flag = 1
					addrList = append(addrList, resolver.Address{Addr: string(ev.Kv.Value)})
					log.Debug("", "after add, new list: ", addrList)
				}
			case mvccpb.DELETE: // 摘除服务
				log.Debug("remove addr key: ", string(ev.Kv.Key), "value:", string(ev.Kv.Value))
				i := strings.LastIndexAny(string(ev.Kv.Key), "/")
				if i < 0 {
					return
				}
				t := string(ev.Kv.Key)[i+1:]
				log.Debug("remove addr key: ", string(ev.Kv.Key), "value:", string(ev.Kv.Value), "addr:", t)
				if s, ok := remove(addrList, t); ok {
					flag = 1
					addrList = s
					log.Debug("after remove, new list: ", addrList)
				}
			}
		}
		// 更新服务状态
		if flag == 1 {
			r.cc.UpdateState(resolver.State{Addresses: addrList})
			log.Debug("update: ", addrList)
		}
	}
}

var Conn4UniqueList []*grpc.ClientConn // 消息网关所有的连接
var Conn4UniqueListMtx sync.RWMutex    // 获取连接读写锁
var IsUpdateStart bool                 // 获取连接标识位
var IsUpdateStartMtx sync.RWMutex      // 获取连接标识位读写锁

// GetDefaultGatewayConn4Unique - 获取消息网关所有的连接
func GetDefaultGatewayConn4Unique(schema, etcdaddr, operationID string) []*grpc.ClientConn {
	// 加锁不要重复去获取连接信息
	IsUpdateStartMtx.Lock()
	if IsUpdateStart == false {
		Conn4UniqueList = getConn4Unique(schema, etcdaddr, config.Config.RpcRegisterName.OpenImRelayName)
		go func() {
			for {
				select {
				case <-time.After(time.Second * time.Duration(30)):
					Conn4UniqueListMtx.Lock()
					Conn4UniqueList = getConn4Unique(schema, etcdaddr, config.Config.RpcRegisterName.OpenImRelayName)
					Conn4UniqueListMtx.Unlock()
				}
			}
		}()
	}
	IsUpdateStart = true
	IsUpdateStartMtx.Unlock()

	Conn4UniqueListMtx.Lock()
	var clientConnList []*grpc.ClientConn
	for _, v := range Conn4UniqueList {
		clientConnList = append(clientConnList, v)
	}
	Conn4UniqueListMtx.Unlock()

	grpcConns := clientConnList
	if len(grpcConns) > 0 {
		return grpcConns
	}
	log.NewWarn(operationID, utils.GetSelfFuncName(), " len(grpcConns) == 0 ", schema, etcdaddr, config.Config.RpcRegisterName.OpenImRelayName)
	grpcConns = GetDefaultGatewayConn4UniqueFromcfg(operationID)
	log.NewDebug(operationID, utils.GetSelfFuncName(), config.Config.RpcRegisterName.OpenImRelayName, grpcConns)
	return grpcConns
}

// GetDefaultGatewayConn4UniqueFromcfg - 从配置中获取消息网关连接
func GetDefaultGatewayConn4UniqueFromcfg(operationID string) []*grpc.ClientConn {
	rpcRegisterIP := config.Config.RpcRegisterIP
	var err error
	if config.Config.RpcRegisterIP == "" {
		rpcRegisterIP, err = utils.GetLocalIP()
		if err != nil {
			log.Error("", "GetLocalIP failed ", err.Error())
			return nil
		}
	}
	var conns []*grpc.ClientConn
	configPortList := config.Config.RpcPort.OpenImMessageGatewayPort
	for _, port := range configPortList {
		target := rpcRegisterIP + ":" + utils.Int32ToString(int32(port))
		log.Info(operationID, "rpcRegisterIP ", rpcRegisterIP, " port ", configPortList, " grpc target: ", target, " serviceName: ", "msgGateway")
		conn, err := grpc.Dial(target, grpc.WithInsecure())
		if err != nil {
			log.Error(operationID, "grpc.Dail failed ", err.Error())
			continue
		}
		conns = append(conns, conn)
	}
	return conns
}

// getConn4Unique - 根据服务名称获取所有的rpc连接
func getConn4Unique(schema, etcdaddr, servicename string) []*grpc.ClientConn {
	gEtcdCli, err := clientv3.New(clientv3.Config{Endpoints: strings.Split(etcdaddr, ",")})
	if err != nil {
		log.Error("clientv3.New failed", err.Error())
		return nil
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//     "%s:///%s"
	prefix := GetPrefix4Unique(schema, servicename)

	resp, err := gEtcdCli.Get(ctx, prefix, clientv3.WithPrefix())
	//  "%s:///%s:ip:port"   -> %s:ip:port
	allService := make([]string, 0)
	if err == nil {
		for i := range resp.Kvs {
			k := string(resp.Kvs[i].Key)

			b := strings.LastIndex(k, "///")
			k1 := k[b+len("///"):]

			e := strings.Index(k1, "/")
			k2 := k1[:e]
			allService = append(allService, k2)
		}
	} else {
		gEtcdCli.Close()
		//log.Error("gEtcdCli.Get failed", err.Error())
		return nil
	}
	gEtcdCli.Close()

	allConn := make([]*grpc.ClientConn, 0)
	for _, v := range allService {
		r := getConn(schema, etcdaddr, v, "0")
		allConn = append(allConn, r)
	}

	return allConn
}
