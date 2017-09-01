// 供应商：云之讯
// 请求方式：POST
// 请求地址：https://message.ucpaas.com/{version}/Accounts/{accountSid}/Calls/voiceNotify?sig={SigParameter}
// 参数说明：
// SigParameter是REST API 验证参数
// URL后必须带有sig参数，sig= MD5（账户Id + 账户授权令牌 + 时间戳），共32位(注:转成大写)
// 使用MD5加密（账户Id + 账户授权令牌 + 时间戳），共32位
// 时间戳是当前系统时间（24小时制），格式“yyyyMMddHHmmss”。时间戳有效时间为50分钟。
// Authorization是包头验证信息
// 使用Base64编码（账户Id + 冒号 + 时间戳）
// 冒号为英文冒号
// 时间戳是当前系统时间（24小时制），格式“yyyyMMddHHmmss”，需与SigParameter中时间戳相同

package thirdParty

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type UcpaasVoiceNotify struct {
	AppId      string `json:"appId"`
	To         string `json:"to"`
	ToSerNum   string `json:"toSerNum"`
	Type       string `json:"type"`
	PlayTimes  string `json:"playTimes"`
	TemplateId string `json:"templateId"`
	Content    string `json:"content"`
	BillUrl    string `json:"billUrl"`
	UserData   string `json:"userData"`
}

type UcpaasVoiceMsg struct {
	Vn UcpaasVoiceNotify `json:"voiceNotify"`
}

type UcpaasVendorInfo struct {
	httpHost string
	vn       UcpaasVoiceMsg
	auth     AuthInfo
}

var uvi UcpaasVendorInfo = UcpaasVendorInfo{}

//init all vendor info value like sid/token/appid ..
func InitUcpassAndNewReqJsonBody(phone, content string) []byte {

	uvi.httpHost = "message.ucpaas.com"
	uvi.auth.vendor = "ucpaas"

	err := LoadJsonConf(&uvi.auth)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(uvi.auth)

	uvn := UcpaasVoiceNotify{
		AppId:      uvi.auth.appId,
		To:         phone,
		Type:       "2",
		Content:    content,
		PlayTimes:  "2",
		TemplateId: uvi.auth.templateId,
	}
	uvi.vn = UcpaasVoiceMsg{Vn: uvn}

	js, err := json.Marshal(uvi.vn)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return js

}

func PostMsgToUcpaas(body []byte) ([]byte, error) {

	/* 账户(sid) + 授权令牌(token) + 时间戳 */
	sign := uvi.auth.GetSid() + uvi.auth.GetToken() + uvi.auth.GetTimeSec()
	h := md5.New()
	h.Write([]byte(sign))

	SigParameter := hex.EncodeToString(h.Sum(nil))

	sig := strings.ToUpper(SigParameter)

	fmt.Println(string(body))
	return HttpsPostRequstToUcpass(sig, body)
}

func HttpsPostRequstToUcpass(sig string, body []byte) ([]byte, error) {

	uri := "https://" + uvi.httpHost + `/2014-06-30/Accounts/` + uvi.auth.GetSid() + `/Calls/voiceNotify?sig=` + sig

	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}

	fmt.Println(uri)
	//fmt.Println(string(body))

	reqBody := ioutil.NopCloser(strings.NewReader(string(body)))
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("POST", uri, reqBody)

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json;charset=utf-8;")
	req.Header.Set("Connection", "close")

	/* Authorization域  使用Base64编码（账户Id + 冒号 + 时间戳）(time.Now().Format("20060102150405"))*/
	auths := uvi.auth.GetSid() + ":" + uvi.auth.GetTimeSec()

	b64 := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
	b64Auth := b64.EncodeToString([]byte(auths))
	req.Header.Set("Authorization", b64Auth)

	req.Header.Set("Content-Length", strconv.Itoa(len(body)))

	req.Header.Write(os.Stdout)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("http client failed:", err)
		return nil, errors.New("http request failed!")
	}

	defer resp.Body.Close()

	resbody, err := ioutil.ReadAll(resp.Body)
	return resbody, err
}
