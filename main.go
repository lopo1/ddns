package main

import (
	"ddns/config"
	"ddns/log"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"io"
	"net/http"
	"time"
)

func main() {
	log.InitLogger()
	config.InitConfig()
	//timer := time.NewTimer(1 * time.Second)

	//for {
	//	select {
	//	case <-timer.C:
	//		ip := GetIp()
	//		fmt.Println("ip = ", ip)
	//	}
	//}
	//go func() {
	//	for {
	//		select {
	//		case t := <-timer.C:
	//			fmt.Println("Tick at", t)
	//		}
	//	}
	//
	//	//<-timer.C
	//	//ip := GetIp()
	//	//fmt.Println("ip = ", ip)
	//}()
	//
	//// 让主程序持续运行，等待定时器
	//select {}

	ticker := time.NewTicker(5 * time.Second) // 设置1秒的定时周期

	go func() {
		for {
			select {
			case <-ticker.C:
				ip := GetIp()
				fmt.Println("ip = ", ip)
				err := UpDNS(ip)
				if err != nil {
					log.Logger.Error("UpDNS error = ", err)
				}
			}
		}
	}()

	//time.Sleep(10 * time.Second) // 让主程序运行10秒
	//ticker.Stop()
	fmt.Println("Ticker stopped")
	select {}
}

func GetIp() (ip string) {
	// 设置你的AccessKey ID和AccessKey Secret

	// 获取当前公网IP（这里需要实现获取公网IP的逻辑）
	//publicIP := "<your-public-ip>"
	resp, err := http.Get("http://api.ipify.org")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	ip = string(result)

	fmt.Println("My public IP address is:", ip)
	return

}

func UpDNS(ip string) (err error) {
	accessKeyId := config.Global.Aliyun.Id
	accessKeySecret := config.Global.Aliyun.Secret

	// 创建阿里云DNS客户端
	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", accessKeyId, accessKeySecret)
	if err != nil {
		return err
	}
	// 创建API请求并设置参数
	request := alidns.CreateDescribeDomainRecordsRequest()
	request.Scheme = "https"
	request.DomainName = config.Global.Aliyun.Demain

	// 发送请求并获取响应
	response, err := client.DescribeDomainRecords(request)
	if err != nil {
		return err
	}

	// 遍历DNS记录，并找到需要更新的记录
	for _, record := range response.DomainRecords.Record {

		// 发现匹配的记录，更新它
		updateRequest := alidns.CreateUpdateDomainRecordRequest()
		if record.Value == ip {
			log.Logger.Info("已经改了")
			continue
		}
		updateRequest.Scheme = "https"
		updateRequest.RecordId = record.RecordId
		updateRequest.RR = record.RR
		updateRequest.Type = record.Type

		updateRequest.Value = ip

		_, err := client.UpdateDomainRecord(updateRequest)
		if err != nil {
			return err
		}

		fmt.Println("DNS record updated.")

	}
	return
}
