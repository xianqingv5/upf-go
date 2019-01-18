package main

import (
	"database/sql"
	"fmt"
	"github.com/Shopify/sarama"
	. "github.com/aerospike/aerospike-client-go"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"strings"
	"sync"
)

var (
	wg        sync.WaitGroup
	logger    = log.New(os.Stderr, "[upf]", log.LstdFlags)
	topicName = "effectcallbackLog"

	m1 map[string]string

	db *sql.DB
)

func init() {
	logger.Println("init")

	db, _ = sql.Open("mysql", "admin:TH)6*Ca($.$u5)kA)bb+X%k[$wWY45@tcp(211.151.64.236:3306)/category?charset=utf8")
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	db.Ping()

	rows, err := db.Query("select id, ad_name from edison.ad_user where industry = '19'")
	PanicOnError(err)

	m1 = make(map[string]string)
	for rows.Next() {
		var id string
		var ad_name string
		err = rows.Scan(&id, &ad_name)
		PanicOnError(err)
		m1[id] = ad_name
		logger.Printf("id:%s, ad_name:%s", id, ad_name)
	}

}

func main() {
	sarama.Logger = logger
	conf := sarama.NewConfig()
	conf.ClientID = "upf-go"
	// consumer, err := sarama.NewConsumer(strings.Split("madwx21:9092,madwx61:9092,madwx71:909", ","), conf)
	consumer, err := sarama.NewConsumer(strings.Split("mad101:9092,mad102:9092,mad103:909", ","), conf)
	if err != nil {
		logger.Printf("Failed to start consumer: %s", err)
	}
	partitionList, err := consumer.Partitions(topicName)
	if err != nil {
		logger.Println("Failed to get the list of partitions: ", err)
	}

	for partition := range partitionList {
		pc, err := consumer.ConsumePartition(topicName, int32(partition), sarama.OffsetNewest)
		if err != nil {
			logger.Printf("Failed to start consumer for partition %d: %s\n", partition, err)
		}
		defer pc.AsyncClose()
		wg.Add(1)
		go func(sarama.PartitionConsumer) {
			defer wg.Done()
			for msg := range pc.Messages() {
				// fmt.Printf("Partition:%d, Offset:%d, Key:%s, Value:%s", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
				parseLog(string(msg.Value))
			}
		}(pc)
	}

	wg.Wait()
	logger.Println("Done consuming topic effectcallbackLog")
	consumer.Close()
}

func parseLog(msg string) {
	s := strings.Split(msg, "\t")
	adSpaceId := s[10]
	adId := s[11]
	eventType := s[22]
	imei := s[16]
	idfa := s[17]
	platform := s[18]
	ip := s[29]
	reqIp := s[31]
	uc := s[32]

	if eventType == "3" { // 效果数据
		fmt.Printf("广告位:%s, 广告主:%s, 事件类型:%s, imei:%s, idfa:%s, 操作系统:%s, 渠道ID:%s, 流量来源:%s, ip:%s, reqIp:%s, uc:%s", adSpaceId, adId, eventType, imei, idfa, platform, s[7], s[5], ip, reqIp, uc)
		fmt.Println()

		if platform == "2" && len(idfa) == 36 {
			// 查找键值是否存在
			if v, ok := m1[adId]; ok {
				fmt.Println(v)
				toTag(idfa) // 实时打标签
			}
		} else if platform == "2" && idfa == "[PASS_IDFA_HERE]" {
			if v, ok := m1[adId]; ok {
				fmt.Println(v)
				if reqIp != "^" {
					toTag(reqIp)
				} else {
					toTag(ip)
				}

			}
		} else if platform == "2" && idfa == "^" {
			if v, ok := m1[adId]; ok {
				fmt.Println(v)
				if reqIp != "^" {
					toTag(reqIp)
				} else {
					toTag(ip)
				}
			}
		}
	}
}

func toTag(sbid string) {
	// define a client to connect to
	client, err := NewClient("192.168.0.12", 3000)
	PanicOnError(err)
	defer client.Close()
	key, err := NewKey("upf", "mwuser", sbid)
	PanicOnError(err)

	// define some bins with data
	bins := BinMap{
		"gcr": "1",
	}

	// write the bins
	err = client.Put(nil, key, bins)
	PanicOnError(err)

}

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
