package web

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v4/net"
	con "github.com/xrcuo/api_boot/config"
)

//go:embed templates
var templates embed.FS

func Ltml() {
	var (
		up = time.Duration(con.Conf.Updatedelay) * time.Second
		ca = time.Duration(con.Conf.Cacheduration) * time.Second
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
			"SpeedUnit":  con.Conf.Speedunit,
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
	router.Run(":8080")
}
