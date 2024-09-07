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
	"github.com/shirou/gopsutil/v4/process"
)

//go:embed templates
var templates embed.FS

// 定义一个结构体存储系统信息
type SystemInfo struct {
	CPUPercent float64     `json:"cpuPercent"`
	MemoryUsed float64     `json:"memoryUsed"`
	DiskInfos  []*DiskInfo `json:"diskInfos"`
}

// 定义一个结构体存储硬盘信息
type DiskInfo struct {
	MountPoint  string  `json:"mountPoint"`
	Total       uint64  `json:"total"` // 总容量，单位：GB
	Used        uint64  `json:"used"`  // 已用容量，单位：GB
	UsedPercent float64 `json:"usedPercent"`
}

// 定义一个结构体存储网络速度信息
type NetworkSpeed struct {
	Name      string
	BytesRecv uint64
	BytesSent uint64
}

// 定义一个结构体存储进程信息
type ProcessInfo struct {
	Pid        int32   `json:"pid"`
	Name       string  `json:"name"`
	CPUPercent float64 `json:"cpuPercent"`
	MemoryUsed float32 `json:"memoryUsed"` // 内存使用量，单位：MB
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
		infoMutex.RLock()
		defer infoMutex.RUnlock()

		// 使用更有效的方式构建 speeds 切片
		speeds := make([]NetworkSpeed, 0, len(speedCache))
		for _, speed := range speedCache {
			speeds = append(speeds, speed)
		}

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
	// 在后台启动 Web 服务器

}

// 定时更新系统信息
func updateSystemInfo() {
	for {
		// 获取 CPU 使用率
		cpuPercent, _ := cpu.Percent(0, false)

		// 获取内存信息
		memory, _ := mem.VirtualMemory()

		// 获取硬盘信息
		partitions, _ := disk.Partitions(false)
		diskInfos := make([]*DiskInfo, 0, len(partitions))
		for _, partition := range partitions {
			usage, _ := disk.Usage(partition.Mountpoint)
			diskInfos = append(diskInfos, &DiskInfo{
				MountPoint:  partition.Mountpoint,
				Total:       usage.Total / (1024 * 1024 * 1024), // GB
				Used:        usage.Used / (1024 * 1024 * 1024),  // GB
				UsedPercent: usage.UsedPercent,
			})
		}

		// 更新系统信息
		infoMutex.Lock()
		systemInfo.CPUPercent = cpuPercent[0]
		systemInfo.MemoryUsed = memory.UsedPercent
		systemInfo.DiskInfos = diskInfos
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

	// 在循环外部创建 speeds 切片
	speeds := make([]NetworkSpeed, 0, len(speedCache))

	// 限制 speedCache 的容量，例如最多存储 100 个网络接口的信息
	maxCacheSize := 100

	for range ticker.C {
		// 获取当前网络接口统计信息
		currCounters, err := net.IOCounters(true)
		if err != nil {
			fmt.Println("获取网络接口统计信息失败:", err)
			continue
		}

		// 遍历网络接口，更新速度信息
		infoMutex.Lock()
		// 重置 speeds 切片的长度
		speeds = speeds[:0]

		for i, counter := range currCounters {
			if i < len(prevCounters) {
				// 计算速度差值
				bytesRecv := counter.BytesRecv - prevCounters[i].BytesRecv
				bytesSent := counter.BytesSent - prevCounters[i].BytesSent

				// 转换为 kbps
				kbpsRecv := bytesRecv * 8 / 1024
				kbpsSent := bytesSent * 8 / 1024

				// 使用 append 方法添加速度信息
				speeds = append(speeds, NetworkSpeed{
					Name:      counter.Name,
					BytesRecv: kbpsRecv,
					BytesSent: kbpsSent,
				})
			}
		}

		// 如果 speedCache 的容量超过 maxCacheSize，则删除最旧的条目
		if len(speedCache) > maxCacheSize {
			for key := range speedCache {
				delete(speedCache, key)
				if len(speedCache) <= maxCacheSize {
					break
				}
			}
		}

		// 更新 speedCache
		for _, speed := range speeds {
			speedCache[speed.Name] = speed
		}

		infoMutex.Unlock()
		prevCounters = currCounters
	}
}

// 获取进程信息
func getProcessInfo() ([]*ProcessInfo, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("获取进程列表失败: %w", err)
	}

	processInfos := make([]*ProcessInfo, 0, len(processes))
	for _, p := range processes {
		// 使用 defer 释放进程资源
		//defer p.Release()

		cpuPercent, _ := p.CPUPercent()
		memInfo, _ := p.MemoryInfo()
		name, _ := p.Name() // 获取进程名

		processInfos = append(processInfos, &ProcessInfo{
			Pid:        p.Pid,
			Name:       name,
			CPUPercent: cpuPercent,
			MemoryUsed: float32(memInfo.RSS) / (1024 * 1024), // 转换为 MB
		})
	}

	return processInfos, nil
}
