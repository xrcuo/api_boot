package web

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
)

//go:embed templates
var templates embed.FS

// 定义一个结构体存储系统信息
type SystemInfo struct {
	CPUPercent float64 `json:"cpuPercent"`
	MemoryUsed float64 `json:"memoryUsed"`
	DiskUsed   float64 `json:"diskUsed"`
}

// 定义一个结构体存储网络速度信息
type NetworkSpeed struct {
	Name      string
	BytesRecv uint64
	BytesSent uint64
}

// 使用 sync.Map 存储每个接口的网络速度信息
var (
	speedCache  = make(map[string]NetworkSpeed)
	systemInfo  SystemInfo
	infoMutex   sync.RWMutex
	speedUnit   = "kbps" // 默认速度单位
	updateDelay = 1 * time.Second
)

func Ltml() {
	router := gin.Default()

	templ := template.Must(template.ParseFS(templates, "templates/*"))
	//htmlTempl := template.Must(template.New("html").ParseFS(templates, "templates/index.html"))
	router.SetHTMLTemplate(templ)

	// 启动 goroutine 定时更新系统信息和网络速度信息
	go updateSystemInfo()
	go updateNetworkSpeed()

	router.GET("/", func(c *gin.Context) {
		// 获取网络接口信息
		interfaces, err := net.Interfaces()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"Interfaces": interfaces,
			"SpeedUnit":  speedUnit,
		})
	})

	router.GET("/system_info", func(c *gin.Context) {
		infoMutex.RLock()
		defer infoMutex.RUnlock()
		c.JSON(http.StatusOK, systemInfo)
	})

	router.GET("/speed", func(c *gin.Context) {
		speeds := make([]NetworkSpeed, 0, len(speedCache))
		for _, speed := range speedCache {
			speeds = append(speeds, speed)
		}
		c.JSON(http.StatusOK, speeds)
	})

	router.Run(":8080")
}

// 定时更新系统信息
func updateSystemInfo() {
	for {
		// 获取 CPU 使用率
		cpuPercent, _ := cpu.Percent(0, false)

		// 获取内存信息
		memory, _ := mem.VirtualMemory()

		// 获取硬盘使用率
		partitions, _ := disk.Partitions(false)
		for _, partition := range partitions {
			if partition.Mountpoint == "/" { // 获取根分区的使用率
				usage, _ := disk.Usage(partition.Mountpoint)
				infoMutex.Lock()
				systemInfo.DiskUsed = usage.UsedPercent
				infoMutex.Unlock()
				break
			}
		}

		// 更新系统信息
		infoMutex.Lock()
		systemInfo.CPUPercent = cpuPercent[0]
		systemInfo.MemoryUsed = memory.UsedPercent
		infoMutex.Unlock()

		time.Sleep(updateDelay)
	}
}

// 定时更新网络速度信息
func updateNetworkSpeed() {
	ticker := time.NewTicker(updateDelay)
	defer ticker.Stop()

	// 获取初始网络接口统计信息
	prevCounters, _ := net.IOCounters(true)

	for range ticker.C {
		// 获取当前网络接口统计信息
		currCounters, err := net.IOCounters(true)
		if err != nil {
			fmt.Println("获取网络接口统计信息失败:", err)
			continue
		}

		// 遍历网络接口，更新速度信息
		for i, counter := range currCounters {
			if i < len(prevCounters) {
				// 计算速度差值
				bytesRecv := counter.BytesRecv - prevCounters[i].BytesRecv
				bytesSent := counter.BytesSent - prevCounters[i].BytesSent

				// 更新速度缓存
				speedCache[counter.Name] = NetworkSpeed{
					Name:      counter.Name,
					BytesRecv: bytesRecv,
					BytesSent: bytesSent,
				}
			}
		}

		// 更新之前的网络接口统计信息
		prevCounters = currCounters
	}
}
