package main

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"time"
)

func main() {
	registerToNacos()
	time.Sleep(1 * time.Hour)
}
func registerToNacos() {
	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         "", // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}
	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: "192.168.198.134",
			Port:   8848,
			Scheme: "http",
		},
	}

	namingClient, err := clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  &clientConfig,
		ServerConfigs: serverConfigs,
	})
	if err != nil {
		fmt.Printf("client err: %s\n", err)
		return
	}

	result, err := namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          "192.168.199.128",
		Port:        8001,
		Weight:      1,
		Enable:      true,
		Healthy:     true,
		ServiceName: "choco-proxy",
		Ephemeral:   true,
	})
	if err != nil {
		fmt.Printf("eeee %s\n", err)
		return
	}
	if result == true {
		fmt.Println("register to nacos success")
	}

}
