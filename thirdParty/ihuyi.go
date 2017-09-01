//接口类型：互亿无线语音通知接口。
//账户注册：请通过该地址开通账户http://sms.ihuyi.com/register.html
//注意事项：
//（1）调试期间，请仔细阅读接口文档；
//（2）请使用APIID（查看APIID请登录用户中心->语音通知->帐户及签名设置->APIID）及 APIkey来调用接口；
//（3）该代码仅供接入互亿无线语音通知接口参考使用，客户可根据实际需要自行编写；

package thirdParty

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var ihy AuthInfo

func InitAuthInfo() {

	ihy.vendor = "ihuyi"
	err := LoadJsonConf(&ihy)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(ihy)
}

func PostMsgToIhuyi(mobile string) {
	v := url.Values{}
	_now := strconv.FormatInt(time.Now().Unix(), 10)
	//fmt.Printf(_now)
	_account := ihy.GetSid()
	_password := ihy.GetToken()
	_mobile := mobile
	_content := "您的订单号是：0648。已由顺风快递发出，请注意查收。"
	//_content := "生产系统已发生严重故障，请紧急处理。"
	v.Set("account", _account)
	v.Set("password", GetMd5String(_account+_password+_mobile+_content+_now))
	v.Set("mobile", _mobile)
	v.Set("content", _content)
	v.Set("time", _now)
	body := ioutil.NopCloser(strings.NewReader(v.Encode())) //form数据编码
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://api.vm.ihuyi.com/webservice/voice.php?method=Submit&format=json", body)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	resp, err := client.Do(req) //发送
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(data), err)
}

func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
