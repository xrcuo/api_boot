package YiYan

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetYiYan() (H ZW, err error) {
	req, _ := http.NewRequest("GET", "https://iuxcn.cn/Api/YiYan?format=json", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	res, _ := http.DefaultClient.Do(req)
	if res.StatusCode != 200 {
		err = fmt.Errorf("YiYan API Error: %d", res.StatusCode)
		return
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var M Yiyan
	json.Unmarshal([]byte(body), &M)
	H.Code = M.Code
	H.Msg = M.Msg
	H.Text = M.Data.Text
	return H, err
}

type (
	Yiyan struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Text string `json:"text"`
		} `json:"data"`
	}

	ZW struct {
		Code string
		Msg  string
		Text string
	}
)
