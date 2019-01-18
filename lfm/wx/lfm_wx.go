package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var db *sql.DB
var m1 map[int]int
var client = &http.Client{}

func init() {
	db, _ = sql.Open("mysql", "admin:TH)6*Ca($.$u5)kA)bb+X%k[$wWY45@tcp(211.151.64.236:3306)/category?charset=utf8")
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	db.Ping()

	rows, err := db.Query("select id, cateid, device_number from label_device_number where inc_date = '" + GetYesDate() + "' and exchangeid = 100004")
	defer rows.Close()
	checkErr(err)
	m1 = make(map[int]int)
	for rows.Next() {
		var id int
		var cate_id int
		var device_number int
		err = rows.Scan(&id, &cate_id, &device_number)
		checkErr(err)
		m1[cate_id] = device_number
	}
}

// 得到今天的前一天(日期) 比如今天是20160901 得到的日期为20160831
func GetYesDate() string {
	nTime := time.Now()
	yesTime := nTime.AddDate(0, 0, -1)
	return yesTime.Format("20060102")
}

// 检查文件或目录是否存在
// 如果由 filename 指定的文件或目录存在则返回 true，否则返回 false
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

// get file size as how many bytes
func FileSize(file string) (int64, error) {
	f, e := os.Stat(file)
	if e != nil {
		return 0, e
	}
	return f.Size(), nil
}

// 逐行读取文件
func ReadLine(fileName string, handler func(string, string), channelId string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		handler(line, channelId)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
	return nil
}

func pdataAndInDB(line string, channelId string) {
	if line != "" {
		str := strings.Split(line, "\t")
		s0, _ := strconv.Atoi(str[0])
		s1, _ := strconv.Atoi(str[1])
		if v, ok := m1[s0]; ok {
			sum := v + s1
			stmt, err := db.Prepare("update label_device_number set device_number = ? where inc_date = ? and exchangeid = ? and cateid = ?")
			checkErr(err)
			stmt.Exec(sum, GetYesDate(), channelId, s0)
			stmt.Close()
		} else {
			stmt, err := db.Prepare("insert into label_device_number(exchangeid,cateid,device_number,inc_date) values(?,?,?,?)")
			checkErr(err)
			stmt.Exec(channelId, s0, s1, GetYesDate())
			stmt.Close()
		}
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func sendAlarmMail(content string) {
	postValues := url.Values{}
	postValues.Add("subject", "adx通用兴趣流量监控告警")
	postValues.Add("content", "")
	postValues.Add("tos", "songhuiqing@social-touch.com") // songhuiqing@social-touch.com
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

func lfm_adview() {
	adviewFile := "/data/hadoop/analysis/userprofile/user_profile_intt/output/adx/adview/g2_inc_" + GetYesDate() + ".log"
	if Exist(adviewFile) == true {
		if size, _ := FileSize(adviewFile); size > 0 {
			ReadLine(adviewFile, pdataAndInDB, "100004")
		} else {
			fmt.Println("adviewFile size <= 0")
		}
	} else {
		sendAlarmMail("adviewFile没有生成，请检查具体原因!")
	}
}

func main() {
	lfm_adview()
}
