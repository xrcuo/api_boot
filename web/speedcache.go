package web

import (
	"time"

	con "github.com/xrcuo/api_boot/config"
)

// 清理缓存
func cleanCache() {
	var ca = time.Duration(con.Conf.Cacheduration) * time.Second
	// 清理 speedCache
	speedCache.Range(func(key, value interface{}) bool {
		// 检查缓存项是否过期
		if time.Since(value.(NetworkSpeed).LastUpdated) > ca {
			speedCache.Delete(key)
		}
		return true
	})

	// 清理 processCache
	processCache.Range(func(key, value interface{}) bool {
		// 检查缓存项是否过期
		if time.Since(value.(*ProcessInfo).LastUpdated) > ca {
			processCache.Delete(key)
		}
		return true
	})
}
