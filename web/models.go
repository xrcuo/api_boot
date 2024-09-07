package web

import "time"

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

// 定义一个结构体存储进程信息
type ProcessInfo struct {
	Pid         int32     `json:"pid"`
	Name        string    `json:"name"`
	CPUPercent  float64   `json:"cpuPercent"`
	MemoryUsed  float32   `json:"memoryUsed"` // 内存使用量，单位：MB
	LastUpdated time.Time `json:"-"`          // 添加 LastUpdated 字段
}
