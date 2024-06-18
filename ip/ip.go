package ip

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
)

func GetIP(ipx string) (H ZW, err error) {
	// 格式化ip
	s := "//" + ipx
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	// 端口
	host, _, err := net.SplitHostPort(u.Host)
	if err != nil {
		panic(err)
	}

	// 域名解析
	ips, err := net.LookupIP(host)
	if err != nil {
		fmt.Println("域名解析失败:", err)
		return
	}
	var dcip = ips[0]
	req, _ := http.NewRequest("GET", "https://iuxcn.cn/Api/IP?format=json&ip="+dcip.String(), nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	res, _ := http.DefaultClient.Do(req)
	if res.StatusCode != 200 {
		err = fmt.Errorf("IP code: %d", res.StatusCode)
		return
	}
	defer res.Body.Close()
	// 解析json
	body, _ := ioutil.ReadAll(res.Body)
	var M Ip
	json.Unmarshal([]byte(body), &M)
	H.Code = M.Code
	H.Msg = M.Msg
	H.Area = M.Data.Location.Area
	H.Ip = M.Data.Location.Ip
	return H, err
}

type (
	Ip struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Location struct {
				Ip   string `json:"ip"`
				Area string `json:"area"`
			} `json:"location"`
		} `json:"data"`
	}
	ZW struct {
		Code string
		Msg  string

		Ip   string
		Area string
	}
)
