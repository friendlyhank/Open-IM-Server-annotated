package unrelation

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/mw/specialerror"
	"github.com/friendlyhank/open-im-server-annotated/v3/pkg/common/config"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	maxRetry         = 10 // number of retries
	mongoConnTimeout = 10 * time.Second
)

type Mongo struct {
	db     *mongo.Client
	config *config.GlobalConfig
}

// NewMongo - 初始化mongodb client
func NewMongo(config *config.GlobalConfig) (*Mongo, error) {
	specialerror.AddReplace(mongo.ErrNoDocuments, errs.ErrRecordNotFound)
	uri := buildMongoURI(config)

	var mongoClient *mongo.Client
	var err error

	// Retry connecting to MongoDB
	for i := 0; i <= maxRetry; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), mongoConnTimeout)
		defer cancel()
		mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err == nil {
			if err = mongoClient.Ping(ctx, nil); err != nil {
				return nil, errs.Wrap(err, uri)
			}
			return &Mongo{db: mongoClient, config: config}, nil
		}
		if shouldRetry(err) {
			time.Sleep(time.Second) // exponential backoff could be implemented here
			continue
		}
	}
	return nil, errs.Wrap(err, uri)
}

func buildMongoURI(config *config.GlobalConfig) string {
	uri := os.Getenv("MONGO_URI")
	if uri != "" {
		return uri
	}

	if config.Mongo.Uri != "" {
		return config.Mongo.Uri
	}

	username := os.Getenv("MONGO_OPENIM_USERNAME")
	password := os.Getenv("MONGO_OPENIM_PASSWORD")
	address := os.Getenv("MONGO_ADDRESS")
	port := os.Getenv("MONGO_PORT")
	database := os.Getenv("MONGO_DATABASE")
	maxPoolSize := os.Getenv("MONGO_MAX_POOL_SIZE")

	if username == "" {
		username = config.Mongo.Username
	}
	if password == "" {
		password = config.Mongo.Password
	}
	if address == "" {
		address = strings.Join(config.Mongo.Address, ",")
	} else if port != "" {
		address = fmt.Sprintf("%s:%s", address, port)
	}
	if database == "" {
		database = config.Mongo.Database
	}
	if maxPoolSize == "" {
		maxPoolSize = fmt.Sprint(config.Mongo.MaxPoolSize)
	}

	if username != "" && password != "" {

		return fmt.Sprintf("mongodb://%s:%s@%s/%s?maxPoolSize=%s", username, password, address, database, maxPoolSize)
	}
	return fmt.Sprintf("mongodb://%s/%s?maxPoolSize=%s", address, database, maxPoolSize)
}

func shouldRetry(err error) bool {
	if cmdErr, ok := err.(mongo.CommandError); ok {
		return cmdErr.Code != 13 && cmdErr.Code != 18
	}
	return true
}

// GetClient returns the MongoDB client.
func (m *Mongo) GetClient() *mongo.Client {
	return m.db
}

// GetDatabase returns the specific database from MongoDB.
func (m *Mongo) GetDatabase(database string) *mongo.Database {
	return m.db.Database(database)
}
