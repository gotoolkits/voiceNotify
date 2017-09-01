package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/gotoolkits/voiceSender/thirdParty"
)

func main() {

	var vendor, phone string
	flag.StringVar(&vendor, "v", "", "The third-party vendor name")
	flag.StringVar(&phone, "p", "", "Receive phone numbers")

	flag.Parse()

	if vendor == "" || phone == "" {
		flag.Usage()
	}

	if strings.Contains(vendor, "ucpaas") {

		data := `{"messages":"this is test"}`
		body := thirdParty.InitUcpassAndNewReqJsonBody(phone, data)

		resp, err := thirdParty.PostMsgToUcpaas(body)

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(resp))

	}

	if strings.Contains(vendor, "ihuyi") {
		thirdParty.InitAuthInfo()
		thirdParty.PostMsgToIhuyi(phone)
	}

	if strings.Contains(vendor, "ytx") {

		body := thirdParty.InitYtxAndNewReqJsonBody(phone)
		resp, err := thirdParty.PostMsgToYtx(body)

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(resp))

	}

}
