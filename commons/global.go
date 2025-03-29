package commons

import "github.com/sirupsen/logrus"

// 全局实例对象
var (
	Conf    = &ServerConfig{} // 全局配置文件
	Log     = logrus.New()    // Log API
	ProxyIp = ""              // 代理 ip
)
