package user

import (
	registry "github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/config"
	"google.golang.org/grpc"
)

func Start(config *config.GlobalConfig, client registry.SvcDiscoveryRegistry, server *grpc.Server) error {
	return nil
}
