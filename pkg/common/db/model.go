package db

import (
	"Open_IM/pkg/common/config"
	"context"
	"fmt"
	"github.com/dtm-labs/rockscache"
	go_redis "github.com/go-redis/redis/v8"
	"time"
)

var DB DataBases

type DataBases struct {
	MysqlDB mysqlDB                  // 数据库连接
	RDB     go_redis.UniversalClient // redis连接
	Rc      *rockscache.Client       // 强一直性缓存锁
	WeakRc  *rockscache.Client       // 弱一致性缓存锁
}

func init() {
	fmt.Println("init mysql redis mongo ")
	// 初始化数据库
	initMysqlDB()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if config.Config.Redis.EnableCluster {
		// 设置redis集群
		DB.RDB = go_redis.NewClusterClient(&go_redis.ClusterOptions{
			Addrs:    config.Config.Redis.DBAddress,
			Username: config.Config.Redis.DBUserName,
			Password: config.Config.Redis.DBPassWord, // no password set
			PoolSize: 50,
		})
		_, err := DB.RDB.Ping(ctx).Result()
		if err != nil {
			fmt.Println("redis cluster failed address ", config.Config.Redis.DBAddress)
			panic(err.Error() + " redis cluster " + config.Config.Redis.DBUserName + config.Config.Redis.DBPassWord)
		}
	} else {
		DB.RDB = go_redis.NewClient(&go_redis.Options{
			Addr:     config.Config.Redis.DBAddress[0],
			Username: config.Config.Redis.DBUserName,
			Password: config.Config.Redis.DBPassWord, // no password set
			DB:       0,                              // use default DB
			PoolSize: 100,                            // 连接池大小
		})
		_, err := DB.RDB.Ping(ctx).Result()
		if err != nil {
			panic(err.Error() + " redis " + config.Config.Redis.DBAddress[0] + config.Config.Redis.DBUserName + config.Config.Redis.DBPassWord)
		}
	}
	// 强一致性缓存，当一个key被标记删除，其他请求线程会被锁住轮询直到新的key生成，适合各种同步的拉取, 如果弱一致可能导致拉取还是老数据，毫无意义
	DB.Rc = rockscache.NewClient(DB.RDB, rockscache.NewDefaultOptions())
	DB.Rc.Options.StrongConsistency = true

	// 弱一致性缓存，当一个key被标记删除，其他请求线程直接返回该key的value，适合高频并且生成很缓存很慢的情况 如大群发消息缓存的缓存
	DB.WeakRc = rockscache.NewClient(DB.RDB, rockscache.NewDefaultOptions())
	DB.WeakRc.Options.StrongConsistency = false
	fmt.Println("init mysql redis mongo ok ")
}
