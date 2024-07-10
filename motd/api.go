package motd

import (
	"encoding/json"
	"fmt"
)

// MCSAF函数用于获取服务器信息
func MCSAF(cb string) (C MC, err error) {
	// 调用MotdBE函数获取服务器信息
	data, err := MotdBE(cb)
	if err != nil {
		fmt.Println("获取服务器信息失败:", err)
		return
	}
	// 将服务器信息转换为json格式
	jsonData, err := json.Marshal(data)

	if err != nil {
		fmt.Println("json转换失败:", err)
		return
	}
	// 定义MotdBEInfo结构体
	var M MotdBEInfo
	// 将json数据转换为MotdBEInfo结构体
	json.Unmarshal([]byte(jsonData), &M)
	// 将MotdBEInfo结构体的数据赋值给MC结构体
	C.Status = M.Status
	C.Host = M.Host
	C.Motd = M.Motd
	C.Agreement = M.Agreement
	C.Version = M.Version
	C.Online = M.Online
	C.Max = M.Max
	C.LevelName = M.LevelName
	C.GameMode = M.GameMode
	C.Delay = M.Delay
	C.Ip = M.Ip
	C.Area = M.Area
	C.Text = M.Text
	C.Ipv4 = M.Ipv4
	C.Ipv6 = M.Ipv6
	return C, err

}

// MC结构体用于存储服务器信息
type MC struct {
	Status    string //服务器状态
	Host      string //服务器Host
	Motd      string //Motd信息
	Agreement int    //协议版本
	Version   string //支持的游戏版本
	Online    int    //在线人数
	Max       int    //最大在线人数
	LevelName string //存档名字
	GameMode  string //游戏模式
	Delay     int64  //连接延迟
	Favicon   string //服务器图标(base64)
	Ip        string //服务器IP
	Area      string //服务器区域
	Text      string //服务器文字信息
	Ipv4      string //服务器IPV4
	Ipv6      string //服务器IPV6
}
