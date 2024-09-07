package web

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/nyancatda/AyaLog"
	"github.com/shirou/gopsutil/v4/process"
)

// 获取进程信息
func getProcessInfo() ([]*ProcessInfo, error) {
	defer func() {
		if r := recover(); r != nil {
			AyaLog.Info("System", "获取进程信息时发生错误:", r)
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
			AyaLog.Info("System", "获取进程 %d 的内存信息失败: %s\n", p.Pid, err)
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
