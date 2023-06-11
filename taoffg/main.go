package main

import (
	"xiyu.com/common"
)

func main() {
	conf := "../config/tao.yaml"
	baInfra := common.BaseInit(conf)
	baInfra.Logger.Info("Hello go")
}
