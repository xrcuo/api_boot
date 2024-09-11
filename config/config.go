package con

import (
	_ "embed"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

//go:embed default_config.yaml
var defConfig string

var (
	Conf     *Config
	confLock sync.RWMutex
	confData []byte
)

// 生成配置文件
func genConfig() error {
	sb := strings.Builder{}
	sb.WriteString(defConfig)
	err := os.WriteFile("webui.yaml", []byte(sb.String()), 0644)
	if err != nil {
		return err
	}
	return nil
}

// 解析配置文件
func Parse() {
	if err := loadConfig(); err != nil {
		logrus.Fatal("初始化配置文件失败:", err)
	}

	// 启动配置文件监听
	go watchConfig()
}

// 加载配置文件
func loadConfig() error {
	content, err := os.ReadFile("webui.yaml")
	if err != nil {
		if os.IsNotExist(err) {
			err = genConfig()
			if err != nil {
				return err
			}
			logrus.Warn("未检测到 webui.yaml，已自动于同目录生成，请配置并重新启动")
			logrus.Warn("将于 5 秒后退出...")
			time.Sleep(5 * time.Second)
			os.Exit(-1)
		}
		return err
	}

	confLock.Lock()
	defer confLock.Unlock()

	confData = content
	tmpConf := &Config{}
	if err := yaml.Unmarshal(confData, tmpConf); err != nil {
		return err
	}

	// 配置校验
	if err := validateConfig(tmpConf); err != nil {
		return err
	}

	Conf = tmpConf
	logrus.Info("配置文件已加载")
	return nil
}

// 监听配置文件变化
func watchConfig() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logrus.Fatal("创建文件监听器失败:", err)
	}
	defer watcher.Close()

	err = watcher.Add("webui.yaml")
	if err != nil {
		logrus.Fatal("添加文件监听失败:", err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				logrus.Info("配置文件已修改，正在重新加载...")
				if err := loadConfig(); err != nil {
					logrus.Error("重新加载配置文件失败:", err)
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			logrus.Error("文件监听错误:", err)
		}
	}
}

// 校验配置
func validateConfig(conf *Config) error {
	// TODO: 添加配置校验逻辑
	if conf == nil {
		return nil
	}
	return nil
}
