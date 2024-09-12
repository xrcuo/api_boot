package webw

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"runtime/debug"
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

//go:embed static
var staticFiles embed.FS

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
	Name        string
	BytesRecv   uint64
	BytesSent   uint64
	LastUpdated time.Time `json:"-"` // 添加 LastUpdated 字段
}

// 为 NetworkSpeed 添加 GetLastUpdated 方法
func (ns NetworkSpeed) GetLastUpdated() time.Time {
	return ns.LastUpdated
}

// 定义一个结构体存储进程信息
type ProcessInfo struct {
	Pid         int32     `json:"pid"`
	Name        string    `json:"name"`
	CPUPercent  float64   `json:"cpuPercent"`
	MemoryUsed  float32   `json:"memoryUsed"` // 内存使用量，单位：MB
	LastUpdated time.Time `json:"-"`          // 添加 LastUpdated 字段
}

// 为 ProcessInfo 添加 GetLastUpdated 方法
func (pi ProcessInfo) GetLastUpdated() time.Time {
	return pi.LastUpdated
}

// 使用 sync.Map 存储每个接口的网络速度信息和进程信息缓存
var (
	speedCache    = sync.Map{} // 使用 sync.Map 存储网络速度信息
	processCache  = sync.Map{} // 使用 sync.Map 缓存进程信息
	systemInfo    SystemInfo
	infoMutex     sync.RWMutex
	speedUnit     = "kbps" // 默认速度单位
	updateDelay   = 1 * time.Second
	cacheDuration = 5 * time.Second // 缓存持续时间
)

func Ltml() {
	router := gin.Default()
	static, _ := fs.Sub(staticFiles, "static")
	router.StaticFS("/static", http.FS(static))

	templ := template.Must(template.ParseFS(templates, "templates/*"))
	router.SetHTMLTemplate(templ)

	// 启动 goroutine 定时更新系统信息和网络速度信息
	go func() {
		ticker := time.NewTicker(updateDelay)
		defer ticker.Stop()

		for range ticker.C {
			updateSystemInfo()
			updateNetworkSpeed()
		}
	}()

	// 启动 goroutine 定时清理缓存
	go func() {
		ticker := time.NewTicker(cacheDuration)
		defer ticker.Stop()

		for range ticker.C {
			cleanSpeedCache(&speedCache, cacheDuration)
			cleanProcessCache(&processCache, cacheDuration)
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

// 定时更新系统信息
func updateSystemInfo() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("更新系统信息时发生错误:", r)
			debug.PrintStack()
		}
	}()

	cpuPercent, _ := cpu.Percent(0, false)
	memory, _ := mem.VirtualMemory()

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

	infoMutex.Lock()
	systemInfo.CPUPercent = cpuPercent[0]
	systemInfo.MemoryUsed = memory.UsedPercent
	systemInfo.DiskInfos = diskInfos
	infoMutex.Unlock()
}

// 获取之前的网络接口统计信息
func getPreviousNetworkSpeed(counterName string) (NetworkSpeed, bool) {
	prevCounter, ok := speedCache.Load(counterName)
	if !ok {
		return NetworkSpeed{}, false
	}

	prevSpeedInfo, ok := prevCounter.(NetworkSpeed)
	if !ok {
		fmt.Printf("类型断言失败: %s\n", counterName)
		return NetworkSpeed{}, false
	}

	return prevSpeedInfo, true
}

// 定时更新网络速度信息
func updateNetworkSpeed() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("更新网络速度信息时发生错误:", r)
			debug.PrintStack()
		}
	}()

	currCounters, err := net.IOCounters(true)
	if err != nil {
		fmt.Println("获取网络接口统计信息失败:", err)
		return
	}

	for _, counter := range currCounters {
		// 从 speedCache 中获取现有的 NetworkSpeed 对象，如果不存在则创建
		speed, _ := speedCache.LoadOrStore(counter.Name, NetworkSpeed{
			Name:        counter.Name,
			LastUpdated: time.Now(),
		})

		// 使用类型断言获取 NetworkSpeed 对象
		speedInfo, ok := speed.(NetworkSpeed)
		if !ok {
			fmt.Printf("类型断言失败: %s\n", counter.Name)
			continue
		}

		// 获取之前的网络接口统计信息
		prevSpeedInfo, ok := getPreviousNetworkSpeed(counter.Name)
		if !ok {
			continue
		}

		// 计算速度差值
		bytesRecv := counter.BytesRecv - prevSpeedInfo.BytesRecv
		bytesSent := counter.BytesSent - prevSpeedInfo.BytesSent

		// 转换为 kbps
		kbpsRecv := bytesRecv * 8 / 1024
		kbpsSent := bytesSent * 8 / 1024

		// 更新速度值
		speedInfo.BytesRecv = kbpsRecv
		speedInfo.BytesSent = kbpsSent
		speedInfo.LastUpdated = time.Now()

		// 更新 speedCache
		speedCache.Store(counter.Name, speedInfo)
	}
}

// 获取进程信息
func getProcessInfo() ([]*ProcessInfo, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("获取进程信息时发生错误:", r)
			debug.PrintStack()
		}
	}()

	processes, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("获取进程列表失败: %w", err)
	}

	processInfos := make([]*ProcessInfo, 0, len(processes))
	for _, p := range processes {
		// 检查进程信息是否已缓存
		if cachedInfo, ok := processCache.Load(p.Pid); ok {
			// 使用类型断言获取缓存的 ProcessInfo 对象
			processInfo, ok := cachedInfo.(*ProcessInfo)
			if ok {
				processInfos = append(processInfos, processInfo)
				continue // 跳过后续获取进程信息的步骤
			} else {
				// 如果类型断言失败，则从缓存中删除该条目
				processCache.Delete(p.Pid)
			}
		}

		// 获取进程信息
		cpuPercent, _ := p.CPUPercent()
		memInfo, err := p.MemoryInfo()
		if err != nil {
			fmt.Printf("获取进程 %d 的内存信息失败: %s\n", p.Pid, err)
			continue
		}
		name, _ := p.Name()

		processInfo := &ProcessInfo{
			Pid:         p.Pid,
			Name:        name,
			CPUPercent:  cpuPercent,
			MemoryUsed:  float32(memInfo.RSS) / (1024 * 1024), // 转换为 MB
			LastUpdated: time.Now(),
		}

		// 将进程信息添加到结果切片
		processInfos = append(processInfos, processInfo)

		// 将进程信息缓存起来
		processCache.Store(p.Pid, processInfo)
	}

	return processInfos, nil
}

// 清理 speedCache
func cleanSpeedCache(cache *sync.Map, cacheDuration time.Duration) {
	cache.Range(func(key, value interface{}) bool {
		speed, ok := value.(NetworkSpeed)
		if !ok {
			return true // 跳过非 NetworkSpeed 类型的条目
		}

		if time.Since(speed.GetLastUpdated()) > cacheDuration {
			cache.Delete(key)
		}
		return true
	})
}

// 清理 processCache
func cleanProcessCache(cache *sync.Map, cacheDuration time.Duration) {
	cache.Range(func(key, value interface{}) bool {
		process, ok := value.(ProcessInfo)
		if !ok {
			return true // 跳过非 ProcessInfo 类型的条目
		}

		if time.Since(process.GetLastUpdated()) > cacheDuration {
			cache.Delete(key)
		}
		return true
	})
}
