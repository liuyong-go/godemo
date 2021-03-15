package main

import (
	"fmt"

	"github.com/liuyong-go/godemo/pkg/conf"
)

func main() {
	conf.ConfPath = "/Users/liuyong/go/src/godemo/demos/toml/"
	conf.InitConfig()
	fmt.Println(conf.Conf.DB)
	fmt.Println(conf.Conf.Server.Http.Addr)
}
