package sendMailAndSMS

import (
	"fmt"
	"github.com/c4pt0r/ini"
	"io/ioutil"
	"net/http"
	"net/url"
)

var client = &http.Client{}
var conf = ini.NewConf("config.ini")
var (
	mailsubject = conf.String("mail", "subject", "")
	mailcontent = conf.String("mail", "content", "")
	mailtos     = conf.String("mail", "tos", "")

	smscontent = conf.String("sms", "content", "")
	smstos     = conf.String("sms", "tos", "")
)

func main() {

}

func sendAlarmMail() {
	postValues := url.Values{}
	postValues.Add("subject", *mailsubject)
	postValues.Add("content", *mailcontent)
	postValues.Add("tos", *mailtos) // songhuiqing@social-touch.com
	resp, err := client.PostForm("http://c.fuhaijun.com/mail/", postValues)

	if err != nil {
		fmt.Println(err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
	}
}

func sendAlarmSMS(content string, tos string) {
	postValues := url.Values{}
	postValues.Add("content", *smscontent)
	postValues.Add("tos", *smstos) // 13911045897
	resp, err := client.PostForm("http://c.fuhaijun.com/sms/", postValues)

	if err != nil {
		fmt.Println(err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
	}
}
