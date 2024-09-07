package web

import (
	"runtime/debug"

	"github.com/nyancatda/AyaLog"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

// 定时更新系统信息
func updateSystemInfo() {
	defer func() {
		if r := recover(); r != nil {
			AyaLog.Info("System", "更新系统信息时发生错误:", r)
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
