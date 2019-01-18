package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/c4pt0r/ini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/cron"
)

var db *sql.DB
var client = &http.Client{}
var conf = ini.NewConf("cpa_target_features.ini")
var dbName = "mysql-p"
var (
	dataSourceName = conf.String(dbName, "dataSourceName", "dataSourceName")
	maxOpenConns   = conf.Int(dbName, "maxOpenConn", 2000)
	maxIdleConns   = conf.Int(dbName, "maxIdleConn", 1000)

	mailsubject = conf.String("mail", "subject", "")
	mailcontent = conf.String("mail", "content", "")
	mailtos     = conf.String("mail", "tos", "")

	smscontent = conf.String("sms", "content", "")
	smstos     = conf.String("sms", "tos", "")
)

func init() {
	conf.Parse()
	db, _ = sql.Open("mysql", *dataSourceName)
	db.SetMaxOpenConns(*maxOpenConns)
	db.SetMaxIdleConns(*maxIdleConns)
	db.Ping()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func truncateCTF() {
	res, err := db.Exec("TRUNCATE table cpa_target_features")
	checkErr(err)
	fmt.Println(res.LastInsertId())
}

func insertCTF() {
	/*res, err := db.Exec("insert INTO cpa_target_features(dsp_id,campaign_id,target_features,nshow,nclick,nconvert,cost,cpa) " +
	"select dsp_id, tt1.campaign_id,target_features,nshow,nclick,nconvert,cost,cpa " +
	"from cpa_target_features_raw tt1, " +
	"(select t1.campaign_id as campaign_id, " +
	"		 t2.cpa_lower_exp as cpa_lower_exp, " +
	"		 t2.cpa_upper_exp as cpa_upper_exp " +
	"from edison.ad_campaign t1, edison.ad_plan t2 " +
	"where t1.plan_id=t2.plan_id and t2.cpa_lower_exp>0 and t2.cpa_upper_exp>0 " +
	"group by t1.campaign_id, t2.cpa_lower_exp, t2.cpa_upper_exp " +
	") tt2 " +
	"where tt1.campaign_id=tt2.campaign_id and cpa<=(tt2.cpa_upper_exp/100)*2")*/

	res, err := db.Exec("insert INTO cpa_target_features(dsp_id,campaign_id,target_features,nshow,nclick,nconvert,cost,cpa,type) " +
		"select dsp_id, tt1.campaign_id,target_features,nshow,nclick,nconvert,cost,cpa, " +
		"case " +
		"when nconvert>0 and cpa>(tt2.cpa_upper_exp/100)*2 then 1 " +
		"when nconvert>0 and cpa<=(tt2.cpa_upper_exp/100)*2 then 2 " +
		"else 3 " +
		"end as type " +
		"from cpa_target_features_raw tt1,  " +
		"(select t1.campaign_id as campaign_id,  " +
		"t2.cpa_lower_exp as cpa_lower_exp,  " +
		"t2.cpa_upper_exp as cpa_upper_exp  " +
		"from edison.ad_campaign t1, edison.ad_plan t2  " +
		"where t1.plan_id=t2.plan_id and t2.cpa_lower_exp>0 and t2.cpa_upper_exp>0 " +
		"group by t1.campaign_id, t2.cpa_lower_exp, t2.cpa_upper_exp  " +
		") tt2  " +
		"where tt1.campaign_id=tt2.campaign_id and (cpa<=(tt2.cpa_upper_exp/100)*2 or nconvert>0)")
	checkErr(err)
	id, err := res.LastInsertId()
	if err != nil {
		fmt.Println("insertCTF err")
		checkErr(err)
		sendAlarmMail()
	} else {
		if id <= 0 {
			fmt.Println("id <= 0")
			//sendAlarmMail()
		}
	}
}

func main() {
	ec := flag.String("ec", "1", "minute") // 执行周期
	flag.Parse()

	fmt.Println(*ec)
	spec := "0 */" + *ec + " * * * " // 每隔*ec分钟执行一次
	c := cron.New()
	c.AddFunc(spec, ctf)
	c.Start()
	select {}
}

func ctf() {
	// TRUNCATE 表cpa_target_features
	truncateCTF()

	// 往cpa_target_features表中写入数据
	insertCTF()
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
