package zookeeper

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/OpenIMSDK/tools/errs"

	"github.com/OpenIMSDK/tools/log"

	"github.com/OpenIMSDK/tools/discoveryregistry"
	openkeeper "github.com/OpenIMSDK/tools/discoveryregistry/zookeeper"
	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/config"
)

/*
 * 基于zookeeper的注册发现
 */

func NewZookeeperDiscoveryRegister(config *config.GlobalConfig) (discoveryregistry.SvcDiscoveryRegistry, error) {
	schema := getEnv("ZOOKEEPER_SCHEMA", config.Zookeeper.Schema)
	zkAddr := getZkAddrFromEnv(config.Zookeeper.ZkAddr)
	username := getEnv("ZOOKEEPER_USERNAME", config.Zookeeper.Username)
	password := getEnv("ZOOKEEPER_PASSWORD", config.Zookeeper.Password)

	zk, err := openkeeper.NewClient(
		zkAddr,
		schema,
		openkeeper.WithFreq(time.Hour),
		openkeeper.WithUserNameAndPassword(username, password),
		openkeeper.WithRoundRobin(),
		openkeeper.WithTimeout(10),
		openkeeper.WithLogger(log.NewZkLogger()),
	)
	if err != nil {
		uriFormat := "address:%s, username:%s, password:%s, schema:%s."
		errInfo := fmt.Sprintf(uriFormat,
			config.Zookeeper.ZkAddr,
			config.Zookeeper.Username,
			config.Zookeeper.Password,
			config.Zookeeper.Schema)
		return nil, errs.Wrap(err, errInfo)
	}
	return zk, nil
}

// getEnv returns the value of an environment variable if it exists, otherwise it returns the fallback value.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// getZkAddrFromEnv returns the Zookeeper addresses combined from the ZOOKEEPER_ADDRESS and ZOOKEEPER_PORT environment variables.
// If the environment variables are not set, it returns the fallback value.
// getZkAddrFromEnv - 从环境中获取zookeeper地址
func getZkAddrFromEnv(fallback []string) []string {
	address, addrExists := os.LookupEnv("ZOOKEEPER_ADDRESS")
	port, portExists := os.LookupEnv("ZOOKEEPER_PORT")

	if addrExists && portExists {
		addresses := strings.Split(address, ",")
		for i, addr := range addresses {
			addresses[i] = addr + ":" + port
		}
		return addresses
	}
	return fallback
}
