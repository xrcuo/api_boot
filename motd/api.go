package motd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetBdsInfo(host string) (H ZW, err error) {
	req, _ := http.NewRequest("GET", "https://mot-api.aqco.top/api/bds?host="+host, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	res, _ := http.DefaultClient.Do(req)
	if res.StatusCode != 200 {
		err = fmt.Errorf("bds code: %d", res.StatusCode)
		return
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var M MotdBEInfo
	json.Unmarshal([]byte(body), &M)
	H.Status = M.Status
	H.Host = M.Host
	H.Motd = M.Motd
	H.Agreement = M.Agreement
	H.Version = M.Version
	H.Online = M.Online
	H.Max = M.Max
	H.LevelName = M.LevelName
	H.GameMode = M.GameMode
	H.Delay = M.Delay
	H.Ip = M.Ip
	H.Area = M.Area
	H.Text = M.Text
	return H, err
}

type ZW struct {
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
	Ip        string //服务器IP
	Area      string //服务器区域
	Text      string //服务器文字信息
}
