package web

import (
	"reflect"
	"sync"
	"time"

	"github.com/nyancatda/AyaLog"
)

// 清理缓存 (使用泛型函数)
func cleanCache[T any](cache *sync.Map, cacheDuration time.Duration) {
	cache.Range(func(key, value interface{}) bool {
		item, ok := value.(T)
		if !ok {
			return true // 跳过非预期类型的条目
		}

		// 使用反射获取 GetLastUpdated 方法
		method := reflect.ValueOf(item).MethodByName("GetLastUpdated")
		if !method.IsValid() {
			return true // 如果没有该方法，则跳过
		}

		// 调用 GetLastUpdated 方法获取最后更新时间
		results := method.Call(nil)
		if len(results) == 0 {
			AyaLog.Warning("GetLastUpdated 方法没有返回值")
			return true
		}
		lastUpdated, ok := results[0].Interface().(time.Time)
		if !ok {
			AyaLog.Warning("GetLastUpdated 方法返回值类型错误")
			return true
		}

		if time.Since(lastUpdated) > cacheDuration {
			cache.Delete(key)
		}
		return true
	})
}
