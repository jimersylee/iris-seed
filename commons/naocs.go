package commons

import (
	"encoding/json"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"strings"
	"time"
)

//var NamingClient = naming_client.NamingClient{}

func RegisterToNacos(config *ServerConfig) {
	// 创建clientConfig
	namespace := ""
	if config.Nacos.Namespace != "public" {
		namespace = config.Nacos.Namespace
	}
	clientConfig := constant.ClientConfig{
		NamespaceId:         namespace, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "info",
	}
	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: config.Nacos.IpAddr,
			Port:   config.Nacos.Port,
			Scheme: "http",
		},
	}

	namingClient, err := clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  &clientConfig,
		ServerConfigs: serverConfigs,
	})
	if err != nil {
		Log.Errorf("client err: %s\n", err)
		return
	}

	//先立马注册上去,然后靠定时任务来更新健康状态和元数据
	//registerInstance(config, namingClient)
	var ch chan int
	ticker := time.NewTicker(time.Second * time.Duration(10))
	go func() {
		for range ticker.C {
			registerInstance(config, namingClient)
		}
		ch <- 1
	}()

}

func registerInstance(config *ServerConfig, namingClient naming_client.INamingClient) {
	//如果有外网ip,则认为是健康的
	Log.Infof("start to registe to nacos")
	health := false
	if ProxyIp != "" {
		health = true
	}
	intranetIp := GetIntranetIp()
	var metaData = map[string]string{}
	reportBody := &ReportBody{}
	reportBody.IntranetIp = intranetIp
	reportBody.InternetIp = ProxyIp
	reportBody.Hostname = config.Report.HostName
	reportBody.Proxies = CreateReportBody(config)
	reportBody.Platform = config.App.Platform
	reportBody.Tags = config.App.Tags
	proxyBytes, _ := json.Marshal(reportBody)
	metaData["info"] = string(proxyBytes)
	metaData["hostname"] = reportBody.Hostname
	metaData["internetIp"] = reportBody.InternetIp
	metaData["intranetIp"] = reportBody.IntranetIp

	if len(reportBody.Tags) != 0 {
		metaData["tags"] = strings.Join(reportBody.Tags, ",")
	}

	regIP := ProxyIp
	if config.Nacos.RegIp == "private" {
		regIP = reportBody.IntranetIp
	}

	param := vo.RegisterInstanceParam{
		Ip:          regIP,
		Port:        config.App.Port,
		Weight:      1,
		Enable:      true,
		Healthy:     health,
		ServiceName: config.App.Name,
		Ephemeral:   true,
		Metadata:    metaData,
	}

	Log.Infof("register instance: %v", param)
	result, err := namingClient.RegisterInstance(param)
	if err != nil || result == false {
		Log.Errorf("注册实例失败, err: %v, result: %v", err, result)
		return
	}
}
