module Open_IM

go 1.15

require (
	github.com/Shopify/sarama v1.32.0
	github.com/antonfisher/nested-logrus-formatter v1.3.0
	github.com/bytedance/sonic v1.8.8 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/gin-gonic/gin v1.9.0
	github.com/go-playground/validator/v10 v10.13.0
	github.com/go-redis/redis/v8 v8.11.5
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.4.2
	github.com/golang/protobuf v1.5.2
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/gorilla/websocket v1.4.2
	github.com/jinzhu/copier v0.3.4
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/klauspost/cpuid/v2 v2.2.4 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v1.0.4 // indirect
	github.com/mattn/go-isatty v0.0.18 // indirect
	github.com/pelletier/go-toml/v2 v2.0.7 // indirect
	github.com/pkg/errors v0.9.1
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5
	github.com/sirupsen/logrus v1.8.1
	github.com/ugorji/go/codec v1.2.11 // indirect
	go.etcd.io/etcd/api/v3 v3.5.4
	go.etcd.io/etcd/client/v3 v3.5.4
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.19.1 // indirect
	golang.org/x/arch v0.3.0 // indirect
	golang.org/x/crypto v0.8.0 // indirect
	golang.org/x/net v0.9.0
	google.golang.org/genproto v0.0.0-20211028162531-8db9c33dc351 // indirect
	google.golang.org/grpc v1.45.0
	google.golang.org/protobuf v1.30.0
	gopkg.in/yaml.v3 v3.0.1
	gorm.io/driver/mysql v1.3.5
	gorm.io/gorm v1.23.8
)

replace github.com/Shopify/sarama => github.com/Shopify/sarama v1.29.0
