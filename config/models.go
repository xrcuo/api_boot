package con

// 定义一个Config结构体，包含SpeedUnit、UpdateDelay、Cacheduration和Webui四个字段
type (
	Config struct {
		// 速度单位
		Speedunit string
		// 更新延迟
		Updatedelay int
		// 缓存持续时间
		Cacheduration int
		// Webui配置
		Webui struct {
			// Webui标志
			Flags string
			// Webui主机
			Host string
			// Webui端口
			Port int
		}
	}
)
