package utils

import (
	"math/rand"
	"strconv"
	"time"
)

// 生成operationID
func OperationIDGenerator() string {
	return strconv.FormatInt(time.Now().UnixNano()+int64(rand.Uint32()), 10)
}
