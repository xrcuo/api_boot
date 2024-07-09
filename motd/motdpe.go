package motd

import (
	"encoding/hex"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/xrcuo/api_boot/YiYan"
	"github.com/xrcuo/api_boot/ip"
)

// MotdBE信息
type MotdBEInfo struct {
	Status    string `json:"status"`     //服务器状态
	Host      string `json:"host"`       //服务器Host
	Motd      string `json:"motd"`       //Motd信息
	Agreement int    `json:"agreement"`  //协议版本
	Version   string `json:"version"`    //支持的游戏版本
	Online    int    `json:"online"`     //在线人数
	Max       int    `json:"max"`        //最大在线人数
	LevelName string `json:"level_name"` //存档名字
	GameMode  string `json:"gamemode"`   //游戏模式
	Delay     int64  `json:"delay"`      //连接延迟
	Ip        string `json:"ip"`         //服务器IP
	Area      string `json:"area"`       //服务器区域
	Text      string `json:"text"`       //随机一言
	Ipv4      string `json:"ipv4"`       //服务器IPv4
	Ipv6      string `json:"ipv6"`       //服务器IPv6
}

// MotdBE 获取BE服务器信息
// Host 服务器Host
// 返回值：服务器信息，错误信息
func MotdBE(Host string) (*MotdBEInfo, error) {
	if Host == "" {
		MotdInfo := &MotdBEInfo{
			Status: "offline",
		}
		return MotdInfo, nil
	}

	// 创建连接
	socket, err := net.Dial("udp", Host)
	if err != nil {
		MotdInfo := &MotdBEInfo{
			Status: "offline",
		}
		return MotdInfo, err
	}
	defer socket.Close()
	// 发送数据
	time1 := time.Now().UnixNano() / 1e6 //记录发送时间
	senddata, _ := hex.DecodeString("0100000000240D12D300FFFF00FEFEFEFEFDFDFDFD12345678")
	_, err = socket.Write(senddata)
	if err != nil {
		MotdInfo := &MotdBEInfo{
			Status: "offline",
		}
		return MotdInfo, err
	}
	// 接收数据
	UDPdata := make([]byte, 4096)
	socket.SetReadDeadline(time.Now().Add(5 * time.Second)) //设置读取五秒超时
	_, err = socket.Read(UDPdata)
	if err != nil {
		MotdInfo := &MotdBEInfo{
			Status: "offline",
		}
		return MotdInfo, err
	}
	time2 := time.Now().UnixNano() / 1e6 //记录接收时间
	//解析数据
	if err == nil {
		MotdData := strings.Split(string(UDPdata), ";")
		Agreement, _ := strconv.Atoi(MotdData[2])
		Online, _ := strconv.Atoi(MotdData[4])
		Max, _ := strconv.Atoi(MotdData[5])
		Z, _ := ip.GetIP(Host)
		W, _ := YiYan.GetYiYan()
		MotdInfo := &MotdBEInfo{
			Status:    "online",
			Host:      Host,
			Motd:      MotdData[1],
			Agreement: Agreement,
			Version:   MotdData[3],
			Online:    Online,
			Max:       Max,
			LevelName: MotdData[7],
			GameMode:  MotdData[8],
			Delay:     time2 - time1,
			Ip:        Z.Ip,
			Area:      Z.Area,
			Text:      W.Text,
			Ipv4:      MotdData[10],
			Ipv6:      MotdData[11],
		}
		return MotdInfo, nil
	}

	MotdInfo := &MotdBEInfo{
		Status: "offline",
	}
	return MotdInfo, err
}
