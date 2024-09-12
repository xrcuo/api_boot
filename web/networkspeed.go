package web

import (
	"runtime/debug"
	"time"

	"github.com/nyancatda/AyaLog"
	"github.com/shirou/gopsutil/v4/net"
)

// 定时更新网络速度信息
func updateNetworkSpeed() {
	defer func() {
		if r := recover(); r != nil {
			AyaLog.Warning("更新网络速度信息时发生错误:", r)
			debug.PrintStack()
		}
	}()

	currCounters, err := net.IOCounters(true)
	if err != nil {
		AyaLog.Warning("获取网络接口统计信息失败:", err)
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
			AyaLog.Warning("类型断言失败: %s\n", counter.Name)
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

// 获取之前的网络接口统计信息
func getPreviousNetworkSpeed(counterName string) (NetworkSpeed, bool) {
	prevCounter, ok := speedCache.Load(counterName)
	if !ok {
		return NetworkSpeed{}, false
	}

	prevSpeedInfo, ok := prevCounter.(NetworkSpeed)
	if !ok {
		AyaLog.Warning("类型断言失败: 预期类型 NetworkSpeed，实际类型 %T\n", prevCounter)
		return NetworkSpeed{}, false
	}

	return prevSpeedInfo, true
}
