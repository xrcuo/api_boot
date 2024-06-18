package ping

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Get(ip string) (H ZW, err error) {
	req, _ := http.NewRequest("GET", "https://iuxcn.cn/Api/Ping?format=json&ip="+ip, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	res, _ := http.DefaultClient.Do(req)
	if res.StatusCode != 200 {
		err = fmt.Errorf("Ping code: %d", res.StatusCode)
		return
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var M Ping
	json.Unmarshal([]byte(body), &M)
	H.Code = M.Code
	H.Msg = M.Msg
	H.Node = M.Data.Node
	H.Ip = M.Data.Ip
	H.Host = M.Data.Host
	H.Domain_ip = M.Data.Domain_ip
	H.Ping_min = M.Data.Ping_min
	H.Ping_avg = M.Data.Ping_avg
	H.Ping_max = M.Data.Ping_max
	H.Location = M.Data.Location
	return H, err
}

// Ping struct
type (
	Ping struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Ip        string `json:"ip"`
			Node      string `json:"node"`
			Host      string `json:"host"`
			Domain_ip string `json:"domain_ip"`
			Ping_min  string `json:"ping_min"`
			Ping_avg  string `json:"ping_avg"`
			Ping_max  string `json:"Ping_max"`
			Location  string `json:"location"`
		} `json:"data"`
	}

	ZW struct {
		Code      string
		Msg       string
		Node      string
		Ip        string
		Host      string
		Domain_ip string
		Ping_min  string
		Ping_avg  string
		Ping_max  string
		Location  string
	}
)
