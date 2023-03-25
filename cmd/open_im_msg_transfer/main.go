package main

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"flag"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	flag.Parse()
	log.NewPrivateLog(constant.LogFileName)
}
