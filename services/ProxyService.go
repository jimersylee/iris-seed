package services

import (
	"fmt"
	"github.com/kataras/iris"
)

type ProxyInterface interface {
	//注册ip
	ProxyCommunity(ctx iris.Context) bool
	ProxyApi(ctx iris.Context)
}

var ProxyService = NewProxyService()

func NewProxyService() *ProxyServiceImpl {
	return &ProxyServiceImpl{}
}

type ProxyServiceImpl struct {
}

func (p *ProxyServiceImpl) ProxyCommunity(ctx iris.Context) bool {

	panic("implement me")
}

func (p *ProxyServiceImpl) ProxyApi(ctx iris.Context) {
	uri := ctx.Request().RequestURI
	fmt.Println(uri)
	uriNeed := uri[len("/api/steamapi/") : len(uri)-1]
	route := ctx.GetCurrentRoute()
	fmt.Println(route)
	fmt.Println(uriNeed)
	//todo：找出能用的ip
	ip := "182.255.44.92"
	webUrl := "http://api.steampowered.com/" + uriNeed
	responseStr := IpService.Proxy(ip, webUrl)
	fmt.Println(responseStr)
	ctx.WriteString(responseStr)
}
