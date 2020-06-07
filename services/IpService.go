package services

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type IpServiceInterface interface {
	//注册ip
	Register(ip string) bool
	Check(ip string) bool
}

var IpService = NewIpService()

func NewIpService() *IpServiceImpl {
	return &IpServiceImpl{}
}

type IpServiceImpl struct {
}

func (i *IpServiceImpl) Register(ip string) bool {
	return true
}

func (i *IpServiceImpl) Check(ip string) bool {
	return true
}

func (i *IpServiceImpl) Proxy(ip string, webUrl string) string {
	proxyUrl := "http://" + ip + ":60002"
	fmt.Println("use " + proxyUrl)
	proxy, _ := url.Parse(proxyUrl)
	tr := &http.Transport{
		Proxy:           http.ProxyURL(proxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	request, _ := http.NewRequest("GET", webUrl, nil)
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36")
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 5, //超时时间
	}

	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("出错了", err)
		return ""
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(body)
	fmt.Println(bodyString)
	return bodyString
}
