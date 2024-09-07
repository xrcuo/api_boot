package web

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nyancatda/AyaLog"
	"github.com/shirou/gopsutil/v4/net"
	con "github.com/xrcuo/api_boot/config"
)

//go:embed templates
var templates embed.FS

var conf *con.Config

func Ltml() {
	// 加载配置文件
	con.Parse()
	AyaLog.Info("System", "webui.yaml 加载成功")
	conf = con.Conf

	// 启用 WebUI
	okk := conf.Webui.Flags
	if okk != "false" {
		AyaLog.Info("System", "Webui启动服务")
		go Start()

	} else {
		AyaLog.Info("System", "Webui停止服务")
	}

}

func Start() {
	var (
		up = time.Duration(conf.Updatedelay) * time.Second
		ca = time.Duration(conf.Cacheduration) * time.Second
		//omain = fmt.Sprintf("%s:%d", conf.Webui.Host, conf.Webui.Port)
	)
	router := gin.Default()

	templ := template.Must(template.ParseFS(templates, "templates/*"))
	router.SetHTMLTemplate(templ)

	// 启动 goroutine 定时更新系统信息和网络速度信息
	go func() {
		ticker := time.NewTicker(up)
		defer ticker.Stop()

		for range ticker.C {
			updateSystemInfo()
			updateNetworkSpeed()
		}
	}()

	// 启动 goroutine 定时清理缓存
	go func() {
		ticker := time.NewTicker(ca)
		defer ticker.Stop()

		for range ticker.C {
			cleanCache()
		}
	}()

	router.GET("/", func(c *gin.Context) {
		// 获取网络接口信息
		interfaces, err := net.Interfaces()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"Interfaces": interfaces,
			"SpeedUnit":  conf.Speedunit,
		})
	})

	router.GET("/system_info", func(c *gin.Context) {
		infoMutex.RLock()
		defer infoMutex.RUnlock()
		c.JSON(http.StatusOK, systemInfo)
	})

	router.GET("/speed", func(c *gin.Context) {
		infoMutex.RLock()
		defer infoMutex.RUnlock()

		speeds := make([]NetworkSpeed, 0)
		speedCache.Range(func(key, value interface{}) bool {
			speed, ok := value.(NetworkSpeed)
			if ok {
				speeds = append(speeds, speed)
			}
			return true
		})

		c.JSON(http.StatusOK, speeds)
	})

	router.GET("/processes", func(c *gin.Context) {
		processInfos, err := getProcessInfo()
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("获取进程信息失败: %s", err))
			return
		}
		c.JSON(http.StatusOK, processInfos)
	})
	// 在后台启动 Web 服务器
	go func() {
		if err := router.Run(conf.Webui.Host + ":" + strconv.Itoa(conf.Webui.Port)); err != nil {
			AyaLog.Error("System", err)
		}
	}()
}
