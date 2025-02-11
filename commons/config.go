package commons

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

const (
	ProdEnv                = "prod"
	AppPort                = 10000
	ReportDuration         = 30
	SS5ProxyPort           = 19825
	SS5ProxyConnectTimeOut = 5
	HttpProxyPort          = 21520
	LogsPath               = "./logs"
	InfoLogLevel           = "info"
)

type ServerConfig struct {
	App struct {
		Env             string   `yaml:"env"`             // 运行环境，dev 与 prod, 默认为生产环境(prod)，如果为开发环境，不生成文件日志
		Port            uint64   `yaml:"port"`            // 应用本身暴露出的 http 服务的端口
		ChangeIpCommand string   `yaml:"changeIpCommand"` // agent 端用来更换 ip 的命令
		Platform        string   `yaml:"platform"`        //代理运行的云服务平台
		Tags            []string `yaml:"tags"`            // 标签
		Name            string   `yaml:"name" default:""`
	} `yaml:"app"`

	Report struct {
		Url      string `yaml:"url"`      // 上报的 url
		HostName string `yaml:"hostname"` // 上报信息中的 hostname, 如果没有指定，就获取本机的 hostname
		Ip       string `yaml:"ip"`       // 代理ip，如果不指定，那就请求 http://httpbin.org/ip 获取自己的外网 ip
		Duration int    `yaml:"duration"` // 上报代理信息的周期
	} `yaml:"report"`

	NeedReport bool

	Proxy struct {
		Socks5 struct {
			Enable         bool     `yaml:"enable"` // 是否开启 socks5 代理
			Port           int      `yaml:"port"`
			ConnectTimeout int      `yaml:"connectTimeout"` // 代理服务器与目标网站连接超时时间, 默认5秒
			UserAndPass    []string `yaml:"userAndPass"`
			UserAndPassMap map[string]string
		} `yaml:"socks5"`
		Http struct {
			Enable         bool     `yaml:"enable"`
			Port           int      `yaml:"port"`
			UserAndPass    []string `yaml:"userAndPass"`
			UserAndPassMap map[string]string
		} `yaml:"http"`
		ShadowSocks struct {
			Enable    bool     `yaml:"enable"`    //是否开启shadowSocks代理
			Port      int      `yaml:"port"`      //监听端口
			Cipher    string   `yaml:"cipher"`    //加密方式
			Password  []string `yaml:"password"`  //密码
			WhiteList []string `yaml:"whiteList"` //域名白名单
			BlackList []string `yaml:"blackList"` //域名黑名单
		} `yaml:"shadowSocks"`
	} `yaml:"proxy"`
	Log struct {
		Path  string `yaml:"path"`  // 日志存储的文件目录, 默认为当前目录下的 logs 目录下
		Level string `yaml:"level"` // 日志打印级别, 默认为 info
		SLS   struct {
			AccessKeyId     string `yaml:"accessKeyId"`
			AccessKeySecret string `yaml:"accessKeySecret"`
			Project         string `yaml:"project"`
			Logstore        string `yaml:"logstore"`
			Endpoint        string `yaml:"endpoint"`
			Topic           string `yaml:"topic"`
		}
	} `yaml:"log"`
	Nacos struct {
		Namespace string `yaml:"namespace"` //nacos的命名空间,public或者特定的命名空间
		IpAddr    string `yaml:"ipAddr"`    // nacos服务的ip或者域名
		Port      uint64 `yaml:"port"`      //nacos服务的端口
		RegIp     string `yaml:"regIp"`     //注册到nacos的ip, private / public 内网或者外网，默认外网
	} `yaml:"nacos"`
}

func (conf *ServerConfig) InitConfig(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(file, conf); err != nil {
		return err
	}
	// 初始化两个 map
	conf.Proxy.Socks5.UserAndPassMap = make(map[string]string)
	conf.Proxy.Http.UserAndPassMap = make(map[string]string)

	for _, info := range conf.Proxy.Socks5.UserAndPass {
		split := strings.Split(info, ":")
		conf.Proxy.Socks5.UserAndPassMap[split[0]] = split[1]
	}

	for _, info := range conf.Proxy.Http.UserAndPass {
		split := strings.Split(info, ":")
		conf.Proxy.Http.UserAndPassMap[split[0]] = split[1]
	}

	// 初始化一些默认值，如果用户没有设置

	// 应用相关设置
	if conf.App.Env == "" {
		// 环境设置，默认为生产环境
		conf.App.Env = ProdEnv
	}
	if conf.App.Port == 0 {
		conf.App.Port = AppPort
	}

	// 上报功能相关设置
	if conf.Report.Url != "" {
		conf.NeedReport = true
	}

	if conf.Report.HostName == "" {
		hostname, err := os.Hostname()
		if err != nil {
			return err
		}
		conf.Report.HostName = hostname
	}

	if conf.Report.Duration == 0 {
		conf.Report.Duration = ReportDuration
	}

	// 代理相关设置
	if conf.Proxy.Socks5.Enable {
		if conf.Proxy.Socks5.Port == 0 {
			conf.Proxy.Socks5.Port = SS5ProxyPort
		}

		if conf.Proxy.Socks5.ConnectTimeout == 0 {
			conf.Proxy.Socks5.ConnectTimeout = SS5ProxyConnectTimeOut
		}
	}

	if conf.Proxy.Http.Enable {
		if conf.Proxy.Http.Port == 0 {
			conf.Proxy.Http.Port = HttpProxyPort
		}
	}
	// 不再设置默认 logPath
	//if conf.Log.Path == "" {
	//	conf.Log.Path = LogsPath
	//}
	if conf.Log.Level == "" {
		conf.Log.Level = InfoLogLevel
	}

	return nil
}

// PrettyString returns a pretty-printed JSON representation of the ServerConfig
func (conf *ServerConfig) PrettyString() string {
	jsonBytes, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error marshaling config to JSON: %v", err)
	}
	return string(jsonBytes)
}
