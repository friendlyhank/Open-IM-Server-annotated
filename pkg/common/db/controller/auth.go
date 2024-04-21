package controller

import (
	"context"

	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/config"
)

type AuthDatabase interface {
	// If the result is empty, no error is returned.
	GetTokensWithoutError(ctx context.Context, userID string, platformID int) (map[string]int, error)
	// Create token
	CreateToken(ctx context.Context, userID string, platformID int) (string, error)
}

type authDatabase struct {
	config *config.GlobalConfig
}

func NewAuthDatabase(config *config.GlobalConfig) AuthDatabase {
	return &authDatabase{config: config}
}

func (a authDatabase) GetTokensWithoutError(ctx context.Context, userID string, platformID int) (map[string]int, error) {
	//TODO implement me
	panic("implement me")
}

func (a authDatabase) CreateToken(ctx context.Context, userID string, platformID int) (string, error) {
	//TODO implement me
	panic("implement me")
}
