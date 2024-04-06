package discoveryregister

import (
	"errors"
	"os"

	"github.com/OpenIMSDK/tools/errs"

	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/discoveryregister/zookeeper"

	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/config"
)

// NewDiscoveryRegister creates a new service discovery and registry client based on the provided environment type.
func NewDiscoveryRegister(config *config.GlobalConfig) (discoveryregistry.SvcDiscoveryRegistry, error) {

	if os.Getenv("ENVS_DISCOVERY") != "" {
		config.Envs.Discovery = os.Getenv("ENVS_DISCOVERY")
	}

	switch config.Envs.Discovery {
	case "zookeeper":
		return zookeeper.NewZookeeperDiscoveryRegister(config)
	default:
		return nil, errs.Wrap(errors.New("envType not correct"))
	}
}
