// 供应商：云讯
//请求方式：POST
//沙箱环境http://sandbox.ytx.net
//正式环境http://api.ytx.net
//签名sign:
//MD5加密（账户Id + 账户授权令牌 +时间戳)，
//例如：Sign=AAABBBCCCDDDEEEFFFGGG *时间戳需与Authorization中时间戳相同(时间戳格式:yyyyMMddHHmmss) 注:MD5加密32位,无论大小写
//Authorization:
// 通信`平台API接口，包头验证信息：base64加密(账户Id + "|" + 时间戳)
// 说明：时间戳有效时间为24小时 格式"yyyyMMddHHmmss"，如：20140416142030

package thirdParty

import (
	"crypto/md5"
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

type YtxVoiceNotifyParam struct {
	Action     string   `json:"action"`
	Dst        string   `json:"dst"`
	AppId      string   `json:"appid"`
	TemplateId string   `json:"templateId"`
	Datas      []string `json:"datas"`
}

type YtxVendorInfo struct {
	httpHost string
	vn       YtxVoiceNotifyParam
	auth     AuthInfo
}

var vinfo YtxVendorInfo = YtxVendorInfo{}

//init all vendor info value like sid/token/appid ..
func InitYtxAndNewReqJsonBody(phone string) []byte {

	vinfo.httpHost = "sandbox.ytx.net"
	vinfo.auth.vendor = "ytx"

	err := LoadJsonConf(&vinfo.auth)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(vinfo.auth)

	vinfo.auth.SetSid(vinfo.auth.sid)
	vinfo.auth.SetToken(vinfo.auth.token)

	vinfo.vn = YtxVoiceNotifyParam{
		AppId:      vinfo.auth.appId,
		Dst:        phone,
		Action:     "templateNoticeCall",
		TemplateId: vinfo.auth.templateId,
		Datas:      []string{"66537", "2"},
	}
	js, err := json.Marshal(vinfo.vn)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	return js
}

func PostMsgToYtx(body []byte) ([]byte, error) {

	/* 账户(sid) + 授权令牌(token) + 时间戳 */
	sign := vinfo.auth.GetSid() + vinfo.auth.GetToken() + vinfo.auth.GetTimeSec()
	h := md5.New()
	h.Write([]byte(sign))

	SigParameter := hex.EncodeToString(h.Sum(nil))

	sig := strings.ToUpper(SigParameter)

	fmt.Println(string(body))
	return HttpPostRequstToYtx(sig, body)
}

func HttpPostRequstToYtx(sig string, body []byte) ([]byte, error) {

	uri := "http://" + vinfo.httpHost + `/201512/sid/` + vinfo.auth.GetSid() + `/call/NoticeCall.wx?Sign=` + sig
	fmt.Println(uri)
	//fmt.Println(string(body))

	reqBody := ioutil.NopCloser(strings.NewReader(string(body)))
	client := &http.Client{}
	req, err := http.NewRequest("POST", uri, reqBody)

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json;charset=utf-8;")
	//req.Header.Set("Connection", "close")

	/* Authorization域  使用Base64编码（账户Id + "|" + 时间戳）(time.Now().Format("20060102150405"))*/
	auths := vinfo.auth.GetSid() + "|" + vinfo.auth.GetTimeSec()

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
