package con

import (
	_ "embed"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

//go:embed default_config.yaml
var defConfig string

var Conf *Config

// 生成配置文件
func genConfig() error {
	sb := strings.Builder{}
	sb.WriteString(defConfig)
	err := os.WriteFile("boot.yaml", []byte(sb.String()), 0644)
	if err != nil {
		return err
	}
	return nil
}

// 解析配置文件
func Parse() {
	content, err := os.ReadFile("boot.yaml")
	if err != nil {
		err = genConfig()
		if err != nil {
			panic("无法生成设置文件: boot.yaml, 请确认是否给足系统权限")
		}
		logrus.Warn("未检测到 boot.yaml，已自动于同目录生成，请配置并重新启动")
		logrus.Warn("将于 5 秒后退出...")
		os.Exit(-1)
	}

	Conf = &Config{}
	err = yaml.Unmarshal(content, Conf)
	if err != nil {
		logrus.Fatal("解析 boot.yaml 失败，请检查格式、内容是否输入正确")
	}
}
