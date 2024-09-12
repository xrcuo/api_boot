package web

import (
	"sync"
	"time"
)

// 使用 sync.Map 存储每个接口的网络速度信息和进程信息缓存
var (
	speedCache   = sync.Map{} // 使用 sync.Map 存储网络速度信息
	processCache = sync.Map{} // 使用 sync.Map 缓存进程信息
	systemInfo   SystemInfo   // 存储系统信息
	infoMutex    sync.RWMutex // 使用读写锁保护 systemInfo
)

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

// 定义一个自定义错误类型
type systemInfoError struct {
	message string
	cause   error
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
