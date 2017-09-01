package thirdParty

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type AuthInfo struct {
	vendor     string
	sid        string
	token      string
	appId      string
	templateId string
	timestamp  int
}

func (p AuthInfo) GetSid() string {
	return p.sid
}
func (p *AuthInfo) SetSid(sid string) {
	p.sid = sid
}

func (p *AuthInfo) SetToken(token string) {
	p.token = token
}

func (p AuthInfo) GetToken() string {
	return p.token
}

func (p *AuthInfo) GetTimeSec() string {

	timestamp, _ := strconv.Atoi(time.Now().Format("20060102150405"))
	//大于5分钟 则更新时间戳
	if timestamp-p.timestamp > 300 {
		p.timestamp = timestamp
	}
	return strconv.Itoa(p.timestamp)
}

// SigParameter的生成
func MakeUperMd5SigWithTime(auth AuthInfo) string {
	/* 账户(sid) + 授权令牌(token) + 时间戳 */
	sign := auth.GetSid() + auth.GetToken() + (time.Now().Format("20060102150405"))
	h := md5.New()
	h.Write([]byte(sign))

	SigParameter := hex.EncodeToString(h.Sum(nil))
	return strings.ToUpper(SigParameter)
}

func LoadJsonConf(v *AuthInfo) error {

	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/voiceSender/")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	} else {
		v.sid = viper.GetString(v.vendor + "." + "sid")
		v.token = viper.GetString(v.vendor + "." + "token")
		v.appId = viper.GetString(v.vendor + "." + "appId")
		v.templateId = viper.GetString(v.vendor + "." + "templateId")
	}
	return nil
}
