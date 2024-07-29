package control

import (
	"github.com/sirupsen/logrus"
	con "github.com/xrcuo/api_boot/config"
)

var conf *con.Config

func init() {
	// 加载配置文件
	con.Parse()
	logrus.Info("config.yaml 加载成功")
	conf = con.Conf

}
