package main

import (
	"github.com/nyancatda/AyaLog"
	con "github.com/xrcuo/api_boot/config"
	"github.com/xrcuo/api_boot/web"
)

func main() {
	web.Ltml()

}

var conf *con.Config

func init() {
	// 加载配置文件
	con.Parse()
	AyaLog.Info("System", "config.yaml 加载成功")
	conf = con.Conf
}
